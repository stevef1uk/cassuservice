package handler

import (
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

	// Write out the imports required
	//output.WriteString( COMMONIMPORTS + extraImports + IMPORTSEND  )

	return doNeedTimeImports

}


// Write out the types for the UDT
func addStruct( debug bool, parserOutput parser.ParseOutput, dontUpdate bool, output  *os.File ) {

	for i := 0; i < parserOutput.TypeIndex; i++ {
		v := parserOutput.TypeDetails[i]
		output.WriteString( "\ntype " + v.TypeName + " struct {")
		for j := 0; j < v.TypeFields.FieldIndex ; j++ {
			revisedFieldName := CapitaliseSplitFieldName(debug, strings.ToLower( v.TypeFields.DbFieldDetails[j].DbFieldName ), dontUpdate )
			output.WriteString( "\n    " + revisedFieldName + " ")
			output.WriteString( mapCassandraTypeToGoType( debug, v.TypeFields.DbFieldDetails[j], false, true, true )  + " `" + `cql:"` + revisedFieldName + `"` +"`")
		}
		output.WriteString("\n}\n" )
	}

}




// Entry point
func CreateCode( debug bool, generateDir string,  goPathForRepo string,  parserOutput parser.ParseOutput, cassandraConsistencyRequired string, endPointNameOverRide string, overridePrimaryKeys int, allowFiltering bool, dontUpdate bool, logExtraInfo bool   ) {

	output := CreateFile( debug, generateDir, "/data" )
	defer output.Close()


	WriteHeaderPart( debug, parserOutput, goPathForRepo, endPointNameOverRide, dontUpdate, output )
	addStruct( debug, parserOutput,dontUpdate, output )

}


