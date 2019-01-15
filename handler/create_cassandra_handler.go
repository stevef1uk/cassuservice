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
	WriteSwaggerParttoFile(  COMMONIMPORTS + extraImports + IMPORTSEND , "codegen-get", output, &tmpData)


	return doNeedTimeImports

}


// Write out the types for the UDT. Ensure structure type is lowercased to prevent unintentional export of structure beyond package
func addStruct( debug bool, parserOutput parser.ParseOutput, dontUpdate bool, output  *os.File ) {

	for i := 0; i < parserOutput.TypeIndex; i++ {
		v := parserOutput.TypeDetails[i]
		output.WriteString( "\ntype " + CapitaliseSplitFieldName( debug, strings.ToLower(v.TypeName),dontUpdate)  + " struct {")
		for j := 0; j < v.TypeFields.FieldIndex ; j++ {
			revisedFieldName := CapitaliseSplitFieldName(debug, strings.ToLower( v.TypeFields.DbFieldDetails[j].DbFieldName ), dontUpdate )
			revisedType := CapitaliseSplitFieldName( debug, strings.ToLower(v.TypeName),dontUpdate)
			output.WriteString( "\n    " + revisedFieldName + " ")
			output.WriteString( mapCassandraTypeToGoType( debug, revisedFieldName, strings.ToLower(v.TypeFields.DbFieldDetails[j].DbFieldType), revisedType, v.TypeFields.DbFieldDetails[j], parserOutput, false, false, true, true )  + " `" + `cql:"` + strings.ToLower( v.TypeFields.DbFieldDetails[j].DbFieldName ) + `"` +"`")
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
func retArrayTypes(debug bool, field parser.FieldDetails, dontUpdate bool ) string {
	ret := ""
	v :=  field
	if ( v.DbFieldType == "map" ) {
		ret = ret + "payLoad." + CapitaliseSplitFieldName( debug, v.DbFieldName, dontUpdate ) + " = " + v.DbFieldName
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


func writeField( debug bool, parserOutput parser.ParseOutput, field parser.FieldDetails, dontUpdate bool, output  *os.File) {

	fieldName := strings.ToLower(field.DbFieldName)
	if field.DbFieldCollectionType != "" {
		collectionofUDT := swagger.IsFieldTypeUDT(  parserOutput, field.DbFieldCollectionType )
		fieldType :=  mapCassandraTypeToGoType( debug, strings.ToLower(field.DbFieldName), strings.ToLower(field.DbFieldCollectionType), field.DbFieldCollectionType, field, parserOutput, collectionofUDT,  false, false, false)
		if debug {fmt.Println("writeField name =", field.DbFieldName, " fieldType = ", fieldType) }
		if strings.ToLower(fieldType ) == "map" {
			output.WriteString( INDENT_1 + "var " + strings.ToLower( field.DbFieldName  )+ " models." +  CapitaliseSplitFieldName( debug, strings.ToLower(field.DbFieldName),dontUpdate) )
		}

	} else {
		fieldType :=  mapCassandraTypeToGoType( debug, strings.ToLower(field.DbFieldName), strings.ToLower(field.DbFieldType), field.DbFieldCollectionType, field, parserOutput, false,  false, false, false)
		if debug {fmt.Println("writeField name =", field.DbFieldName, " fieldType = ", fieldType) }
		output.WriteString( INDENT_1 + "var " + fieldName + " " + fieldType )
	}
}


// Function that writes out the variable types for the table
func WriteVars(  debug bool, parserOutput parser.ParseOutput, goPathForRepo string, doNeedTimeImports bool, dontUpdate bool, output  *os.File ) {
	// Need to output parameter variable

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
			fieldName := CapitaliseSplitFieldName(debug, strings.ToLower(v.DbFieldName), dontUpdate)
			output.WriteString( INDENT_1 + fieldName + " := &" + strings.ToLower( v.DbFieldType ) + "{}" )

		} else {
			if debug {fmt.Println("WriteVars writing field") }
			writeField( debug, parserOutput, v, dontUpdate, output)
		}
	}

}


// Entry point
func CreateCode( debug bool, generateDir string,  goPathForRepo string,  parserOutput parser.ParseOutput, cassandraConsistencyRequired string, endPointNameOverRide string, overridePrimaryKeys int, allowFiltering bool, dontUpdate bool, logExtraInfo bool   ) {

	output := CreateFile( debug, generateDir, "/data" )
	defer output.Close()


	doNeedTimeImports := WriteHeaderPart( debug, parserOutput, goPathForRepo, endPointNameOverRide, dontUpdate, output )
	addStruct( debug, parserOutput,dontUpdate, output )
	// Write out the static part of the header
	output.WriteString( "\n" + HEADER)
	WriteVars( debug, parserOutput, goPathForRepo, doNeedTimeImports,dontUpdate, output )

}


