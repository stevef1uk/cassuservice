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
	//genString := getServiceName ( parserOutput.TableDetails.TableName, endPointNameOverRide )

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
func addStruct( debug bool, parserOutput parser.ParseOutput, dontUpdate bool, output  *os.File ) {

	for i := 0; i < parserOutput.TypeIndex; i++ {
		v := parserOutput.TypeDetails[i]
		output.WriteString( "\ntype " + GetFieldName( debug, false, v.TypeName,dontUpdate)  + " struct {")
		for j := 0; j < v.TypeFields.FieldIndex ; j++ {
			revisedFieldName := GetFieldName(debug, false, v.TypeFields.DbFieldDetails[j].DbFieldName , dontUpdate )
			revisedType := GetFieldName( debug, false, v.TypeName,dontUpdate)
			tmp := mapFieldTypeToGoCSQLType( debug,  revisedFieldName, true,false, v.TypeFields.DbFieldDetails[j].DbFieldType, revisedType, v.TypeFields.DbFieldDetails[j], parserOutput, dontUpdate  )
			//debug bool, fieldName string, leaveFieldCase bool, inTable bool, fieldType string, typeName string, fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, dontUpdate bool
			output.WriteString( "\n    " + revisedFieldName + " ")
			output.WriteString( tmp   + " `" + `cql:"` + strings.ToLower( v.TypeFields.DbFieldDetails[j].DbFieldName ) + `"` +"`")
		}
		output.WriteString("\n}\n" )
	}

}


func writeField( debug bool, parserOutput parser.ParseOutput, field parser.FieldDetails, dontUpdate bool, output  *os.File) {

	fieldName := GetFieldName( debug, false, field.DbFieldName, dontUpdate)

	if field.DbFieldCollectionType != "" {
		collectionType := GetFieldName(debug, false, field.DbFieldCollectionType, false )
		fieldType :=  mapTableTypeToGoType( debug, fieldName, collectionType, field.DbFieldCollectionType, field, parserOutput, dontUpdate )
		if debug {fmt.Println("writeField name =", field.DbFieldName, " fieldType = ", fieldType) }
		if strings.ToLower(fieldType ) == "map" {
			output.WriteString( INDENT_1 + "var " + fieldName + " " +  fieldName )
		} else {
			if swagger.IsFieldTypeUDT( parserOutput, field.DbFieldCollectionType ) {
				output.WriteString( INDENT_1 + "var " + fieldName + " " + fieldType )
			} else {
				output.WriteString( INDENT_1 + "var " + fieldName + " []" + fieldType )
			}
		}
	} else {
		fieldType :=  mapTableTypeToGoType( debug, strings.ToLower(field.DbFieldName), field.DbFieldType, field.DbFieldCollectionType, field, parserOutput, dontUpdate )
		if debug {fmt.Println("writeField name =", fieldName, " fieldType = ", fieldType) }
		output.WriteString( INDENT_1 + "var " + fieldName + " " + fieldType )
	}
}


// Function that writes out the variable types for the table & returns the temporary variable created if there is a time field
func WriteVars(  debug bool, parserOutput parser.ParseOutput, goPathForRepo string, doNeedTimeImports bool, dontUpdate bool, endPointNameOverRide string, output  *os.File )  string {
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
			fieldName := GetFieldName(debug, false, v.DbFieldName, dontUpdate)
			output.WriteString( INDENT_1 + fieldName + " := &" + strings.ToLower( v.DbFieldType ) + "{}" )

		} else {
			if debug {fmt.Println("WriteVars writing field") }
			writeField( debug, parserOutput, v, dontUpdate, output)
		}
	}

	if doNeedTimeImports {
		tmpTimeVar = createTempVar( TMP_TIME_VAR_PREFIX )
		output.WriteString( INDENT_1 + tmpTimeVar + " := strfmt.NewDateTime().String()" )
	}
	output.WriteString( "\n" + INDENT_1 + SELECT_OUTPUT + " := map[string]interface{}{}\n")


	return tmpTimeVar
}


func buildSelectParams ( debug bool, parserOutput parser.ParseOutput,  dontUpdate bool )  string {
	ret := ""
	for i :=0; i < parserOutput.TableDetails.TableFields.FieldIndex; i++ {
		if i > 0 {
			ret = ret + ", "
		}
		ret = ret + strings.ToLower( parserOutput.TableDetails.TableFields.DbFieldDetails[i].DbFieldName )
	}
	return ret
}


