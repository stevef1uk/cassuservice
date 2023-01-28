package main

import (
	"flag"
	"github.com/stevef1uk/cassuservice/handler"
	"github.com/stevef1uk/cassuservice/parser"
	"github.com/stevef1uk/cassuservice/swagger"
	"io/ioutil"
	"os"
	"os/exec"
)
import "fmt"

const (
	SWAGGER_FILE = "swagger.txt"
)


func main() {

	allowFilteringPtr := flag.Bool("allowFiltering", false, "Set flag true to add Allow Filtering on Select queries")
	consistencyPtr := flag.String( "consistency", "gocql.One", "Set required Cassandra Read Consistency level, default = gocql.One")
	debugPtr := flag.Bool("debug", false, "set -debug=true to debug code")
	endPointPtr := flag.String( "endPoint", "", "Set to override the endpoint for uService, which will be by default the table name.")
	filePtr := flag.String("file", "", "set file to the full path of the Cassandra DDL file to process")
	goPackageNamePtr := flag.String("goPackageName", "", "set goPackageName to the desired Go package name e.g. github.com/stevef1uk/test4 (this is used to create the import statements in the generated code) ")
	primaryKeysPtr := flag.Int( "numberOfPrimaryKeys", 0, "Set to override the number of primary key fields to use for the select, defaults to that of the table definition")
	logNoDataPtr := flag.Bool("logNoData", true, "Set logNoData to false to supress logging of No Data & any error message from the select")
	outputPtr := flag.String("dirToGenerateIn", "/tmp", "set dirToGenerateIn to the full path of the directory where the output of swagger is defaults to /tmp")
	pathNamePtr := flag.String("pathNamePtr", "", "if auto patching of configure_simple.go isn't working set pathNamePtr to the full path of the directory where the output of swagger is")
	postPtr := flag.Bool("post", false, "set to true to add post method as well")
	swaggerPtr := flag.String("swaggerPath", "", "set to path of swagger executable if not on Mac/Linux ")

    //_ = swaggerPtr

    flag.Parse()

	b, err := ioutil.ReadFile(*filePtr)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	input := string(b)
	parse1 := parser.ParseText( *debugPtr, parser.Setup, parser.Reset, input )

	swagger := swagger.CreateSwagger( *debugPtr, parse1, *endPointPtr, *postPtr )

	pathName := os.Getenv("GOPATH")  + string(os.PathSeparator)  + "src" + string(os.PathSeparator) + *goPackageNamePtr
	if *pathNamePtr != "" {
		pathName = *pathNamePtr
	}
	swaggerFile := handler.CreateFile( *debugPtr , pathName , "", SWAGGER_FILE )
	err = ioutil.WriteFile(swaggerFile.Name(), []byte(swagger), 0666)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	handler.CreateCode( *debugPtr, *outputPtr, *goPackageNamePtr, parse1,  *consistencyPtr,  *endPointPtr, *primaryKeysPtr,  *allowFilteringPtr, *logNoDataPtr, *postPtr   )

	if *swaggerPtr != "" {
		os.Setenv("PATH", *swaggerPtr )
	} else {
		os.Setenv("PATH", "/usr/bin:/sbin:/usr/local/bin:/bin")
	}
	// Latest version of swagger only works for json format so convert the generated swagger file from yaml
	command := "yq"
	args := []string{"-o=json",  "./" + pathName + string(os.PathSeparator)+ SWAGGER_FILE, ">",  "/tmp/tswagger" )
	if err := exec.Command(command, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Swagger conversion from yaml to json failed")
		os.Exit(1)
	}
	command = "swagger"
	args = []string{"generate", "server", "-f", pathName + string(os.PathSeparator)+ SWAGGER_FILE, "-t", pathName }
	if err := exec.Command(command, args...).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintln(os.Stderr, "Swagger failed")
		os.Exit(1)
	}

	//tableName := handler.GetFieldName(  *debugPtr, false, parse1.TableDetails.TableName, false )
	tableName := parse1.TableDetails.TableName
	ret := handler.SpiceInHandler( false , pathName, tableName, *endPointPtr, *postPtr  )
	_ = ret
}
