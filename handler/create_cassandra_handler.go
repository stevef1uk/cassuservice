package handler

import (
	"fmt"
	"github.com/stevef1uk/cassuservice/swagger"
	"log"
	"strconv"
	//"fmt"
	"github.com/stevef1uk/cassuservice/parser"
	"os"
	"strings"

	//"strings"
)



func WriteHeaderPart( debug bool, parserOutput parser.ParseOutput, generateDir string, endPointNameOverRide string, dontUpdate bool, addPost bool, output  *os.File ) bool  {
	doNeedTimeImports := doINeedTime(parserOutput )
	needDecimalImports := doINeedDecimal(parserOutput )

	extraImports := ""
	if doNeedTimeImports {
		extraImports = extraImports + IMPORTSTIMESTAMP
	}
	if needDecimalImports {
		extraImports = extraImports + IMPORTDEC
	}
	if addPost {
		extraImports = extraImports + IMPORTFORPOST
		if doIHaveFloat( parserOutput ) {
			if ( needDecimalImports ) {
				extraImports = extraImports + IMPORTFORPOST2A
			} else {
				extraImports = extraImports + IMPORTFORPOST2
			}
		}
	}

	tmpData := &tableDetails{ generateDir, "", strings.ToLower(parserOutput.TableDetails.TableName),""}
	WriteStringToFileWithTemplates(  COMMONIMPORTS + extraImports + IMPORTSEND , "codegen-get", output, &tmpData)

	return doNeedTimeImports
}


// Write out the types for the UDT. Returns true if a Post specific structure has been added
func addStruct( debug bool, addPost bool, parserOutput parser.ParseOutput, output  *os.File ) bool {
	ret := false

	for i := 0; i < parserOutput.TypeIndex; i++ {
		localModelType := ""
		v := parserOutput.TypeDetails[i]
		output.WriteString( "\ntype " + GetFieldName( debug, false, v.TypeName, true)  + " struct {")
		needToAddPostStruct := false
		for j := 0; j < v.TypeFields.FieldIndex ; j++ {
			revisedFieldName := GetFieldName(debug, false, v.TypeFields.DbFieldDetails[j].OrigFieldName , false )
			revisedType := GetFieldName( debug, false, v.TypeName,true)
			tmp := mapFieldTypeToGoCSQLType( debug,  revisedFieldName, true,false, v.TypeFields.DbFieldDetails[j].DbFieldType, revisedType, v.TypeFields.DbFieldDetails[j], parserOutput, true  )
			//if swagger.IsFieldTypeUDT(parserOutput,  )
			//debug bool, fieldName string, leaveFieldCase bool, inTable bool, fieldType string, typeName string, fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, dontUpdate bool
			output.WriteString( "\n    " + revisedFieldName + " ")
			output.WriteString( tmp   + " `" + `cql:"` + strings.ToLower( v.TypeFields.DbFieldDetails[j].DbFieldName ) + `"` +"`")

			if addPost && strings.HasPrefix( tmp, MODELS )  { // Need to define a local structure for Posts to use cql annoated type
				ret = true
				needToAddPostStruct = ret
				localType := ""
				if v.TypeFields.DbFieldDetails[j].DbFieldMapType != "" {
					localType = " map[string]string"
				} else  if v.TypeFields.DbFieldDetails[j].DbFieldCollectionType != "" {
					localType = " []* " +  GetFieldName( debug, false, v.TypeFields.DbFieldDetails[j].DbFieldCollectionType,true)
				}
				//localTypeName := GetFieldName( debug, false, v.TypeName, true)
				localModelType = localModelType + "\ntype " + strings.TrimPrefix( tmp, MODELS ) +localType
			}
		}

		if needToAddPostStruct {
			output.WriteString("\n}\n" )
			output.WriteString("\n" + localModelType + "\n")
			output.WriteString( "\ntype " + GetFieldName( debug, false, v.TypeName, true)  + "Post struct {")
			for j := 0; j < v.TypeFields.FieldIndex ; j++ {
				revisedFieldName := GetFieldName(debug, false, v.TypeFields.DbFieldDetails[j].OrigFieldName , false )
				revisedType := GetFieldName( debug, false, v.TypeName,true)
				tmp := mapFieldTypeToGoCSQLType( debug,  revisedFieldName, true,false, v.TypeFields.DbFieldDetails[j].DbFieldType, revisedType, v.TypeFields.DbFieldDetails[j], parserOutput, true  )
				if strings.HasPrefix( tmp, MODELS )  {
					tmp = strings.TrimPrefix( tmp, MODELS )

				}
					//debug bool, fieldName string, leaveFieldCase bool, inTable bool, fieldType string, typeName string, fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, dontUpdate bool
				output.WriteString( "\n    " + revisedFieldName + " ")
				output.WriteString( tmp   + " `" + `cql:"` + strings.ToLower( v.TypeFields.DbFieldDetails[j].DbFieldName ) + `"` +"`")

			}
		}
		output.WriteString("\n}\n" )
	}
    return ret
}


