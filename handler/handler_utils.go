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

	var ret string = ""
	if debug { fmt.Printf("CapitaliseSplitFieldName entry field  = %s, len = %d\n ",fieldName, len(fieldName) ) }

	if dontUpdate {
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


func mapCassandraTypeToGoType( debug bool, fieldType string, collectionofUDT bool, smallInt bool, smallFloat bool  ) string {
	var text string = ""
	switch strings.ToLower(fieldType) {
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
		if ( fieldType == swagger.TIMEUUID ) {
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
		fmt.Printf("Field type not recognised %q = ", fieldType)
		panic(1)
	}

	if debug { fmt.Printf("mapCassandraTypeToGoType returning %s from field %q\n", text, fieldType ) }
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