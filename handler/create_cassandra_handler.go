package handler

import (
	"fmt"
	"github.com/stevef1uk/cassuservice/swagger"
	"strconv"
	//"fmt"
	"github.com/stevef1uk/cassuservice/parser"
	"os"
	"strings"

	//"strings"
)



func WriteHeaderPart( debug bool, parserOutput parser.ParseOutput, generateDir string, endPointNameOverRide string, dontUpdate bool, output  *os.File ) bool  {
	doNeedTimeImports := doINeedTime(parserOutput )
	needDecimalImports := doINeedDecimal(parserOutput )

	extraImports := ""
	if doNeedTimeImports {
		extraImports = extraImports + IMPORTSTIMESTAMP
	}
	if needDecimalImports {
		extraImports = extraImports + IMPORTDEC
	}

	tmpData := &tableDetails{ generateDir, "", "",""}
	WriteStringToFileWithTemplates(  COMMONIMPORTS + extraImports + IMPORTSEND , "codegen-get", output, &tmpData)

	return doNeedTimeImports
}


// Write out the types for the UDT. Ensure structure type is lowercased to prevent unintentional export of structure beyond package
func addStruct( debug bool, parserOutput parser.ParseOutput, output  *os.File ) {

	for i := 0; i < parserOutput.TypeIndex; i++ {
		v := parserOutput.TypeDetails[i]
		output.WriteString( "\ntype " + GetFieldName( debug, false, v.TypeName, true)  + " struct {")
		for j := 0; j < v.TypeFields.FieldIndex ; j++ {
			revisedFieldName := GetFieldName(debug, false, v.TypeFields.DbFieldDetails[j].OrigFieldName , false )
			revisedType := GetFieldName( debug, false, v.TypeName,true)
			tmp := mapFieldTypeToGoCSQLType( debug,  revisedFieldName, true,false, v.TypeFields.DbFieldDetails[j].DbFieldType, revisedType, v.TypeFields.DbFieldDetails[j], parserOutput, true  )
			//debug bool, fieldName string, leaveFieldCase bool, inTable bool, fieldType string, typeName string, fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, dontUpdate bool
			output.WriteString( "\n    " + revisedFieldName + " ")
			output.WriteString( tmp   + " `" + `cql:"` + strings.ToLower( v.TypeFields.DbFieldDetails[j].DbFieldName ) + `"` +"`")
		}
		output.WriteString("\n}\n" )
	}

}


func writeField( debug bool, parserOutput parser.ParseOutput, field parser.FieldDetails, output  *os.File) {

	fieldName := GetFieldName( debug, false, field.OrigFieldName, false)

	if field.DbFieldCollectionType != "" {
		collectionType := GetFieldName(debug, false, field.DbFieldCollectionType, true )
		fieldType :=  mapTableTypeToGoType( debug, fieldName, collectionType, field.DbFieldCollectionType, field, parserOutput, true )
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
		fieldType :=  mapTableTypeToGoType( debug, strings.ToLower(field.DbFieldName), field.DbFieldType, field.DbFieldCollectionType, field, parserOutput, true )
		if debug {fmt.Println("writeField name =", fieldName, " fieldType = ", fieldType) }
		output.WriteString( INDENT_1 + "var " + fieldName + " " + fieldType )
	}
	// Not all vars used so handle any compilation errors
	output.WriteString( INDENT_1 + "_ = " + fieldName  )
}