func writeField( debug bool, inTable bool, parserOutput parser.ParseOutput, field parser.FieldDetails, output  *os.File) {

	fieldName := GetFieldName( debug, false, field.OrigFieldName, false)

	if field.DbFieldCollectionType != "" {
		collectionType := GetFieldName(debug, false, field.DbFieldCollectionType, true )
		fieldType :=  mapTableTypeToGoType( debug, inTable, fieldName, collectionType, field.DbFieldCollectionType, field, parserOutput, true )
		if debug {fmt.Println("writeField name =", field.DbFieldName, " fieldType = ", fieldType) }
		if strings.ToLower(fieldType ) == "map" {
			output.WriteString( INDENT_1 + "var " + fieldName + " " +  fieldType )
		} else {
			if swagger.IsFieldTypeUDT( parserOutput, field.DbFieldCollectionType ) {
				output.WriteString( INDENT_1 + "var " + fieldName + " " + fieldType )
			} else {
				output.WriteString( INDENT_1 + "var " + fieldName + " []" + fieldType )
			}
		}
	} else {
		fieldType :=  mapTableTypeToGoType( debug, inTable, strings.ToLower(field.DbFieldName), field.DbFieldType, field.DbFieldCollectionType, field, parserOutput, true )
		if debug {fmt.Println("writeField name =", fieldName, " fieldType = ", fieldType) }
		output.WriteString( INDENT_1 + "var " + fieldName + " " + fieldType )
	}
	// Not all vars used so handle any compilation errors
	output.WriteString( INDENT_1 + "_ = " + fieldName  )
}


// Function that writes out the variable types for the table & returns the temporary variable created if there is a time field
func WriteVars(  debug bool, inTable bool, parserOutput parser.ParseOutput, goPathForRepo string, doNeedTimeImports bool, addPost bool, endPointNameOverRide string, output  *os.File )  string {
	tmpTimeVar := ""

	const UDTTYPE = `
    {{.FieldName}} := &{{.TypeName}}{}
    _= {{.FieldName}}
`

	for i := 0; i < parserOutput.TableDetails.TableFields.FieldIndex; i++ {
		v := parserOutput.TableDetails.TableFields.DbFieldDetails[i]
		if debug {fmt.Println("WriteVars v =", v.DbFieldName) }
		// If field type is a UDT
		if swagger.IsFieldTypeUDT(parserOutput, v.DbFieldType) {
			if debug {fmt.Println("WriteVars Found UDT = ", v.DbFieldType)}
			// Process UDT
			fieldName := GetFieldName(debug, false, v.OrigFieldName, false)
			fieldType := GetFieldName(debug, false, v.DbFieldType, true)
			output.WriteString( INDENT_1 + fieldName + " := &" + fieldType  + "{}" + INDENT_1 + "_ = " +  fieldName )

		} else {
			if debug {fmt.Println("WriteVars writing field") }
			writeField( debug, inTable, parserOutput, v, output)
		}
	}

	if doNeedTimeImports {
		tmpTimeVar = createTempVar( TMP_TIME_VAR_PREFIX )
		//output.WriteString( INDENT_1 + tmpTimeVar + " := strfmt.NewDateTime().String()" )
	}

	// Handle cases where we donn't access any structures in the model direwctory
	if addPost {
		output.WriteString(INDENT_1 + "_ = " + MODELS + swagger.CapitaliseSplitTableName(debug, parserOutput.TableDetails.TableName) + "{}")
	}

	output.WriteString( "\n" + INDENT_1 + SELECT_OUTPUT + " := map[string]interface{}{}\n")


	return tmpTimeVar
}


func buildSelectParams ( debug bool, parserOutput parser.ParseOutput )  string {
	ret := ""
	for i :=0; i < parserOutput.TableDetails.TableFields.FieldIndex; i++ {
		if i > 0 {
			ret = ret + ", "
		}
		ret = ret + strings.ToLower( parserOutput.TableDetails.TableFields.DbFieldDetails[i].DbFieldName )
	}
	return ret
}


func createSelectString( debug bool, parserOutput parser.ParseOutput, tableName string,  timeVar string, cassandraConsistencyRequired string,  overridePrimaryKeys int, allowFiltering bool, logExtraInfo bool, output  *os.File )  string {

	ret := buildSelectParams( debug, parserOutput )
	// First build primary key conditions
	pkNum := parserOutput.TableDetails.PkIndex
	if overridePrimaryKeys != 0 {
		if debug {fmt.Println("createSelectString primary fields constrained for select statement") }
		pkNum = overridePrimaryKeys
	}
	whereClause := " WHERE "
	varsClause := ""
	for i:= 0; i < pkNum; i++ {
		v := swagger.FindFieldByname( parserOutput.TableDetails.DbPKFields[i], parserOutput.TableDetails.TableFields.FieldIndex, parserOutput.TableDetails.TableFields )
		if i > 0 {
			whereClause = whereClause + "and "
			varsClause = varsClause + ","
		}
		whereClause = whereClause + strings.ToLower( v.DbFieldName ) + " = ? "

		fieldName := GetFieldName(debug, false, v.OrigFieldName, false)
		if swagger.IsFieldTypeATime( v.DbFieldType) {
			// Need to parse the received parameter
			output.WriteString( INDENT_1 + fieldName + ",_ = time.Parse(time.RFC3339,params." + fieldName + ".String() ) ")
			varsClause = varsClause + fieldName
		} else {
			varsClause = varsClause + "params." + fieldName
		}
	}

	consistency := "gocql.One"
	if cassandraConsistencyRequired != "" {
		consistency = cassandraConsistencyRequired
	}

	ret = ret + " FROM " + strings.ToLower( parserOutput.TableDetails.TableName ) + whereClause
	if allowFiltering {
		ret = ret + "ALLOW FILTERING"
	}
	ret = ret +   "`," + varsClause + ")"
	ret = ret + ".Consistency(" + consistency + ").MapScan(codeGenRawTableResult); err != nil {"
	if logExtraInfo {
		ret = ret + INDENT_1 + `  log.Println("No data? ", err)`
	}
	tableName = GetFieldName(debug, false, tableName, false)
	ret = ret + INDENT_1 + "  return " + OPERATIONS + "NewGet" + tableName + "BadRequest()" +  INDENT_1 + "}"
	ret = ret + INDENT_1 + PAYLOAD + " := " + OPERATIONS + "NewGet" +  tableName + "OK()"
	ret = ret + INDENT_1 + PAYLOAD + "." + PAYLOAD_STRUCT + " = make([]*" + OPERATIONS + "Get" + tableName + "OKBodyItems0,1)"
	ret = ret + INDENT_1 + PAYLOAD + "." + PAYLOAD_STRUCT + "[0] = new(" + OPERATIONS + "Get" + tableName + "OKBodyItems0)"
	ret = ret + INDENT_1 + PARAMS_RET + " := " + PAYLOAD + "." + PAYLOAD_STRUCT + "[0]"

	return ret
}