func createSelectString( debug bool, parserOutput parser.ParseOutput, timeVar string, cassandraConsistencyRequired string,  overridePrimaryKeys int, allowFiltering bool, dontUpdate bool, logExtraInfo bool, output  *os.File )  string {

	ret := buildSelectParams( debug, parserOutput, dontUpdate )
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

		fieldName := GetFieldName(debug, false, v.DbFieldName, dontUpdate)
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
		consistency = strings.ToLower(cassandraConsistencyRequired)
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
func handleStructVarConversion(  debug bool, recursing bool, indexCounter int, timeFound bool, inDent string, theStructVar string, destVar string,  theType string,  fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, timeVar string, dontUpdate bool ) string {

	ret := ""

	fieldName := GetFieldName( debug, false, fieldDetails.DbFieldName, false)
	sourceVar := theStructVar + "." +  fieldName
	switch ( strings.ToLower( fieldDetails.DbFieldType ) ) {
	case "int":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := int64(" + sourceVar + ")"
		ret = ret + INDENT_1 + inDent + destVar + " = &" + tmp_var
	case "float":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + inDent + tmp_var + " := float64(" + sourceVar + ")"
		ret = ret + INDENT_1 + inDent + destVar + " = &" + tmp_var
	case "date": fallthrough
	case "timestamp": fallthrough
	case "timeuuid":
		tmp, tmp1 := ProcessTime( timeFound, INDENT_1 + inDent, timeVar, sourceVar )
		ret = ret + tmp
		ret = ret  + INDENT_1 + inDent + destVar + " = &" + tmp1
	case "set": fallthrough
	case "list":
		//tmp_var := createTempVar( fieldName )
		collectionType := GetFieldName(debug, recursing, fieldDetails.DbFieldCollectionType, dontUpdate )
		if swagger.IsFieldTypeUDT( parserOutput, collectionType ) {
			ret = ret +  setUpStruct( debug ,true ,timeFound , inDent, destVar, sourceVar,  collectionType, parserOutput, timeVar, dontUpdate )
		} else {
			ret = "OOPS"
			//ret = CopyArrayElements( debug, false, INDENT_1 + INDENT, tmp_var, tmp_var,  fieldDetails, parserOutput, dontUpdate  )
		}
	default:
		ret = INDENT_1 + inDent + destVar + " = &" + sourceVar

	}
	return ret
}


var indexCounter int = 0
func setUpStruct ( debug bool, recursing bool,  timeFound bool, inDent string, destField string, theVar string, theType string,  parserOutput parser.ParseOutput, timeVar string, dontUpdate bool  ) string {

	indexCounter++
	typeStruct := findTypeDetails( debug, theType, parserOutput )
	structName := GetFieldName(  debug, recursing, theType, false )
	tmpStruct := createTempVar( structName )
	iIndex := "i" + strconv.Itoa(indexCounter)

	ret := INDENT_1 + inDent + "for " + iIndex +", _ := range " + theVar + " {" + INDENT_1 + inDent + INDENT + tmpStruct + " := &" + structName + "{"
	for i := 0; i < typeStruct.TypeFields.FieldIndex; i++ {
		fieldName := strings.ToLower(GetFieldName( debug, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldName, false))
		fieldType := mapFieldTypeToGoCSQLType( debug, fieldName, recursing, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldType, structName, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, dontUpdate  )
		if recursing {
			inDent = inDent + INDENT
		}
		ret = ret + INDENT_1 + inDent + INDENT2  + "v[" + fieldName + "].(" + fieldType + "),"
	}
	ret = ret + INDENT_1 + inDent + INDENT + "}"

	// Now process each variable in order to set-up the Payload structure
	for i := 0; i < typeStruct.TypeFields.FieldIndex; i++ {
		fieldName := GetFieldName( debug, false, typeStruct.TypeFields.DbFieldDetails[i].DbFieldName, false)
		tmp := handleStructVarConversion( debug, recursing, indexCounter, timeFound, inDent, tmpStruct, destField + "[" + iIndex + "]." + fieldName, typeStruct.TypeName, typeStruct.TypeFields.DbFieldDetails[i], parserOutput, timeVar, dontUpdate )
		ret = ret + inDent + INDENT2 + tmp
	}

	ret = ret + INDENT_1 + inDent + "}"

	return ret
}


func handleReturnedVar( debug bool, recursing bool, timeFound bool, inTable bool, typeIndex int , fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, timeVar string, dontUpdate bool ) (string, bool) {
    ret := ""
	fieldName := GetFieldName( debug, false, fieldDetails.DbFieldName, false)
	switch ( strings.ToLower( fieldDetails.DbFieldType ) ) {
	case "int":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + tmp_var + " := " +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(int)`
		ret = ret + INDENT_1 + fieldName + " = int64(" + tmp_var + ")"
		ret = ret + INDENT_1 + PARAMS_RET + "." + fieldName + " = &" + fieldName
	case "float":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + tmp_var + " := " +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(float32)`
		ret = ret + INDENT_1 + fieldName + " = float64(" + tmp_var + ")"
		ret = ret + INDENT_1 + PARAMS_RET + "." + fieldName + " = &" + fieldName
	case "date": fallthrough
	case "timestamp": fallthrough
	case "timeuuid":
		timeFound = true
		ret = ret + INDENT_1 + fieldName + " = " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(time.Time)`
		tmp, tmp1 := ProcessTime( timeFound, INDENT_1, timeVar, fieldName )
		ret = ret + tmp
		ret = ret  + INDENT_1 + PARAMS_RET + "." + fieldName + " = &" + tmp1
	case "nottimestamp":
		// @TODO need to check return type!
		ret = ret + INDENT_1 + fieldName + " = " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(string)`
		ret = ret + INDENT_1 + PARAMS_RET + "." + fieldName + " = &" + fieldName
	case "set": fallthrough
	case "list":
		collectionType := GetFieldName(debug, recursing, fieldDetails.DbFieldCollectionType, dontUpdate )
		if swagger.IsFieldTypeUDT( parserOutput, collectionType ) {
			arrayType := collectionType
			tmp_var := createTempVar( collectionType )
			returnedVar := PARAMS_RET + "." + fieldName
			ret = INDENT_1 + tmp_var + ", ok := " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].([]map[string]interface{})`
			ret = ret + INDENT_1 +  "if ! ok {" + INDENT_1 + INDENT + `log.Fatal("handleReturnedVar() - failed to find entry for ` + strings.ToLower(fieldDetails.DbFieldName ) + `", ok )` + INDENT_1 + "}"
			ret = ret + INDENT_1 + returnedVar + " = make([]*" + MODELS + collectionType + ", len(" + tmp_var + ")"
			if ! inTable {
				arrayType = GetFieldName(debug, recursing, parserOutput.TypeDetails[typeIndex].TypeName, dontUpdate) + arrayType
			}
			ret = ret + setUpStruct( debug, recursing, timeFound, INDENT, returnedVar, tmp_var, fieldDetails.DbFieldCollectionType, parserOutput, timeVar, dontUpdate )

		} else {
			ret = CopyArrayElements( debug, inTable, INDENT_1, fieldName, PARAMS_RET + "." + fieldName,  fieldDetails, parserOutput, dontUpdate  )
		}
	default:
		ret = INDENT_1  + PARAMS_RET + "." + fieldName + " = &" + fieldName

	}
    return ret, timeFound
}



func handleSelectReturn( debug bool, parserOutput parser.ParseOutput, timeVar string, dontUpdate bool ) string {
	ret := ""
	tmp := ""
	timeFound := false
	for i :=0; i < parserOutput.TableDetails.TableFields.FieldIndex; i++ {
		tmp, timeFound = handleReturnedVar( debug, false, timeFound, true, 0, parserOutput.TableDetails.TableFields.DbFieldDetails[i], parserOutput, timeVar, dontUpdate )
		ret = ret + tmp
	}
	return ret
}


// Entry point
func CreateCode( debug bool, generateDir string,  goPathForRepo string,  parserOutput parser.ParseOutput, cassandraConsistencyRequired string, endPointNameOverRide string, overridePrimaryKeys int, allowFiltering bool, dontUpdate bool, logExtraInfo bool   ) {

	output := CreateFile( debug, generateDir, "/data" )
	defer output.Close()


	doNeedTimeImports := WriteHeaderPart( debug, parserOutput, goPathForRepo, endPointNameOverRide, dontUpdate, output )
	addStruct( debug, parserOutput,dontUpdate, output )
	// Write out the static part of the header
	tmpName := GetFieldName(debug, false, parserOutput.TableDetails.TableName, false)
	if endPointNameOverRide != "" {
		tmpName = GetFieldName( debug, false, endPointNameOverRide, false)
	}
	tmpData := &tableDetails{ generateDir, strings.ToLower(parserOutput.TableSpace), strings.ToLower(parserOutput.TableDetails.TableName), tmpName}
	WriteStringToFileWithTemplates(  "\n" + HEADER, "header", output, &tmpData)
	//output.WriteString( "\n" + HEADER)
	tmpTimeVar := WriteVars( debug, parserOutput, goPathForRepo, doNeedTimeImports,dontUpdate, endPointNameOverRide, output )
	tmp := createSelectString( debug , parserOutput, tmpTimeVar, cassandraConsistencyRequired, overridePrimaryKeys, allowFiltering, dontUpdate, logExtraInfo, output )
	output.WriteString( INDENT_1 + "if err := " + SESSION_VAR + ".Query(" + "`" + " SELECT " + tmp )
	tmp = handleSelectReturn( debug, parserOutput, tmpTimeVar, dontUpdate )
	tmp = tmp + INDENT_1 + "return operations.NewGet" + tmpName + "OK().WithPayload( " + PAYLOAD + "." + PAYLOAD_STRUCT + ")" + INDENT_1 + "}"
	output.WriteString( tmp )
 

}