// Function that writes out the variable types for the table & returns the temporary variable created if there is a time field
func WriteVars(  debug bool, parserOutput parser.ParseOutput, goPathForRepo string, doNeedTimeImports bool, endPointNameOverRide string, output  *os.File )  string {
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
			writeField( debug, parserOutput, v, output)
		}
	}

	if doNeedTimeImports {
		tmpTimeVar = createTempVar( TMP_TIME_VAR_PREFIX )
		//output.WriteString( INDENT_1 + tmpTimeVar + " := strfmt.NewDateTime().String()" )
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


func createSelectString( debug bool, parserOutput parser.ParseOutput, timeVar string, cassandraConsistencyRequired string,  overridePrimaryKeys int, allowFiltering bool, logExtraInfo bool, output  *os.File )  string {

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

	ret = ret + " FROM " + strings.ToLower( parserOutput.TableDetails.TableName) + whereClause +  "`," + varsClause + ")"
	ret = ret + ".Consistency(" + consistency + ").MapScan(codeGenRawTableResult); err != nil {"
	if logExtraInfo {
		ret = ret + INDENT_1 + `  log.Println("No data? ", err)`
	}
	tableName := GetFieldName(debug, false, parserOutput.TableDetails.TableName, false)
	ret = ret + INDENT_1 + "  return " + OPERATIONS + "NewGet" + tableName + "BadRequest()" +  INDENT_1 + "}"
	ret = ret + INDENT_1 + PAYLOAD + " := " + OPERATIONS + "NewGet" +  tableName + "OK()"
	ret = ret + INDENT_1 + PAYLOAD + "." + PAYLOAD_STRUCT + " = make([]*" + MODELS + "Get" + tableName + "OKBodyItems,1)"
	ret = ret + INDENT_1 + PAYLOAD + "." + PAYLOAD_STRUCT + "[0] = new(" + MODELS + "Get" + tableName + "OKBodyItems)"
	ret = ret + INDENT_1 + PARAMS_RET + " := " + PAYLOAD + "." + PAYLOAD_STRUCT + "[0]"

	return ret
}

// Function called to process a local UDT structure and copy into the go-swagger model's structure type for the UDT
func handleStructVarConversion(  debug bool, recursing bool, indexCounter int, timeFound bool, inDent string, theStructVar string, destVar string,  theType * parser.TypeDetails,  fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, timeVar string ) string {

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
			ret = CopyArrayElements( debug, false, INDENT_1 + INDENT, theStructVar + "." + fieldName, destVar,  fieldDetails, parserOutput )
		}
	default:
		ret = INDENT_1 + inDent + destVar + " = " + sourceVar

	}
	return ret
}


var indexCounter int = 0