// Function called to process a local UDT structure and copy into the go-swagger model's structure type for the UDT
func handleStructVarConversion(  debug bool, recursing bool, inDent string, theStructVar string, destVar string, fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput ) string {

	ret := ""

	fieldName := GetFieldName( debug, false, fieldDetails.OrigFieldName, false)
	sourceVar := theStructVar + "." +  fieldName

	switch ( strings.ToLower( fieldDetails.DbFieldType ) ) {
	case "int":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := int64(" + sourceVar + ")"
		ret = ret + INDENT_1 + inDent + destVar + " = " + tmp_var
	case "float":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := float64(" + sourceVar + ")"
		ret = ret + INDENT_1 + inDent + destVar + " = " + tmp_var
	case "date": fallthrough
	case "timestamp": fallthrough
	case "timeuuid":
		tmp_var := createTempVar( fieldName )
		//tmp, tmp1 := ProcessTime( timeFound, INDENT_1 + inDent, timeVar, theStructVar,  fieldName )
		//ret = ret + tmp
		//ret = ret  + INDENT_1 + inDent + destVar + " = " + tmp1
		ret = ret + INDENT_1 + inDent + tmp_var + " := " + sourceVar + ".String()"
		ret = ret + INDENT_1 + inDent + destVar + " = " + tmp_var
	case "set": fallthrough //
	case "list":
		collectionType := GetFieldName(debug, recursing, fieldDetails.DbFieldCollectionType, true )
		if swagger.IsFieldTypeUDT( parserOutput, collectionType ) {
			ret = INDENT_1 + inDent + destVar + " = " + sourceVar
		} else {
			ret = CopyArrayElements( debug, false,false, INDENT_1 + INDENT, theStructVar + "." + fieldName, destVar,  fieldDetails, parserOutput )
		}
	default:
		if swagger.IsFieldTypeUDT( parserOutput, fieldDetails.DbFieldType ) {
			//tmp := handleStructVarConversion(  debug, true, inDent, theStructVar, destVar, fieldDetails, parserOutput )
			//_ = tmp
			typeDetails := findTypeDetails ( debug , fieldDetails.DbFieldType , parserOutput )
			tmpTypeName := MODELS + GetFieldName( debug, false, fieldDetails.DbFieldType, false)
			ret = ret + INDENT_1 + inDent + INDENT + destVar + " = &" + tmpTypeName + "{}"
			ret = ret + convertToModelType( debug, inDent + INDENT , false, theStructVar, destVar, typeDetails, parserOutput )

		} else {
			ret = INDENT_1 + inDent + destVar + " = " + sourceVar
		}

	}
	return ret
}


var indexCounter int = 0

