package handler

import (
	"fmt"
	//"github.com/go-openapi/strfmt"
	"github.com/stevef1uk/cassuservice/parser"
	"github.com/stevef1uk/cassuservice/swagger"
	"strconv"

	//"github.com/stevef1uk/cassuservice/swagger"
	"os"
	"strings"

	//"regexp"
	//"strings"
	//"log"
)

// Function to create a file that will contain the Cassandra handler code
func CreateFile( debug bool, pathPrefix string, dir string, fileName string ) *os.File {

	// Create the directory if not already there
	fulldirName := pathPrefix  + dir
	if debug { fmt.Println("Dir Name = ", fulldirName )}
	// Create data dir if it doesn't already exist
	if _, err := os.Stat(fulldirName); err != nil {
		if os.IsNotExist(err) {
			// dir does not exist
			if err := os.MkdirAll( fulldirName , 0755) ; err != nil  {
				panic(err)
			}
		}
	}
	fullFileName := fulldirName + "/" + fileName
	if debug { fmt.Println("Generated file name  = ", fullFileName )}
	// Save previous generated file
	os.Rename( fullFileName, fullFileName+".old")
	// Now create new file
	var file, err = os.Create(fullFileName)
	if err != nil  {
		panic(err)
	}

	return file
}


func getServiceName ( tableName, endPointNameOverRide string ) string {
	ret := tableName
	if ( endPointNameOverRide != "") {
		ret = strings.Title( endPointNameOverRide  )
	}
	return ret
}



func GetFieldName(  debug bool, leaveCase bool, fieldName string, dontUpdate bool ) string {
	name := fieldName
	if ! leaveCase {
		name = strings.ToLower(fieldName)
	}
	return CapitaliseSplitFieldName( debug, name, dontUpdate )
}



// Function that renames fields to match that performed for some reason by go-swagger in its generated framework code
// Go-Swagger turns any field that starts with 'id' into ID and this is true for Table names, but not for Types for the first id, but if one has id_id it does produce IdID
// The first character is always capitalised
func Capitiseid( debug bool, fieldName string, dontUpdate bool ) string {

	var ret string = ""
	if debug { fmt.Printf("Capitiseid entry field  = %s, len = %d\n ",fieldName, len(fieldName) ) }

	runes := []rune(fieldName[:])
	last := len( runes  ) - 1

	if last <= 1 {
		if last == 0 {
			return strings.ToUpper(string(runes[0]))
		}
	}

	if ! dontUpdate {
		if (runes[0] == rune('i') || runes[0] == rune('I')) && (runes[1] == rune('d') || runes[1] == rune('D')) {
			runes[0] = rune('I')
			runes[1] = rune('D')
		}
		/*
		for i := 0; i < last; i++ {
			if ! ((i == 0) || (i == last-1)) {
				continue;
			}
			//if debug { fmt.Printf("Capitiseid [0] = %q, [1] = %q\n ", runes[i], runes[i+1]) }
			if (runes[i] == rune('i') || runes[i] == rune('I')) && (runes[i+1] == rune('d') || runes[i+1] == rune('D')) {
				if debug {
					fmt.Printf("Capitiseid match at i= %d\n ", i)
				}
				runes[i] = rune('I')
				runes[i+1] = rune('D')
			}
		}*/

	}

	ret = strings.ToUpper(string(runes[0])) + string(runes[1:])
	if debug {fmt.Printf("Capitiseid returning field  = %s\n ", ret)}
	return ret
}