// Handle the case of a single UDT field, which can only occur in a Table definition right now
func setUpStruct ( debug bool, recursing bool,  timeFound bool, inDent string, inTable bool, raw_data string, destField string, theVar string, theType string,  parserOutput parser.ParseOutput, timeVar string ) string {

	extraVars := ""
	newStr := ""
	ret := ""
	//iIndex := "i" + strconv.Itoa(indexCounter)
	//vIndex := "v" + strconv.Itoa(indexCounter)
	typeStruct := findTypeDetails( debug, theType, parserOutput )
	var structAssignment []bool = make( []bool, typeStruct.TypeFields.FieldIndex)
	structName := GetFieldName(  debug, recursing, theType, false )
	tmpStruct := createTempVar( structName )

	if recursing {
		inDent = inDent + INDENT2
	}
	/*space := ""
	if recursing {
		space = inDent
	}*/

	loopAssignment := INDENT_1 + inDent  + tmpStruct + " := &" + structName + "{"
	for i := 0; i < typeStruct.TypeFields.FieldIndex; i++ {
		fieldName := strings.ToLower( typeStruct.TypeFields.DbFieldDetails[i].DbFieldName )
		fieldType := mapFieldTypeToGoCSQLType( debug, fieldName, true, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldType, structName, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, true  )
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
				extraVars = extraVars +  INDENT_1 + inDent + tmpVar + " := " + raw_data +  `["` + strings.ToLower(fieldName ) + `"].(map[string]string)`
			} else {
				// Handle lists & sets here!
				ret = ret + INDENT_1 + inDent + INDENT2  + tmpVar1 + ","
				extraVars = extraVars +  INDENT_1 + inDent + tmpVar + ":= " + raw_data + `["` + strings.ToLower(fieldName ) + `"].([]map[string]interface{})`
				tmpType := mapFieldTypeToGoCSQLType( debug, fieldName, true, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType, structName, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, true  )
				extraVars = extraVars +  INDENT_1 + inDent + tmpVar1 + ":= make(" + tmpType + ", len(" + tmpVar + ") )"
				extraVars = extraVars + INDENT_1 + inDent + setUpStructs( debug,  true,  timeFound, inDent, false, tmpVar1,  tmpVar, strings.ToLower(typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType),  parserOutput, timeVar )
				structAssignment[i] = true
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
		//typeName := GetFieldName(debug, false, typeStruct.TypeName, false)
		tmpDest := destField + "." + fieldName
		//addMake := INDENT_1 + inDent + destField +  = &" + MODELS + typeName + "{}"
		if structAssignment[i] {
			if typeStruct.TypeFields.DbFieldDetails[i].DbFieldMapType == "" && ! inTable {
				tmpDest = destField
				tmp = INDENT_1 + inDent + tmpDest + " = " + tmpStruct
				goto here
			}
		}
		tmp = handleStructVarConversion(debug, recursing, indexCounter, timeFound, inDent, tmpStruct, tmpDest, typeStruct, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, timeVar )
	here:
		if i == 0 {
			//ret = ret + inDent + INDENT2 + addMake + tmp
			ret = ret + inDent + INDENT2  + tmp
		} else {
			ret = ret + inDent + INDENT2  + tmp
		}
	}

	//ret = ret + INDENT_1 + inDent + "}"

	return ret
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
	loopAssignment := INDENT_1 + inDent  + tmpStruct + " := &" + structName + "{"
	for i := 0; i < typeStruct.TypeFields.FieldIndex; i++ {
		fieldName := strings.ToLower( typeStruct.TypeFields.DbFieldDetails[i].DbFieldName )
		fieldType := mapFieldTypeToGoCSQLType( debug, fieldName, true, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldType, structName, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, true  )
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
				extraVars = extraVars +  INDENT_1 + inDent + tmpVar + " := " +vIndex +  `["` + strings.ToLower(fieldName ) + `"].(map[string]string)`
			} else {
				// Handle lists & sets here!
				ret = ret + INDENT_1 + inDent + INDENT2  + tmpVar1 + ","
				extraVars = extraVars +  INDENT_1 + inDent + tmpVar + ":= " + vIndex + `["` + strings.ToLower(fieldName ) + `"].([]map[string]interface{})`
				tmpType := mapFieldTypeToGoCSQLType( debug, fieldName, true, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType, structName, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, true  )
				extraVars = extraVars +  INDENT_1 + inDent + tmpVar1 + ":= make(" + tmpType + ", len(" + tmpVar + ") )"
				extraVars = extraVars + INDENT_1 + inDent + setUpStructs( debug,  true,  timeFound, inDent, false, tmpVar1,  tmpVar, strings.ToLower(typeStruct.TypeFields.DbFieldDetails[i].DbFieldCollectionType),  parserOutput, timeVar )
				structAssignment[i] = true
			}
		} else {
			ret = ret + INDENT_1 + inDent + INDENT2  + vIndex + `["` + strings.ToLower(fieldName ) + `"].(` + fieldType + "),"
		}
	}
	ret = ret + INDENT_1 + inDent + INDENT + "}"
	ret = loopStart + INDENT_1 + extraVars + loopAssignment + INDENT_1 + ret + INDENT_1 + newStr

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
		tmp = handleStructVarConversion(debug, recursing, indexCounter, timeFound, inDent, tmpStruct, tmpDest, typeStruct, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, timeVar )
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
			ret = CopyArrayElements( debug, inTable, INDENT_1 + inDent, tmp_var, PARAMS_RET + "." + fieldName,  fieldDetails, parserOutput )
		}
		// ret = ret + INDENT_1 + inDent + PARAMS_RET + "." + fieldName + " = " + fieldName
		// We don't need the above because it won't give us what we need
	case "map" :
		indexCounter++
		iIndex := "i" + strconv.Itoa(indexCounter)
		tmp_var := createTempVar( fieldName )
		mapTypeInGo := "string" // This will always be the case as the swagger generated for maps is always []map[string]string
		ret = INDENT_1 + inDent + tmp_var + ", ok := " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(map[string]` + mapTypeInGo + ")"
		ret = ret + INDENT_1 + inDent +  "if ! ok {" + INDENT_1 + INDENT + `log.Fatal("handleReturnedVar() - failed to find entry for ` + strings.ToLower(fieldDetails.DbFieldName ) + `", ok )` + INDENT_1 + "}"
		ret = ret + INDENT_1 + inDent +  PARAMS_RET + "." + fieldName + " = make(map[string]string,len(" + tmp_var + "))"
		ret = ret + INDENT_1 + inDent + "for " + iIndex +", v := range " + tmp_var + " {" + INDENT_1 + inDent + INDENT + PARAMS_RET + "." + fieldName  + "[" +  iIndex + "] = v"  +  INDENT_1 + inDent + "}"

	default:
		if swagger.IsFieldTypeUDT( parserOutput, fieldDetails.DbFieldType ) {
			// UDTs can only appear as singular fields in table and not in other UDTs
			tmp_var := createTempVar( fieldDetails.DbFieldName )
			ret_struct := createTempVar( fieldDetails.DbFieldName )
			local_struct := createTempVar( fieldDetails.DbFieldName )
			theType := GetFieldName(debug, recursing, fieldDetails.DbFieldType, false )
			ret = INDENT_1 + inDent + tmp_var + ", ok := " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(map[string]interface{})`
			ret = ret + INDENT_1 + inDent + ret_struct + " := &" + MODELS + theType + "{}"
			ret = ret + INDENT_1 + inDent +  "if ! ok {" + INDENT_1 + INDENT + `log.Fatal("handleReturnedVar() - failed to find entry for ` + strings.ToLower(fieldDetails.DbFieldName ) + `", ok )` + INDENT_1 + "}"
			ret = ret + setUpStruct( debug, recursing, timeFound, "", inTable,  tmp_var, ret_struct , local_struct, fieldDetails.DbFieldType, parserOutput, timeVar )
			ret = ret + INDENT_1 + inDent +  PARAMS_RET + "." + fieldName + " = " +  ret_struct
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


// Entry point
func CreateCode( debug bool, generateDir string,  goPathForRepo string,  parserOutput parser.ParseOutput, cassandraConsistencyRequired string, endPointNameOverRide string, overridePrimaryKeys int, allowFiltering bool, logExtraInfo bool   ) {
	indexCounter = 0
	counter = 0
	output := CreateFile( debug, generateDir, "/data", MAINFILE )
	defer output.Close()


	doNeedTimeImports := WriteHeaderPart( debug, parserOutput, goPathForRepo, endPointNameOverRide, false, output )
	addStruct( debug, parserOutput, output )
	// Write out the static part of the header
	tmpName := GetFieldName(debug, false, parserOutput.TableDetails.TableName, false)
	if endPointNameOverRide != "" {
		tmpName = GetFieldName( debug, false, endPointNameOverRide, false)
	}
	tmpData := &tableDetails{ generateDir, strings.ToLower(parserOutput.TableSpace), strings.ToLower(parserOutput.TableDetails.TableName), tmpName}
	WriteStringToFileWithTemplates(  "\n" + HEADER, "header", output, &tmpData)
	//output.WriteString( "\n" + HEADER)
	tmpTimeVar := WriteVars( debug, parserOutput, goPathForRepo, doNeedTimeImports, endPointNameOverRide, output )
	tmp := createSelectString( debug , parserOutput, tmpTimeVar, cassandraConsistencyRequired, overridePrimaryKeys, allowFiltering, logExtraInfo, output )
	output.WriteString( INDENT_1 + "if err := " + SESSION_VAR + ".Query(" + "`" + " SELECT " + tmp )
	tmp = handleSelectReturn( debug, parserOutput, tmpTimeVar )
	tmp = tmp + INDENT_1 + "return operations.NewGet" + tmpName + "OK().WithPayload( " + PAYLOAD + "." + PAYLOAD_STRUCT + ")" + INDENT_1 + "}"
	output.WriteString( tmp )
 

}