// Handle the case of a single UDT field, which can only occur in a Table definition right now. Sadly not true as they can be in UDTs too!
func setUpStruct ( debug bool, recursing bool,  timeFound bool, inDent string, inTable bool, raw_data string, destField string, theVar string, theType string,  parserOutput parser.ParseOutput, timeVar string ) (string,string) {

	extraVars := ""
	newStr := ""
	ret := ""

	typeStruct := findTypeDetails( debug, theType, parserOutput )
	var structAssignment []bool = make( []bool, typeStruct.TypeFields.FieldIndex)
	structName := GetFieldName(  debug, recursing, theType, false )
	tmpStruct := createTempVar( structName )

	if recursing {
		inDent = inDent + INDENT2
	}

	loopAssignment := INDENT_1 + inDent  + tmpStruct + " := &" + structName + "{"
	if ! inTable {
		loopAssignment = INDENT_1 + inDent  + tmpStruct + " := " + structName + "{"
	}

	for i := 0; i < typeStruct.TypeFields.FieldIndex; i++ {
		fieldName := strings.ToLower( typeStruct.TypeFields.DbFieldDetails[i].DbFieldName )
		fieldType := mapFieldTypeToGoCSQLType( debug, fieldName, true, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldType, structName, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, true  )
		if recursing {
			inDent = inDent + INDENT
		}

		//if ( typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType != "" && ( swagger.IsFieldTypeUDT( parserOutput, typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType ) ) )  || typeStruct.TypeFields.DbFieldDetails[i].DbFieldMapType != ""  {
		if ( typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType != "" && ( swagger.IsFieldTypeUDT( parserOutput, typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType ) ) )  ||
			typeStruct.TypeFields.DbFieldDetails[i].DbFieldMapType != "" || ( swagger.IsFieldTypeUDT( parserOutput, typeStruct.TypeFields.DbFieldDetails[i].DbFieldType ) ) {
			// Deal with the more complex types
			tmpVar := createTempVar( fieldName )
			tmpVar1 := createTempVar( fieldName )
			// Note as there seems no way of mapping a Map type in Swagger to anything other than string:string we are a bit stuffed here!
			if typeStruct.TypeFields.DbFieldDetails[i].DbFieldMapType != "" {
				ret = ret + INDENT_1 + inDent + INDENT2  + tmpVar + ","
				extraVars = extraVars +  INDENT_1 + inDent + tmpVar + " := " + raw_data +  `["` + strings.ToLower(fieldName ) + `"].(map[string]string)`
			} else {
				if  typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType != "" {
					// Handle lists & sets here!
					ret = ret + INDENT_1 + inDent + INDENT2 + tmpVar1 + ","
					extraVars = extraVars + INDENT_1 + inDent + tmpVar + ":= " + raw_data + `["` + strings.ToLower(fieldName) + `"].([]map[string]interface{})`
					tmpType := mapFieldTypeToGoCSQLType(debug, fieldName, true, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType, structName, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, true)
					extraVars = extraVars + INDENT_1 + inDent + tmpVar1 + ":= make(" + tmpType + ", len(" + tmpVar + ") )"
					extraVars = extraVars + INDENT_1 + inDent + setUpStructs(debug, true, timeFound, inDent, false, tmpVar1, tmpVar, strings.ToLower(typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType), parserOutput, timeVar)
					structAssignment[i] = true
				} else { // Case of single UTD at this point. We can't set this up, that needs to be done below
					typeName := GetFieldName(  debug, recursing, typeStruct.TypeFields.DbFieldDetails[i].DbFieldType, false )
					ret = ret + INDENT_1 + inDent + INDENT2 + typeName + "{},"
				}
			}
		} else {
			ret = ret + INDENT_1 + inDent + INDENT2  + raw_data + `["` + strings.ToLower(fieldName ) + `"].(` + fieldType + "),"
		}
	}
	ret = ret + INDENT_1 + inDent + INDENT + "}"
	ret =  INDENT_1 + extraVars + loopAssignment + INDENT_1 + ret + INDENT_1 + newStr

	// Now process each variable in order to set-up the Payload structure
	tmp := ""
	for i := 0; i < typeStruct.TypeFields.FieldIndex; i++ {
		fieldName := GetFieldName(debug, false, typeStruct.TypeFields.DbFieldDetails[i].OrigFieldName, false)
		tmpDest := destField + "." + fieldName
		if structAssignment[i] {
			if typeStruct.TypeFields.DbFieldDetails[i].DbFieldMapType == "" && ! inTable {
				tmpDest = destField
				tmp = INDENT_1 + inDent + tmpDest + " = " + tmpStruct
				goto here
			}
		} else if swagger.IsFieldTypeUDT( parserOutput, typeStruct.TypeFields.DbFieldDetails[i].DbFieldType ) {
			// Case of a UTD within a UTD
			tmpVar := createTempVar( fieldName )
			ret = ret + INDENT_1 + inDent + tmpVar + ",ok := " + raw_data + `["` + strings.ToLower( fieldName ) + `"].(map[string]interface{})`
			ret = ret + INDENT_1 + inDent +  "if ! ok {" + INDENT_1 + INDENT + `log.Fatal("handleReturnedVar() - failed to find entry for ` +  fieldName + `", ok )` + INDENT_1 + "}"
			typeName := GetFieldName(  debug, recursing, typeStruct.TypeFields.DbFieldDetails[i].DbFieldType, false )
			ret = ret + INDENT_1 + inDent + destField + "." + fieldName + " = &" + "models." + typeName + "{}"
			destField = destField + "." + GetFieldName(debug, false, typeStruct.TypeFields.DbFieldDetails[i].OrigFieldName, false)
			tmpStruct, _ := setUpStruct(debug , recursing ,  timeFound , inDent , false, tmpVar, destField, theVar, typeName , parserOutput, timeVar )
			ret = ret + INDENT_1 + tmpStruct
			return ret, tmpStruct
		}

		tmp = handleStructVarConversion(debug, recursing, inDent, tmpStruct, tmpDest, typeStruct.TypeFields.DbFieldDetails[i], parserOutput)
	here:
		ret = ret + inDent + INDENT2  + tmp
	}

	return ret, tmpStruct
}


// Returns the value type for maps as a string and a boolean indicating it it is a UTD
func getMapType( debug bool, recursing bool, inTable bool, mapType string, fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput  ) (string,bool) {
	mapTypeInGo := ""
	isUDT := false
	if swagger.IsFieldTypeUDT( parserOutput, mapType  ) { // Map is <text, frozen <UDT>> where text has to be string owing to swagger limitation
		mapTypeInGo = GetFieldName(debug, recursing, mapType, false )
		isUDT = true
	} else {
		mapTypeInGo = basicMapCassandraTypeToGoType( debug , recursing , inTable, mapType , mapType , mapType,  fieldDetails, parserOutput, false , true  )
	}
	return mapTypeInGo, isUDT
}