// Function that renames fields to match that performed for some reason by go-swagger in its generated framework code e.g. My_List becomes MyList & address_id becomes AddressID
func CapitaliseSplitFieldName ( debug bool, fieldName string, dontUpdate bool ) string {
//@TODO remove
debug = false
	var ret string = ""
	if debug { fmt.Printf("CapitaliseSplitFieldName entry field  = %s, len = %d\n ",fieldName, len(fieldName) ) }

	if fieldName == ""{
		ret = fieldName
		if debug { fmt.Printf("CapitaliseSplitFieldName fieldName empty\n ") }
	} else {
		tmpFields := strings.Split(fieldName, "_" )
		if debug {fmt.Printf("CapitaliseSplitFieldName tmpFields  = %q\n ", tmpFields)}
		first := dontUpdate
		for _, v := range tmpFields {
			v = Capitiseid( debug, v, first )
			first = false
			v = strings.ToUpper(string(v[0])) + v[1:]
			ret = ret + v
		}
		ret = strings.ToUpper(string(ret[0])) + ret[1:]

	}


	if debug {fmt.Printf("Capitiseid returning field  = %s\n ", ret)}
	return ret
}





// THis function returns the Go types for UDT fields
func basicMapCassandraTypeToGoType( debug bool, leaveFieldCase bool, inTable bool, fieldName string, fieldType string, typeName string,  fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, dontUpdate bool, makeSmall bool ) string {
	text := ""
	leaveCase := false
	if leaveFieldCase {
		leaveCase = true
	}

	if debug {fmt.Printf("basicMapCassandraTypeToGoType %s %s\n ", fieldName,fieldType )}
	switch strings.ToLower(fieldType) {
	case "int": fallthrough
	case "bigint": fallthrough
	case "counter": fallthrough
	case "varint":
		if makeSmall {
			text = "int"
		} else {
			text = "int64"
		}
	case "uuid":
		text = "string"
	case "date": fallthrough
	case "timestamp":
		text = "time.Time"
	case "timeuuid":
		text = "gocql.UUID"
	case "boolean":
		text = "bool"
	//case "decimal":
		//text = "*inf.Dec" // this is in the gopkg.in/inf.v0 package
	case "float": fallthrough
	case "double":
		if makeSmall {
			text = "float32"
		} else {
			text = "float64"
		}
	case "decimal":
		text = "*inf.Dec"
	case "text":
		text = "string"
	case "varchar":
		text = "string"
	case "ascii":
		text = "string"
	case "blob":
		text = "string"
	case "inet":
		text = "string"
	case "smallint":
		text = "int16"
	case "list": fallthrough
	case "set":
		if ! swagger.IsFieldTypeUDT( parserOutput, fieldDetails.DbFieldCollectionType) {
			if ! makeSmall {
				text = "[]"
			}
			text =  text + basicMapCassandraTypeToGoType( debug, true, inTable, fieldName, fieldDetails.DbFieldCollectionType, typeName,   fieldDetails , parserOutput, dontUpdate, makeSmall )
		} else {
			fieldName = GetFieldName( debug, leaveCase, fieldName, dontUpdate)
			text = text + MODELS
			if ! inTable {
				text = text + fieldName
			} else {
				text = text + typeName + fieldName
			}
		}
		//text = text + basicMapCassandraTypeToGoType( debug, true, inTable, fieldName, fieldDetails.DbFieldCollectionType, typeName,   fieldDetails , parserOutput, dontUpdate )
	case "map":
		fieldName = GetFieldName( debug, leaveCase, fieldName, dontUpdate)
		if inTable {
			text = MODELS + fieldName
		} else {
			typeName = GetFieldName( debug, leaveCase, typeName, dontUpdate)
			text = MODELS + typeName + fieldName
		}
	default:
		if debug {fmt.Printf("basicMapCassandraTypeToGoType TYPE NOT MATCHED!!!!\n " )}
		fieldName = GetFieldName( debug, leaveCase, fieldName, dontUpdate)
		if inTable {
			text =  MODELS + fieldName
		} else {
			typeName = GetFieldName( debug, leaveCase, typeName, dontUpdate)
			if fieldDetails.DbFieldCollectionType != "" || fieldDetails.DbFieldMapType != "" {
				text =  MODELS + typeName + fieldName
			} else {
				text = typeName
			}
		}

		//panic(1)
	}

	if debug { fmt.Printf("basicMapCassandraTypeToGoType returning %s from field %s type %s\n", text, fieldName, fieldType ) }
	return text
}


