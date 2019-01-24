package handler

import (
	"fmt"
	"github.com/stevef1uk/cassuservice/swagger"
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
			tmp := basicMapCassandraTypeToGoType( debug, true, false, revisedFieldName, v.TypeFields.DbFieldDetails[j].DbFieldType, revisedType, v.TypeFields.DbFieldDetails[j], parserOutput, dontUpdate  )
			output.WriteString( "\n    " + revisedFieldName + " ")
			output.WriteString( tmp   + " `" + `cql:"` + strings.ToLower( v.TypeFields.DbFieldDetails[j].DbFieldName ) + `"` +"`")
		}
		output.WriteString("\n}\n" )
	}

}

/*
func setUpArrayTypes(  debug bool, field parser.FieldDetails,  dontUpdate bool ) string {
	ret := ""
	tmpVar := createTempVar( field.DbFieldName)

	if  swagger.IsFieldTypeATime( field.DbFieldType ) {
		ret = ret + `
        ` + tmpVar + " = strfmt.NewDateTime().String()" + `
        ` + "_ = " + tmpVar + `
		` + strings.ToLower(field.DbFieldName) + " = " + RAWRESULT + `["` + strings.ToLower(field.DbFieldName) + `"].([]` +
			mapCassandraTypeToGoType( debug, field,false,   false, false ) +  `)
		` + "payLoad." + Capitiseid( debug, field.DbFieldName, dontUpdate) + " = make([] string, len( payLoad." +strings.ToLower(field.DbFieldName) + ") )" + `
		for i := 0; i < len(` + strings.ToLower(field.DbFieldName) + `); i++ {
			payLoad.` + Capitiseid( debug, field.DbFieldName, dontUpdate) + "[i] = " + strings.ToLower(field.DbFieldName) + "[i].String()" + `
		}`
	} else {
		if ( strings.ToLower(field.DbFieldType) == "decimal") {
			ret = ret + `
    payLoad.` + Capitiseid(debug, field.DbFieldName, dontUpdate) + " = make([]float64, len(" + strings.ToLower(field.DbFieldName) + ") )" + `
    for i := 0; i < len( payLoad.` + strings.ToLower(field.DbFieldName) + `); i++ {
        ` + tmpVar + ", err := strconv.ParseFloat( " + strings.ToLower(field.DbFieldName) + "[i].String(), 64 )" + `
        if ( err != nil ) {
            log.Println("error parsing decimal value for field %s\n",` + field.DbFieldName + `)
        }
` + `
        payLoad.` + Capitiseid( debug, field.DbFieldName, dontUpdate) + "[i] = " + tmpVar + `
    }`
		} else {
			ret = ret + `
		` + strings.ToLower(field.DbFieldName) + " = " + RAWRESULT + `["` + strings.ToLower(field.DbFieldName) + `"].([]` + mapCassandraTypeToGoType( debug, field,false,   false, false ) + `)`
			ret = ret + `
		` + "payLoad." + Capitiseid(debug, field.DbFieldName, dontUpdate) + " = make([]" + mapCassandraTypeToGoType( true, field,false,   false, false ) + ",len(" + strings.ToLower(field.DbFieldName) + ") )" + `
		for i := 0; i < len( payload.` + strings.ToLower(field.DbFieldName) + `); i++ {
			payLoad.` + Capitiseid(debug,field.DbFieldName, dontUpdate) + "[i] = " + mapCassandraTypeToGoType( true, field,false,   false, false ) + "(" + strings.ToLower(field.DbFieldName) + "[i])" + `
		}`
		}
	}
	if debug { fmt.Printf("setUpArrayTypes returning %s\n", ret ) }
	return ret
}
*/


// Setup array types from gocql select result
/*
func retArrayTypes(debug bool, field parser.FieldDetails, dontUpdate bool ) string {
	ret := ""
	v :=  field
	if ( v.DbFieldType == "map" ) {
		ret = ret + "payLoad." + GetFieldName( debug, false, v.DbFieldName, dontUpdate ) + " = " + v.DbFieldName
	} else {
		switchValue := strings.ToLower( v.DbFieldCollectionType )
		switch switchValue {
		case "float", "int", "varint", "boolean", "uuid", "bigint", "counter", "decimal", "double", "text", "varchar", "ascii", "blob", "inet", swagger.DATE, swagger.TIMESTAMP, swagger.TIMEUUID :
			//ret = setUpArrayTypes(  debug, v,  dontUpdate  )

		default:
			ret = ret + "payLoad." + CapitaliseSplitFieldName( debug, v.DbFieldName, dontUpdate ) + " = " + v.DbFieldName
		}
	}

	return ret
}
*/

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

	/*
	tableName := parserOutput.TableDetails.TableName

	if endPointNameOverRide != "" {
		tableName = endPointNameOverRide
	}
	tableName = GetFieldName( debug, false, tableName, false)
	//output.WriteString( INDENT_1 + "ret := " + OPERATIONS + "NewGet" + tableName + "OK()\n" )
	//output.WriteString( INDENT_1 + "ret." + PAYLOAD_STRUCT + "= " + OPERATIONS + "NewGet" + tableName + "OK()\n" )
	*/
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
	ret = ret + INDENT_1 + PAYLOAD + "." + PAYLOAD_STRUCT + " = make([]*models.Get" + tableName + "OKBodyItems,1)"
	ret = ret + INDENT_1 + PARAMS_RET + " := " + PAYLOAD + "." + PAYLOAD_STRUCT + "[0]"

	return ret
}




func handleReturnedVar( debug bool, timeFound bool, inTable bool, typeIndex int , fieldDetails parser.FieldDetails, parserOutput parser.ParseOutput, timeVar string, dontUpdate bool ) (string, bool) {
    ret := ""
	fieldName := GetFieldName( debug, false, fieldDetails.DbFieldName, false)
	switch ( strings.ToLower( fieldDetails.DbFieldType ) ) {
	case "int":
		tmp_var := createTempVar( fieldName )
		ret = ret + INDENT_1 + tmp_var + " := " +  SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + `"].(int)`
		ret = ret + INDENT_1 + fieldName + " = int64(" + tmp_var + ")"
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
		collectionType := GetFieldName(debug, false, fieldDetails.DbFieldCollectionType, dontUpdate )
		if swagger.IsFieldTypeUDT( parserOutput, collectionType ) {
			arrayType := collectionType
			tmp_var := createTempVar( collectionType )
			ret = INDENT_1 + tmp_var + ", ok := " + SELECT_OUTPUT + `["` + strings.ToLower(fieldDetails.DbFieldName) + ` "].([]map[string]interface{})`
			ret = ret + INDENT_1 +  "if != ok {" + INDENT_1 + INDENT + `log.Fatal("handleReturnedVar() - failed to find entry for ` + strings.ToLower(fieldDetails.DbFieldName ) + `", ok )` + INDENT_1 + "}"
			if ! inTable {
				arrayType = GetFieldName(debug, false, parserOutput.TypeDetails[typeIndex].TypeName, dontUpdate) + arrayType
			}
			ret = ret + setUpStruct( debug, tmp_var, fieldDetails.DbFieldCollectionType, parserOutput)

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
		tmp, timeFound = handleReturnedVar( debug, timeFound, true, 0, parserOutput.TableDetails.TableFields.DbFieldDetails[i], parserOutput, timeVar, dontUpdate )
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
	output.WriteString( tmp )
 

}