// Function to handle UDTs in collection types
func setUpStructs ( debug bool, recursing bool,  timeFound bool, inDent string, inTable bool, destField string, theVar string, theType string,  parserOutput parser.ParseOutput, timeVar string ) string {

	indexCounter++
	extraVars := ""
	newStr := ""
	ret := ""
	typeStruct := findTypeDetails( debug, theType, parserOutput )
	var structAssignment []bool = make( []bool, typeStruct.TypeFields.FieldIndex)
	structName := GetFieldName(  debug, recursing, theType, false )
	tmpStruct := createTempVar( structName )
	iIndex := "i" + strconv.Itoa(indexCounter)
	vIndex := "v" + strconv.Itoa(indexCounter)

	if recursing {
		inDent = inDent + INDENT2
	}
	space := ""
	if recursing {
		space = inDent
	}
	loopStart := INDENT_1  + space + "for " + iIndex + ", " + vIndex + " := range " + theVar + " {"
	beforeLoop := ""
	loopAssignment := INDENT_1 + inDent  + tmpStruct + " := &" + structName + "{"
	for i := 0; i < typeStruct.TypeFields.FieldIndex; i++ {
		fieldName := strings.ToLower( typeStruct.TypeFields.DbFieldDetails[i].DbFieldName )
		fieldType := mapFieldTypeToGoCSQLType( debug, fieldName, true, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldType, structName, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, true  )

		if swagger.IsFieldTypeUDT( parserOutput, typeStruct.TypeFields.DbFieldDetails[i].DbFieldType ) {
			tmpVar := createTempVar( fieldName )
			//tmpVar1 := createTempVar( fieldName )
			tmpModelVar1 := createTempVar( fieldName )
			beforeLoop = beforeLoop + INDENT_1 + inDent  + tmpVar + " := " + vIndex + `["` + fieldName + `"].(map[string]interface{})`
			tmpTypeName := GetFieldName(debug, recursing, strings.ToLower(typeStruct.TypeFields.DbFieldDetails[i].DbFieldType) , false )
			beforeLoop = beforeLoop + INDENT_1 + inDent   + tmpModelVar1 + " := &" + MODELS + tmpTypeName + "{}"
			//beforeLoop = beforeLoop + INDENT_1 + inDent + INDENT2  + tmpVar1 + " := &" + tmpTypeName + "{}"
			tmptmp, tmpVar2 := setUpStruct ( debug , recursing ,  timeFound , inDent , false , tmpVar, tmpModelVar1, "SJFSJF", tmpTypeName ,  parserOutput, timeVar )
			_ = tmpVar2
			beforeLoop = beforeLoop + INDENT_1 + inDent + INDENT2 + tmptmp
			ret = ret + INDENT_1 + inDent + INDENT2  + tmpVar2 + ","
			continue
		}

		if recursing {
			inDent = inDent + INDENT
		}
		if ( typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType != "" && ( swagger.IsFieldTypeUDT( parserOutput, typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType ) ) )  || typeStruct.TypeFields.DbFieldDetails[i].DbFieldMapType != ""  {
			// Deal with the more complex types
			tmpVar := createTempVar( fieldName )
			tmpVar1 := createTempVar( fieldName )
			// Note as there seems no way of mapping a Map type in Swagger to anything other than string:string we are a bit stuffed here!
			if typeStruct.TypeFields.DbFieldDetails[i].DbFieldMapType != "" {
				ret = ret + INDENT_1 + inDent + INDENT2  + tmpVar + ","
				mapTypeInGo,_ := getMapType(  debug , recursing , inTable , typeStruct.TypeFields.DbFieldDetails[i].DbFieldMapType , typeStruct.TypeFields.DbFieldDetails[i] , parserOutput )
				extraVars = extraVars +  INDENT_1 + inDent + tmpVar + " := " +vIndex +  `["` + strings.ToLower(fieldName ) + `"].(map[string]` + mapTypeInGo + ")"
			} else {
				// Handle lists & sets here!
				ret = ret + INDENT_1 + inDent + INDENT2  + tmpVar1 + ","
				extraVars = extraVars + INDENT_1 + inDent + "if " + vIndex + `["` + strings.ToLower(fieldName) + `"] == nil { ` +  INDENT_1 + inDent + INDENT2 + "continue" + INDENT_1 + inDent  + "}"
				extraVars = extraVars +  INDENT_1 + inDent + tmpVar + ":= " + vIndex + `["` + strings.ToLower(fieldName ) + `"].([]map[string]interface{})`
				tmpType := mapFieldTypeToGoCSQLType( debug, fieldName, true, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType, structName, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, true  )
				extraVars = extraVars +  INDENT_1 + inDent + tmpVar1 + ":= make(" + tmpType + ", len(" + tmpVar + ") )"
				extraVars = extraVars + INDENT_1 + inDent + setUpStructs( debug,  true,  timeFound, inDent, false, tmpVar1,  tmpVar, strings.ToLower(typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType),  parserOutput, timeVar )
				structAssignment[i] = true
			}
		} else {
			//ret = ret + INDENT_1 + inDent +  INDENT2 + "if " + vIndex + `["` + strings.ToLower(fieldName) + `"] == nil { ` +  INDENT_1 + inDent + INDENT2 + "continue" + INDENT_1 + inDent  + "}"
			ret = ret + INDENT_1 + inDent + INDENT2  + vIndex + `["` + strings.ToLower(fieldName ) + `"].(` + fieldType + "),"
		}
	}
	ret = ret + INDENT_1 + inDent + INDENT + "}"
	ret = loopStart + INDENT_1 + extraVars + beforeLoop + loopAssignment + INDENT_1 + ret + INDENT_1 + newStr

	// Now process each variable in order to set-up the Payload structure
	tmp := ""
	for i := 0; i < typeStruct.TypeFields.FieldIndex; i++ {
		fieldName := GetFieldName(debug, false, typeStruct.TypeFields.DbFieldDetails[i].OrigFieldName, false)
		typeName := GetFieldName(debug, false, typeStruct.TypeName, false)
		tmpDest := destField + "[" + iIndex + "]." + fieldName
		addMake := INDENT_1 + inDent + destField + "[" + iIndex + "] = &" + MODELS + typeName + "{}"
		if structAssignment[i] {
			if typeStruct.TypeFields.DbFieldDetails[i].DbFieldMapType == "" && ! inTable {
				tmpDest = destField + "[" + iIndex + "]"
				tmp = INDENT_1 + inDent + tmpDest + " = " + tmpStruct
				goto here
			}
		}
		tmp = handleStructVarConversion(debug, recursing, inDent, tmpStruct, tmpDest, typeStruct.TypeFields.DbFieldDetails[i], parserOutput )
	here:
		if i == 0 {
			ret = ret + inDent + INDENT2 + addMake + tmp
		} else {
			ret = ret + inDent + INDENT2  + tmp
		}
	}

	ret = ret + INDENT_1 + inDent + "}"

	return ret
}