// The Go types are different for UDT types and table field types in some cases. This function deals with table field return types
func mapTableTypeToGoType( debug bool, fieldName string, fieldType string, typeName string, fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, dontUpdate bool  ) string {

	text := ""

	switch strings.ToLower(fieldType) {

	default:
		text = basicMapCassandraTypeToGoType( debug, true, true, fieldName, fieldType, typeName,   fieldDetails , parserOutput, dontUpdate, false )
	}

	if debug { fmt.Printf("mapTableTypeToGoType returning %s from field %s type %s\n", text, fieldName, fieldType ) }
	return text
}


func mapFieldTypeToGoCSQLType( debug bool, fieldName string, leaveFieldCase bool, inTable bool, fieldType string, typeName string, fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, dontUpdate bool  ) string {

	text := ""
	if debug {fmt.Printf("mapFieldTypeToGoCSQLType %s %s\n ", fieldName,fieldType )}

	switch strings.ToLower(fieldType) {
	case "int": fallthrough
	case "varint":
		text = "int"
	case "uuid":
		text = "string"
	case "date":
		text = "strfmt.DateTime"
	case "timeuuid":
		text = "gocql.UUID"
	case "float":
		text = "float32"
	case "list": fallthrough
	case "set":
		if ! swagger.IsFieldTypeUDT( parserOutput, fieldDetails.DbFieldCollectionType) {
			text = "[]"
		}
		retType := basicMapCassandraTypeToGoType( debug, true, true, fieldName, fieldType, typeName,   fieldDetails , parserOutput, dontUpdate, true )
		text = text + retType
	default:
		text = basicMapCassandraTypeToGoType( debug, true, inTable, fieldName, fieldType, typeName,   fieldDetails , parserOutput, dontUpdate, true )
	}

	if debug { fmt.Printf("mapFieldTypeToGoCSQLType returning %s from field %s type %s\n", text, fieldName, fieldType ) }
	return text
}


// Function to return a temporary variable based on string
var counter int = 0
func createTempVar ( fieldName string ) string {
	ret := TEMP_VAR_PREFIX + fieldName + "_" + strconv.Itoa( counter)
	counter = counter + 1
	return ret
}


func existsTimeField( fieldDetails parser.FieldDetails  ) bool {
	ret := false

	if ( swagger.IsFieldTypeATime( fieldDetails.DbFieldType ) ||
		 swagger.IsFieldTypeATime( fieldDetails.DbFieldCollectionType ) ||
		 swagger.IsFieldTypeATime( fieldDetails.DbFieldMapType ) ) {
		ret = true
	}
	return ret
}

/*
func existsFieldType( fieldDetails parser.FieldDetails, fieldType string  ) bool {
	ret := false
	
	if swagger.IsFieldTypeATime( strings.ToUpper( fieldType ) ) {
		ret = existsTimeField( fieldDetails )
	} else if ( ( strings.ToLower( fieldDetails.DbFieldType ) == fieldType ) ||
		 ( strings.ToLower( fieldDetails.DbFieldCollectionType ) == fieldType ) ||
		 (  strings.ToLower( fieldDetails.DbFieldMapType ) == fieldType ) ) {
		ret = true
		}
	return ret
}
*/

func existsFieldType( fieldDetails parser.FieldDetails, fieldType string  ) bool {
	ret := false

	if  ( ( strings.ToLower( fieldDetails.DbFieldType ) == fieldType ) ||
		( strings.ToLower( fieldDetails.DbFieldCollectionType ) == fieldType ) ||
		(  strings.ToLower( fieldDetails.DbFieldMapType ) == fieldType ) ) {
		ret = true
	}
	return ret
}


