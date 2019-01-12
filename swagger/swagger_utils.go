package swagger

import (
	"github.com/stevef1uk/cassuservice/parser"
	"strings"
	"log"
)

// Utility function to identify if the passed fieldType is a User Defined Type
func IsFieldTypeUDT( typeDetails parser.ParseOutput, fieldType string ) bool {
	ret := false
	for _, v := range typeDetails.TypeDetails {
		if v.TypeName == fieldType {
			ret = true
			break
		}
	}

	return ret
}

// Simple function to return true if the string passed is a Cassandra time feld
func IsFieldTypeATime(  fType string ) bool {
	ret := false
	if fType == TIMESTAMP || fType == DATE || fType == TIMEUUID {
		ret = true;
	}
	return ret
}

func findFieldByname( fieldName string, noFields int,  fields parser.AllFieldDetails )  parser.FieldDetails {
	var ret parser.FieldDetails
	for i := 0; i < noFields;  i++ {
		if  fields.DbFieldDetails[i].DbFieldName== fieldName {
			ret = fields.DbFieldDetails[i]
			break
		}
	}
	return ret
}

func mapCassandraTypeToSwaggerType( checkType bool, fieldType string  ) string {
	var text string = ""
	switch strings.ToLower(fieldType) {
	case "int":
		text = "integer"
	case "varint":
		text = "integer"
	case "date":
		text = "string"
	case "bigint":
		text = "integer"
	case "uuid":
		text = "string"
	case "counter":
		text = "integer"
	case "decimal":
		text = "number"
	case "float":
		text = "number"
	case "double":
		text = "number"
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
	case "boolean":
		text = "boolean"
	case "list":
		text = "array"
	case "map":
		text = "array"
	case "set":
		text = "array"
	case "timestamp":
		text = "string"
	case "timeuuid":
		text = "string"
	default:
		//log.Println("map func checktype = ",checktype )
		if checkType {
			log.Fatal("Data type not supported in mapCassandraTypeToSwaggerType = ", fieldType )
		}
		text = fieldType
		//log.Fatal("Data type not supported in parse mapType = ", fieldType )
	}

	return text
}

func mapCassandraTypeToSwaggerFormat( fieldType string  ) string {
	var text string = ""
	switch strings.ToLower(fieldType) {
	case "int":
		text = "int32"
	case "uuid":
		text = "string"
	case "varint":
		text = "int64"
	case "date":
		text = "date-time"
	case "timestamp":
		text = "date-time"
	case "boolean":
		text = "boolean"
	case "bigint":
		text = "int64"
	case "counter":
		text = "int64"
	case "decimal":
		text = "float"
	case "float":
		text = "numberfloat"
	case "double":
		text = "numberdouble"
	case "text":
		text = "string"
	case "ascii":
		text = "string"
	case "varchar":
		text = "string"
	case "blob":
		text = "string"
	case "inet":
		text = "string"
	case "timeuuid":
		text = "date-time"
	default:
		log.Fatal("Field type not supported in mapCassandraTypeToSwaggerFormat = ", fieldType)
	}
	return text

}