func handleReturnedVar( debug bool, recursing bool, timeFound bool, inDent string, inTable bool, typeIndex int , fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, timeVar string ) (string, bool) {
	indexCounter++
    ret := ""
	fieldName := GetFieldName( debug, false, fieldDetails.OrigFieldName, false)
	switch ( strings.ToLower( fieldDetails.DbFieldType ) ) {
	case "boolean":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := " +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(bool)`
		ret = ret + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = &" + tmp_var
	case "int":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := " +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(int)`
		ret = ret + INDENT_1 + inDent + fieldName + " = int64(" + tmp_var + ")"
		ret = ret + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = &" + fieldName
	case "bigint":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := " +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(int64)`
		ret = ret + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = &" + tmp_var
	case "float":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := " +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(float32)`
		ret = ret + INDENT_1 + inDent + fieldName + " = float64(" + tmp_var + ")"
		ret = ret + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = &" + fieldName
	case "decimal":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + fieldName + " = " +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(*inf.Dec)`
		ret = ret + INDENT_1 + inDent + tmp_var + ",_ := " + "strconv.ParseFloat(" + fieldName + ".String(), 64 )"
		ret = ret + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = &" + tmp_var
	case "text":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := " +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(string)`
		ret = ret + INDENT_1 + inDent + PARAMS_RET + "." + fieldName  + " = &" + tmp_var
	case "blob":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := " +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].([]uint8)`
		ret = ret + INDENT_1 + inDent + fieldName + " = string(" + tmp_var + ")"
		ret = ret + INDENT_1 + inDent + PARAMS_RET + "." + fieldName  + " = &" + fieldName
	case "date": fallthrough
	case "timestamp":
		timeFound = true
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + fieldName + " = " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(time.Time)`
		ret = ret + INDENT_1 + inDent + tmp_var + " := " + fieldName + ".String()"
		//tmp, tmp1 := ProcessTime( timeFound, INDENT_1, timeVar, "", fieldName )
		//ret = ret + tmp
		//ret = ret  + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = &" + tmp1
		ret = ret  + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = &" + tmp_var
	case "timeuuid":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + fieldName + " = " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(gocql.UUID)`
		ret = ret + INDENT_1 + inDent + tmp_var + " := " + fieldName + ".String()"
		ret = ret  + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = &" + tmp_var
	case "uuid":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(gocql.UUID)`
		ret = ret + INDENT_1 + inDent + fieldName + " = " + tmp_var + ".String()"
		ret = ret  + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = &" + fieldName
	case "set": fallthrough
	case "list":
		collectionType := GetFieldName(debug, recursing, fieldDetails.DbFieldCollectionType, false )
		if swagger.IsFieldTypeUDT( parserOutput, collectionType ) {
			arrayType := collectionType
			tmp_var := createTempVar( collectionType )
			returnedVar := PARAMS_RET + "." + fieldName
			ret = INDENT_1 + inDent + tmp_var + ", ok := " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].([]map[string]interface{})`
			ret = ret + INDENT_1 + inDent +  "if ! ok {" + INDENT_1 + INDENT + `log.Fatal("handleReturnedVar() - failed to find entry for ` + strings.ToLower(fieldDetails.DbFieldName ) + `", ok )` + INDENT_1 + "}"
			ret = ret + INDENT_1 + inDent + returnedVar + " = make([]*" + MODELS + collectionType + ", len(" + tmp_var + "))"
			if ! inTable {
				arrayType = GetFieldName(debug, recursing, parserOutput.TypeDetails[typeIndex].TypeName, true) + arrayType
			}
			ret = ret + setUpStructs( debug, recursing, timeFound, INDENT, inTable, returnedVar, tmp_var, fieldDetails.DbFieldCollectionType, parserOutput, timeVar )

		} else {
			tmp_var := createTempVar( fieldName )
			ret = CopyArrayElements( debug, true, inTable, INDENT_1 + inDent, tmp_var, PARAMS_RET + "." + fieldName,  fieldDetails, parserOutput )
		}
		// ret = ret + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = " + fieldName
		// We don't need the above because it won't give us what we need
	case "map" :
		indexCounter++
		iIndex := "i" + strconv.Itoa(indexCounter)
		tmp_var := createTempVar( fieldName )
		mapTypeInGo, isUDT := getMapType(  debug , recursing , inTable , fieldDetails.DbFieldMapType , fieldDetails , parserOutput )
		mapTypeToUse := mapTypeInGo
		tmpMapVarType := mapTypeInGo
		tmpMapVar :=""
		tmp1 := ""
		if isUDT {
			mapTypeToUse = "map[string]interface{}"
			tmpMapVar = createTempVar( fieldName )
			tmpMapVarType = MODELS + mapTypeInGo
			tmp := INDENT_1 + inDent + INDENT + tmpMapVar + " := " + mapTypeInGo + "{}"
			//ts := findTypeDetails( debug, fieldDetails.DbFieldMapType, parserOutput)
			mapFieldType, uDTTypeDetails := manageMap(debug, recursing, inDent + INDENT, inTable, true, tmpMapVar,"v", findTypeDetails( debug, mapTypeInGo, parserOutput ), parserOutput, timeVar )
			if uDTTypeDetails != nil {
				log.Fatal( "Sorry currently unable to handle map types that contain UDTs themselves ")
			}
			tmp1 = tmp + mapFieldType
		}
		ret = INDENT_1 + inDent + tmp_var + ", ok := " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(map[string]` + mapTypeToUse + ")"
		ret = ret + INDENT_1 + inDent +  "if ! ok {" + INDENT_1 + INDENT + `log.Fatal("handleReturnedVar() - failed to find entry for ` + strings.ToLower(fieldDetails.DbFieldName ) + `", ok )` + INDENT_1 + "}"
		ret = ret + INDENT_1 + inDent +  PARAMS_RET + "." + fieldName + " = make(map[string]" + tmpMapVarType + ",len(" + tmp_var + "))"
		//ret = ret + INDENT_1 + inDent + "for " + iIndex +", v := range " + tmp_var + " {" +  INDENT_1 + inDent + INDENT + PARAMS_RET + "." + fieldName  + "[" +  iIndex + "] = v" +  INDENT_1 + inDent + "}" // Modify this part!
		ret = ret + INDENT_1 + inDent + "for " + iIndex +", v := range " + tmp_var + " {"
		if tmp1 != "" { // Processing a UDT
			ret = ret + INDENT_1 + inDent + tmp1;
			typeDetails := findTypeDetails( debug, mapTypeInGo, parserOutput )
			dest := PARAMS_RET + "." + fieldName  + "[" +  iIndex + "]"
			tmpModelType := createTempVar( fieldName )
			ret = ret + INDENT_1 + inDent + INDENT + tmpModelType + " := " + tmpMapVarType + "{}"
			ret = ret + convertToModelType( debug, inDent + INDENT , inTable, tmpMapVar, tmpModelType, typeDetails, parserOutput )
			ret = ret + INDENT_1 + inDent + INDENT + dest + " = " + tmpModelType
			//@TODO
		} else  {
			ret = ret + INDENT_1 + inDent + INDENT + PARAMS_RET + "." + fieldName  + "[" +  iIndex + "] = v"
		}
		ret = ret + INDENT_1 + inDent + "}"

	default:
		if swagger.IsFieldTypeUDT( parserOutput, fieldDetails.DbFieldType ) {
			if inTable {
				// UDTs can only appear as singular fields in table and not in other UDTs. Nope, but too hard to handled :-( !
				tmp_var := createTempVar(fieldDetails.DbFieldName)
				ret_struct := createTempVar(fieldDetails.DbFieldName)
				local_struct := createTempVar(fieldDetails.DbFieldName)
				theType := GetFieldName(debug, recursing, fieldDetails.DbFieldType, false)
				ret = INDENT_1 + inDent + tmp_var + ", ok := " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(map[string]interface{})`
				ret = ret + INDENT_1 + inDent + ret_struct + " := &" + MODELS + theType + "{}"
				ret = ret + INDENT_1 + inDent + "if ! ok {" + INDENT_1 + INDENT + `log.Fatal("handleReturnedVar() - failed to find entry for ` + strings.ToLower(fieldDetails.DbFieldName) + `", ok )` + INDENT_1 + "}"
				tmpStruct,_ := setUpStruct(debug, recursing, timeFound, "", inTable, tmp_var, ret_struct, local_struct, fieldDetails.DbFieldType, parserOutput, timeVar)
				ret = ret + tmpStruct
				ret = ret + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = " + ret_struct
			} else { // UDT in UDT @TODO
				ret = ret + INDENT_1  + inDent + "SJF HERE"
			}
		} else {
			ret = INDENT_1  + inDent + PARAMS_RET + "." + fieldName + " = &" + fieldName
		}


	}
    return ret, timeFound
}