// Scan through fields and UDT fields to see if a type contained is a time type. Return true if a field is a time field
func doINeedTime(  parserOutput parser.ParseOutput   ) bool {
	ret := false
	for _, v := range parserOutput.TableDetails.TableFields.DbFieldDetails {
		if existsTimeField(v) {
			ret = true;
			break;
		}
	}
	if ! ret {
		for _, v := range parserOutput.TypeDetails {
			for _, k := range v.TypeFields.DbFieldDetails {
				if existsTimeField(k) {
					ret = true;
					break;
				}
			}
			if ret {
				break
			}
		}

	}
	return ret
}

// Scan through table fields to see if type is a float. Return true if a field is float
func doIHaveFloat(  parserOutput parser.ParseOutput   ) bool {
	ret := false
	for _, v := range parserOutput.TableDetails.TableFields.DbFieldDetails {
		if strings.ToLower( v.DbFieldType ) == "float" {
			ret = true;
			break;
		}
	}

	return ret
}


//Scan through fields and UDT fields to see if a type contained is a decimal. Return true if a field is a decimal
func doINeedDecimal(  parserOutput parser.ParseOutput  ) bool {
	ret := false
	for _, v := range parserOutput.TableDetails.TableFields.DbFieldDetails {
		if existsFieldType( v , swagger.DECIMAL ) {
			ret = true;
			break;
		}
	}
	if ! ret {
		for _, v := range parserOutput.TypeDetails {
			for _, k := range v.TypeFields.DbFieldDetails {
				if existsFieldType( k, swagger.DECIMAL ) {
					ret = true;
					break;
				}
			}
			if ret {
				break
			}
		}

	}
	return ret
}


func ProcessTime ( firstTime bool , indent string, timeVar string, typeName string, fieldName string ) (string, string)  {

	equals := " = "
	/*if firstTime {
		equals = " := "
	}*/
	equals = " := "
	if typeName != "" {
		typeName = typeName + "."
	}
	ret := indent + timeVar  + " = " + typeName + fieldName + ".String()"

	tmpV := createTempVar( fieldName )
	tmpV2 := createTempVar( fieldName )
	tmpV3 := createTempVar( fieldName )
	ret = ret + indent + tmpV + " := " + timeVar + `[0:10] + "T" + ` + timeVar + `[11:19] + "." + ` + timeVar + "[20:22]"
	ret = ret + indent + "if " + timeVar + "[22] == ' ' " + "{" +  indent + "  " + timeVar + " = " + tmpV  + ` + "0" + "Z" ` +
		indent + `} else { ` + indent + "  "  + timeVar  + " = " +  tmpV + ` + "Z"` + indent + "}"
	ret = ret + indent + tmpV2 + ", _ " + equals + "strfmt.ParseDateTime(" + timeVar + ")"
	ret = ret + indent + tmpV3 + " := " + tmpV2 + ".String()"



	return ret, tmpV3
}

func findTypeDetails ( debug bool, typeName string, parserOutput parser.ParseOutput ) *parser.TypeDetails {
	ret := &parser.TypeDetails{}
	typeName = strings.ToUpper(( typeName ))
	for i := 0; i < parserOutput.TypeIndex; i++ {
		if  typeName == parserOutput.TypeDetails[i].TypeName {
			ret = &parserOutput.TypeDetails[i]
			break
		}
	}
	return ret
}



