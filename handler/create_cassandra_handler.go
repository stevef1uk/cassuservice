package handler

import (
	//"fmt"
	"github.com/stevef1uk/cassuservice/parser"
	"os"
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


// Entry point
func CreateCode( debug bool, generateDir string,  goPathForRepo string,  parseOutput parser.ParseOutput, cassandraConsistencyRequired string, endPointNameOverRide string, overridePrimaryKeys int, allowFiltering bool, dontUpdate bool, logExtraInfo bool   ) {

	output := CreateFile( debug, generateDir, "/data" )
	defer output.Close()


	WriteHeaderPart( debug, parseOutput, goPathForRepo, endPointNameOverRide, dontUpdate, output )

}