func handleSelectReturn( debug bool, parserOutput parser.ParseOutput, timeVar string ) string {
	ret := ""
	tmp := ""
	timeFound := false
	for i :=0; i < parserOutput.TableDetails.TableFields.FieldIndex; i++ {
		tmp, timeFound = handleReturnedVar( debug, false, timeFound, "", true, 0, parserOutput.TableDetails.TableFields.DbFieldDetails[i], parserOutput, timeVar )
		ret = ret + tmp
	}
	return ret
}


func setupPostParams( debug bool, parserOutput parser.ParseOutput, tableName string ) string {

	ret := INDENT_1
	tmp := buildSelectParams( debug , parserOutput )
	for i, v := range strings.Split(tmp, ", ") {
		ret = ret + processPostField( debug, v, parserOutput, parserOutput.TableDetails.TableFields.DbFieldDetails[i] )
	}

	return ret
}


func createInsert(debug bool, parserOutput parser.ParseOutput, tableName string, cassandraConsistencyRequired string  )string {

	consistency := "gocql.One"
	if cassandraConsistencyRequired != "" {
		consistency = cassandraConsistencyRequired
	}

	ret := INDENT_1 + "m := make(map[string]interface{})"
	ret = ret + INDENT_1 + setupPostParams( debug, parserOutput, swagger.CapitaliseSplitTableName(debug, tableName) )
	ret = ret + INDENT_1 + "if err := " + SESSION_VAR + ".Query(" + "`" + " INSERT INTO " + tableName + "("
	tmp := buildSelectParams( debug , parserOutput )
	tmp1 := ""
	ret = ret + tmp + ") VALUES ("

	s := strings.Split(tmp, ", ")
	for i, v := range s {
		ret = ret + "?"
		tmp1 = tmp1 + "m[" + `"` + v + `"` + "]"
		if i < len(s) -1 {
			tmp1 = tmp1 + ","
			ret = ret + ","
		}
	}
	ret = ret + ")`,"

	ret = ret + tmp1 + ").Consistency(" + consistency + ").Exec(); err != nil {" +  INDENT_1 + INDENT + "return " + tableName + "." + "NewAdd" + swagger.CapitaliseSplitTableName(debug, tableName) + "MethodNotAllowed()" + INDENT_1 + "}"

	return ret
}

