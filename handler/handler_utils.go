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
func CreateFile( debug bool, pathPrefix string, dir string ) *os.File {

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
	fullFileName := fulldirName + "/" + MAINFILE
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
	name := ""
	if leaveCase {
		name = fieldName
	} else {
		name = strings.ToLower(fieldName)
	}
	return CapitaliseSplitFieldName( debug, name, dontUpdate )
}



// Function that renames fields to match that performed for some reason by go-swagger in its generated framework code
func Capitiseid( debug bool, fieldName string, dontUpdate bool ) string {

	var ret string = ""
	if debug { fmt.Printf("Capitiseid entry field  = %s, len = %d\n ",fieldName, len(fieldName) ) }

	if dontUpdate {
		ret = fieldName
		if debug { fmt.Printf("Capitiseid told not to update\n ") }
	} else {
		runes := []rune(fieldName[:])
		last := len( runes  ) - 1

		for i := 0; i < last ; i++   {
				//if debug { fmt.Printf("Capitiseid i   = %d\n ",i ) }
				if ! ( ( i == 0 ) || ( i == last -1 ) ) {
					continue;
				}
				//if debug { fmt.Printf("Capitiseid [0] = %q, [1] = %q\n ", runes[i], runes[i+1]) }
				if (runes[i] == rune('i') || runes[i] == rune('I')) && (runes[i+1] == rune('d') || runes[i+1] == rune('D')) {
					if debug { fmt.Printf("Capitiseid match at i= %d\n ", i)}
					runes[i] = rune('I')
					runes[i+1] = rune('D')
				}
			}
			ret = string(runes)
		}

	if debug {fmt.Printf("Capitiseid returning field  = %s\n ", ret)}
	return ret
}

// Function that renames fields to match that performed for some reason by go-swagger in its generated framework code e.g. My_List becomes MyList & address_id becomes AddressID
func CapitaliseSplitFieldName ( debug bool, fieldName string, dontUpdate bool ) string {
//@TODO remove
debug = false
	var ret string = ""
	if debug { fmt.Printf("CapitaliseSplitFieldName entry field  = %s, len = %d\n ",fieldName, len(fieldName) ) }

	if dontUpdate  || fieldName == ""{
		ret = fieldName
		if debug { fmt.Printf("CapitaliseSplitFieldName told not to update\n ") }
	} else {
		tmpFields := strings.Split(fieldName, "_" )
		if debug {fmt.Printf("CapitaliseSplitFieldName tmpFields  = %q\n ", tmpFields)}
		for _, v := range tmpFields {
			v = Capitiseid( debug, v, dontUpdate )
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
			text = "int32"
		} else {
			text = "int64"
		}
	case "uuid":
		text = "string"
	case "date": fallthrough
	case "timeuuid": fallthrough
	case "timestamp":
			text = "time.Time"
	case "boolean":
		text = "bool"
	case "decimal":
		text = "*inf.Dec" // this is in the gopkg.in/inf.v0 package
	case "float": fallthrough
	case "double":
		if makeSmall {
			text = "float32"
		} else {
			text = "float64"
		}
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
			text =  MODELS + typeName + fieldName
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
		text = "int32"
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
			text = "[]*"
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



func CopyArrayElements( debug bool, inTable bool, inDent string, sourceFieldName string, destFieldName string,  fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, dontUpdate bool  ) string {
	equals := " := "
	/*if inTable {
		equals = " = "
	}
	*/
	arrayType := basicMapCassandraTypeToGoType(debug, false, inTable, fieldDetails.DbFieldName, fieldDetails.DbFieldCollectionType, "", fieldDetails, parserOutput, dontUpdate, true )
	arrayTypeDest := basicMapCassandraTypeToGoType(debug, false, inTable, fieldDetails.DbFieldName, fieldDetails.DbFieldCollectionType, "", fieldDetails, parserOutput, dontUpdate, false )
	ret := INDENT_1 + inDent + sourceFieldName + equals +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].([]` + arrayType + ")"
	ret = ret + inDent + destFieldName + " = " + "make([] " + arrayTypeDest + ", len(" + sourceFieldName + ") )"
	ret = ret + inDent + "for j := 0; j < len(" + sourceFieldName + " ); j++ { "
	switch arrayTypeDest {
	case "float64":
		ret = ret + inDent + INDENT + destFieldName + "[j] = " +  "float64(" + sourceFieldName + "[j])" + inDent + "}"
	case "int64":
		ret = ret + inDent + INDENT + destFieldName + "[j] = " +  "int64(" + sourceFieldName + "[j])" + inDent + "}"
	default:
		if debug {fmt.Printf("CopyArrayElements TYPE NOT MATCHED!!!!\n " )}
	}
	return ret
}


func copyStruct( debug bool, inDent string, recursing bool,  sourceStruct string, sourceField string, destStruct string ,typeDetails *parser.TypeDetails, dontUpdate bool  ) string  {
	typeName := GetFieldName(debug, recursing, typeDetails.TypeName, dontUpdate )
	ret := INDENT_1 + inDent + destStruct + " = &" + typeName + "{"

	for i := 0; i < typeDetails.TypeFields.FieldIndex; i++ {
		if i > 0 {
			ret = ret + ","
		}
		//fieldType := basicMapCassandraTypeToGoType(debug, false, inTable, typeDetails.TypeFields.DbFieldDetails[i].DbFieldName, typeDetails.TypeFields.DbFieldDetails[i].DbFieldCollectionType, typeDetails.TypeName, typeDetails.TypeFields.DbFieldDetails[i], parserOutput, dontUpdate, false )
		fieldName := GetFieldName(debug, recursing, typeDetails.TypeFields.DbFieldDetails[i].DbFieldName, dontUpdate )
		ret = ret + inDent + fieldName + ":" + sourceStruct + "." + sourceField + "." + fieldName
	}
	ret = ret + "}"
	return ret
}

