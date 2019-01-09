package handler

import (
	"fmt"
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
func CreateFile( debug bool, codeBasePath string, generateDir string ) *os.File {

	// Create the directory if not already there
	fulldirName := generateDir  + codeBasePath
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

// Function that renames fields to match that performed for some reason by go-swagger in its generated framework code
func ReviseFieldName ( debug bool, fieldName string, dontUpdate bool ) string {

	var ret string = ""
	if debug { fmt.Printf("ReviseFieldName entry field  = %s, len = %d\n ",fieldName, len(fieldName) ) }

	if dontUpdate {
		ret = fieldName
		if debug { fmt.Printf("ReviseFieldName told not to update\n ") }
	} else {
		runes := []rune(fieldName[:])
		last := len( runes  ) - 1

		for i := 0; i < last ; i++   {
				//if debug { fmt.Printf("ReviseFieldName i   = %d\n ",i ) }
				if ! ( ( i == 0 ) || ( i == last -1 ) ) {
					continue;
				}
				//if debug { fmt.Printf("ReviseFieldName [0] = %q, [1] = %q\n ", runes[i], runes[i+1]) }
				if (runes[i] == rune('i') || runes[i] == rune('I')) && (runes[i+1] == rune('d') || runes[i+1] == rune('D')) {
					if debug { fmt.Printf("ReviseFieldName match at i= %d\n ", i)}
					runes[i] = rune('I')
					runes[i+1] = rune('D')
				}
			}
			ret = string(runes)
		}

	if debug {fmt.Printf("ReviseFieldName returning field  = %s\n ", ret)}
	return ret
}


func mapCassandraTypeToGoType( debug bool, field parser.FieldDetails, collectionofUDT bool, smallInt bool, smallFloat bool, istimeuuid bool  ) string {
	var text string = ""
	switch strings.ToLower(field.DbFieldType) {
	case "int":
		if (smallInt) {
			text = "int" // Reflection of Cassandra Type is int not int64
		} else {
			text = "int64"
		}
	case "uuid":
		text = "string"
	case "date":
		text = "time.Time"
	case "timeuuid":
		if (istimeuuid) {
			text = "gocql.UUID"
		} else {
			text = "time.Time"
		}
	case "timestamp":
		text = "time.Time"
	case "varint":
		text = "int64"
	case "boolean":
		text = "bool"
	case "bigint":
		text = "int64"
	case "counter":
		text = "int64"
	case "decimal":
		text = "*inf.Dec" // this is in the gopkg.in/inf.v0 package
	case "float":
		if (smallFloat) {
			text = "float32" // Reflection of Cassandra Type is int not int64
		} else {
			text = "float64"
		}
	case "double":
		text = "float64"
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
	case "list":
		if (collectionofUDT) {
			text = ""
		} else {
			text = "[]"
		}
	case "set":
		if (collectionofUDT) {
			text = ""
		} else {
			text = "[]"
		}
	case "map":
		if (collectionofUDT) {
			text = "map[string]string" // This is the type required by go-swagger
		} else {
			text = "" // go-swagger will have created a type for the map
		}

	default:
		fmt.Printf("Field type not recognised %q = ", field)
		panic(1)
	}

	if debug { fmt.Printf("mapCassandraTypeToGoType returning %s from field %q\n", text, field ) }
	return text
}

// Function to return a temporary variable based on string
var counter int = 0
func createTempVar ( fieldName string ) string {
	ret := TEMP_VAR_PREFIX + fieldName + "_" + strconv.Itoa( counter)
	counter = counter + 1
	return ret
}

// Helper function for createSelect to handle array fields

func setUpArrayTypes(  debug bool, output string, field parser.FieldDetails,  dontUpdate bool ) string {
	ret := output
	tmpVar := createTempVar( field.DbFieldName)

	if  swagger.IsFieldaTime( field.DbFieldType ) {
		ret = ret + `
        ` + tmpVar + " = strfmt.NewDateTime().String()" + `
        ` + "_ = " + tmpVar + `
		` + strings.ToLower(field.DbFieldName) + "= " + RAWRESULT + `["` + strings.ToLower(field.DbFieldName) + `"].([]` +
		    mapCassandraTypeToGoType( true, field,false,   false, false, false) +  `)
		` + "payLoad." + ReviseFieldName( debug, field.DbFieldName, dontUpdate) + " = make([] string, len(" + strings.ToLower(field.DbFieldName) + ") )" + `
		for i := 0; i < len(` + strings.ToLower(field.DbFieldName) + `); i++ {
			payLoad.` + ReviseFieldName( debug, field.DbFieldType, dontUpdate) + "[i] = " + strings.ToLower(field.DbFieldType) + "[i].String()" + `
		}`
	} else {
		if ( strings.ToLower(field.DbFieldType) == "decimal") {
			ret = ret + `
    payLoad.` + ReviseFieldName(debug, field.DbFieldType, dontUpdate) + " = make([]float64, len(" + strings.ToLower(field.DbFieldType) + ") )" + `
    for i := 0; i < len(` + strings.ToLower(field.DbFieldType) + `); i++ {
        ` + "mytmpdecjf123_" + strings.ToLower(field.DbFieldType) + ", err := strconv.ParseFloat( " + strings.ToLower(field.DbFieldType) + "[i].String(), 64 )" + `
        if ( err != nil ) {
            log.Println("error parsing decimal value for field",` + field.DbFieldType + `)
        }
` + `
        payLoad.` + ReviseFieldName( debug, field.DbFieldName, dontUpdate) + "[i] = " + "mytmpdecjf123_" + strings.ToLower(field.DbFieldName) + `
    }`
		} else {
			ret = ret + `
		` + strings.ToLower(field.DbFieldName) + "= " + RAWRESULT + `["` + strings.ToLower(field.DbFieldName) + `"].([]` + mapCassandraTypeToGoType( true, field,false,   false, false, false) + `)`
			ret = ret + `
		` + "payLoad." + ReviseFieldName(debug, field.DbFieldName, dontUpdate) + " = make([]" + mapCassandraTypeToGoType( true, field,false,   false, false, false) + ",len(" + strings.ToLower(field.DbFieldName) + ") )" + `
		for i := 0; i < len(` + strings.ToLower(field.DbFieldName) + `); i++ {
			payLoad.` + ReviseFieldName(debug,field.DbFieldName, dontUpdate) + "[i] = " + mapCassandraTypeToGoType( true, field,false,   false, false, false) + "(" + strings.ToLower(field.DbFieldName) + "[i])" + `
		}`
		}
	}
	if debug { fmt.Printf("setUpArrayTypes returning %s\n", ret ) }
	return ret
}