func CopyArrayElements( debug bool, inTable bool, inDent string, sourceFieldName string, destFieldName string,  fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput  ) string {
	equals := " := "
	ret := ""
	arrayType := basicMapCassandraTypeToGoType(debug, false, inTable, fieldDetails.DbFieldName, fieldDetails.DbFieldCollectionType, "", fieldDetails, parserOutput, true, true )
	arrayTypeDest := basicMapCassandraTypeToGoType(debug, false, inTable, fieldDetails.DbFieldName, fieldDetails.DbFieldCollectionType, "", fieldDetails, parserOutput, false, false )

	switch arrayType {
	case "*inf.Dec":
		arrayTypeDest = "float64"
	case "gocql.UUID": fallthrough
	case "init8":
		arrayTypeDest = "string"
	}

	if inTable {
		ret = INDENT_1 + inDent + sourceFieldName + equals +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].([]` + arrayType + ")"
	}

	ret = ret + inDent + destFieldName + " = " + "make([] " + arrayTypeDest + ", len(" + sourceFieldName + ") )"
	ret = ret + inDent + "for j := 0; j < len(" + sourceFieldName + " ); j++ { "

	switch arrayType {
	case "*inf.Dec": fallthrough
	case "gocql.UUID": fallthrough
	case "uint8":
		arrayTypeDest = arrayType
	}

	switch arrayTypeDest {
	case "float64":
		ret = ret + inDent + INDENT + destFieldName + "[j] = " +  "float64(" + sourceFieldName + "[j])" + inDent + "}"
	case "int64":
		ret = ret + inDent + INDENT + destFieldName + "[j] = " +  "int64(" + sourceFieldName + "[j])" + inDent + "}"
	case "*inf.Dec":
		tmp_var := createTempVar( sourceFieldName )
		ret = ret + inDent + INDENT + tmp_var + ",_ := " +  "strconv.ParseFloat(" + sourceFieldName + "[j].String(), 64 )"
		ret = ret + inDent + INDENT + destFieldName + "[j] = " +  tmp_var +  inDent + "}"
	case "uint": fallthrough
	case "time.Time": fallthrough
	case "gocql.UUID":
		tmp_var := createTempVar( sourceFieldName )
		ret = ret + inDent + INDENT + tmp_var + " := string(" + sourceFieldName + "[j] )"
		ret = ret + inDent + INDENT + destFieldName + "[j] = " +  tmp_var +  inDent + "}"
	default:
		if debug {fmt.Printf("CopyArrayElements TYPE NOT MATCHED!!!!\n " )}
		ret = ret + inDent + INDENT + destFieldName + "[j] = " + sourceFieldName + "[j]" + inDent + "}"
	}
	return ret
}


func applyTypeConversionForGoSwaggerToGocql( debug bool, output string, suffix string, fieldName string,  fieldType string ) string {

	ret := output
	if debug {fmt.Printf("mapGoSwaggerToGoCSQLFieldType %s %s\n ", fieldName,fieldType )}

	fieldName = suffix + fieldName
	switch strings.ToLower(fieldType) {
	case "int":
		ret = ret + INDENT_1 + "int(" + fieldName + "),"
	case "timestamp":
		ret = ret + INDENT_1 + PARSERTIME_FUNC_NAME + "(" + fieldName + "),"
	case "float":
		ret = ret +  INDENT_1 + "float32(" + fieldName + "),"
	default:
		ret = ret + INDENT_1 + fieldName + ","
	}

	if debug { fmt.Printf("mapGoSwaggerToGoCSQLFieldType returning %s from field %s type %s\n", ret, fieldName, fieldType ) }
	return ret
}


// Function called to convert from go-swagger type tp go-cql type
func copyFromStructToStruc( debug bool, suffix string, dest string, typeDetails * parser.TypeDetails, parserOutput parser.ParseOutput  ) string {
	ret := ""

	if debug {fmt.Printf("copyFromStructToStruc %s %s\n ", suffix, dest )}

	for i := 0; i <  typeDetails.TypeFields.FieldIndex; i++ {
		f := typeDetails.TypeFields.DbFieldDetails[i]
		if swagger.IsFieldTypeUDT(parserOutput, f.DbFieldType) {
			if debug { fmt.Println("copyFromStructToStruc Found UDT = ", f.DbFieldType) }
			// Process UDT
			//fieldName := GetFieldName(debug, false, f.OrigFieldName, false)
			//fieldType := GetFieldName(debug, false, f.DbFieldType, true)
			// Need to recurse here
		} else {
			ret = applyTypeConversionForGoSwaggerToGocql( debug , ret , suffix, CapitaliseSplitFieldName ( debug , strings.ToLower(f.DbFieldName) , false),  f.DbFieldType)
		}
	}

	return ret
}


func processPostField(debug bool, fieldName string,  parserOutput parser.ParseOutput, fieldDetails parser.FieldDetails ) string {
    ret := ""
	if debug {fmt.Printf("processPostField %s %s\n ", fieldName, fieldDetails.DbFieldType )}
	switch strings.ToLower(fieldDetails.DbFieldType) {
	case "timestamp":
		/*
		tmp := createTempVar( fieldName )
		field := GetFieldName(  debug , false, fieldName, false )
		ret = INDENT_1 + "if " + "params.Body." + field + ` != "" { `
		ret = ret + INDENT_1 + INDENT2 + tmp + ",ok" + tmp + " := time.Parse( time.RFC3339,params.Body." + field + ")"
		ret = ret + INDENT_1 + INDENT2 + "if ok" + tmp + " != nil {" + `
` + INDENT2 + INDENT3 +  "log.Println(" + "ok" + tmp + `)
` + INDENT2 + INDENT3 +  `m["` + fieldName + `"] = ""` + INDENT_1 + INDENT2 + "}"
		ret = ret + " else { " + INDENT_1 + INDENT3 + `m["` + fieldName + `"] = ` + tmp + INDENT_1 + INDENT2 +  "}"
		ret = ret + INDENT_1 + "}" + " else {" +  INDENT_1 + INDENT2 +  `m["` + fieldName + `"] = ""` + INDENT_1 +  "}"
		*/
		ret = ret + INDENT_1 + `m["` + fieldName + `"] = ` + PARSERTIME_FUNC_NAME + "(" + "params.Body." + GetFieldName(debug, false, fieldName, false) + ")"
	case "date":
		field := GetFieldName(  debug , false, fieldName, false )
		ret = ret + INDENT_1 + `m["` + fieldName + `"] = ` + "params.Body." + GetFieldName(  debug , false, fieldName, false )
		ret = ret + INDENT_1 + "if " + "params.Body." + field + ` == "" { `
		ret = ret + INDENT_1 + INDENT2 + `m["` + fieldName + `"] =  "1970-01-01"` + INDENT_1 + "}"
	case "float":
		tmp := createTempVar( fieldName )
		tmp1 := createTempVar( fieldName )
		field := GetFieldName(  debug , false, fieldName, false )
		ret = ret + INDENT_1 + tmp + `:= fmt.Sprintf("%f",params.Body.` + field + ")"
		ret = ret + INDENT_1 + tmp1 + `,_ := strconv.ParseFloat(` + tmp + ",32)"
		ret = ret + INDENT_1 + `m["` + fieldName + `"] = float32(` + tmp1 + ")"
	default:
		if swagger.IsFieldTypeUDT( parserOutput, fieldDetails.DbFieldType) {
			fieldName := CapitaliseSplitFieldName(debug, strings.ToLower(fieldDetails.DbFieldName),false)
			dest := createTempVar( fieldName )
			suffix := "params.Body." + fieldName + "."
			ret = INDENT_1 + dest + " := &" + CapitaliseSplitFieldName(debug, strings.ToLower(fieldDetails.DbFieldType),false) + "{"
			//source := "params.Body" + "." + fieldName
			t := copyFromStructToStruc(debug, suffix, fieldName,findTypeDetails (debug ,fieldDetails.DbFieldType, parserOutput), parserOutput)
			ret = ret + t
			ret = ret + INDENT_1 + "}"
			ret = ret + INDENT_1 + `m["` + strings.ToLower(fieldName) + `"] = ` + "&" + dest
		} else {
			ret = ret + INDENT_1 + `m["` + fieldName + `"] = ` + "params.Body." + GetFieldName(debug, false, fieldName, false)
		}
	}
    return ret
}