//NewAddAccountsMethodNotAllowed()

func handlePost(debug bool, parserOutput parser.ParseOutput, cassandraConsistencyRequired string ) string {

	tableName := strings.ToLower(parserOutput.TableDetails.TableName)

	ret:= createInsert( debug, parserOutput, tableName, cassandraConsistencyRequired  )


	ret = ret + INDENT_1 + "return " + tableName + "." + "NewAdd" + swagger.CapitaliseSplitTableName(debug, tableName) + "Created()" + `
}`

	return ret
}

// Entry point
func CreateCode( debug bool, generateDir string,  goPathForRepo string,  parserOutput parser.ParseOutput, cassandraConsistencyRequired string, endPointNameOveride string, overridePrimaryKeys int, allowFiltering bool, logExtraInfo bool, addPost bool    ) {
	indexCounter = 0
	counter = 0
	output := CreateFile( debug, generateDir, "/data", MAINFILE )
	defer output.Close()


	doNeedTimeImports := WriteHeaderPart( debug, parserOutput, goPathForRepo, endPointNameOveride, false, addPost, output )
	if doNeedTimeImports {
		output.WriteString( PARSETIME )
	}

	if addStruct( debug, addPost, parserOutput, output ) {
		tmp := setUpStuctArrayFromSwaggerParams( debug, parserOutput )
		output.WriteString( "\n    " + tmp + "\n")
	}

	// Write out the static part of the header
	tmpName := GetFieldName(debug, false, parserOutput.TableDetails.TableName, false)
	if endPointNameOveride != "" {
		tmpName = GetFieldName( debug, false, endPointNameOveride, false)
	}
	tmpData := &tableDetails{ generateDir, strings.ToLower(parserOutput.TableSpace), tmpName, tmpName}
	WriteStringToFileWithTemplates(  "\n" + HEADER, "header", output, &tmpData)
	tmpTimeVar := WriteVars( debug, true, parserOutput, goPathForRepo, doNeedTimeImports, addPost, endPointNameOveride, output )
	tmp := createSelectString( debug , parserOutput, tmpName, tmpTimeVar, cassandraConsistencyRequired, overridePrimaryKeys, allowFiltering, logExtraInfo, output )
	output.WriteString( INDENT_1 + "if err := " + SESSION_VAR + ".Query(" + "`" + " SELECT " + tmp )
	tmp = handleSelectReturn( debug, parserOutput, tmpTimeVar )
	tmp = tmp + INDENT_1 + "return operations.NewGet" + tmpName + "OK().WithPayload( " + PAYLOAD + "." + PAYLOAD_STRUCT + ")" + INDENT_1 + "}"
	output.WriteString( tmp )

	if addPost {
		tmpData = &tableDetails{ generateDir, strings.ToLower(parserOutput.TableDetails.TableName), swagger.CapitaliseSplitTableName(debug, parserOutput.TableDetails.TableName), tmpName}
		WriteStringToFileWithTemplates(  "\n" + POST_HEADER, "headerpost", output, &tmpData)
		tmp = handlePost( debug , parserOutput, cassandraConsistencyRequired )
		output.WriteString( tmp )
	}
}


