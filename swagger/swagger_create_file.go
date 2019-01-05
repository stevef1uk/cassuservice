package swagger

import (
	"github.com/stevef1uk/cassuservice/parser"
	"strings"
)


func addParametersAndResponses( debug bool, output string, parseOutput  parser.ParseOutput) string {
	ret := output
	tableDetails := parseOutput.TableDetails

	for i :=0;  i < tableDetails.PkIndex; i++ {
		fieldDetails := findFieldByname( tableDetails.DbPKFields[i], tableDetails.FieldIndex, tableDetails.TableFields)
		ret = ret + `
` + "        - " + "name: " + strings.ToLower(tableDetails.DbPKFields[i]) + `
` + "          in: query" + `
` + "          description: Primary Key field in Table" + `
` + "          required: true" + `
` + "          type: " +  mapCassandraTypeToSwaggerType( false, fieldDetails.DbFieldType ) + `
` + "          format: " + mapCassandraTypeToSwaggerFormat( fieldDetails.DbFieldType )
	}

	ret = ret + `
` + "      responses:" + `
` + "        200:" + `
` + "          description: A list of records" + `
` + "          schema:" + `
` + "            type: array" + `
` + "            items:" + `
` + "              required:"

	for i :=0;  i < tableDetails.FieldIndex; i++ {
		ret = ret + `
` + "                - " + strings.ToLower( tableDetails.TableFields.DbFieldDetails[i].DbFieldName)
	}

	ret = ret + `
` + "              properties:"

	for i :=0;  i < tableDetails.FieldIndex; i++ {
		ret = ret + `
` + "                 " + strings.ToLower( tableDetails.TableFields.DbFieldDetails[i].DbFieldName) + ":"
		if  tableDetails.TableFields.DbFieldDetails[i].DbFieldMapType != "" ||
			IsFieldTypeUDT( parseOutput, tableDetails.TableFields.DbFieldDetails[i].DbFieldType ) {
			ret = ret + `
` + "                   $ref: " + `"#/definitions/` + strings.ToLower( tableDetails.TableFields.DbFieldDetails[i].DbFieldName) + `"`
		} else {
			ret = ret + `
` + "                   " + "type:" + mapCassandraTypeToSwaggerType(true, tableDetails.TableFields.DbFieldDetails[i].DbFieldType)
			if tableDetails.TableFields.DbFieldDetails[i].DbFieldCollectionType != "" {
				ret = ret + `
` + "                   items:" + `
` + "                     type: " + mapCassandraTypeToSwaggerType(true, tableDetails.TableFields.DbFieldDetails[i].DbFieldCollectionType)
			} else {
				if IsFieldaTime( tableDetails.TableFields.DbFieldDetails[i].DbFieldType ) {
					ret = ret + `
` + "                   format: date-time"
				}
			}
		}
	}

	ret = ret + `
` + "        400: " + `
` + "          description: Record not found" + `
` + "        default:" + `
` + "          description: Sorry unexpected error"


	return ret
}


// Main function to generate a string containing a swagger file

func CreateSwagger( debug bool, parseOutput parser.ParseOutput ) string {
	retSwagger := HEADER

	//Add tablename & get string
	retSwagger = retSwagger + `
  ` + "/" + strings.ToLower(parseOutput.TableDetails.TableName) + ":" + `
    get: 
      summary: Retrieve some records from the Cassandra table 
      description: Returns rows from the Cassandra table
      parameters:`

	// Add the parameters
	retSwagger = addParametersAndResponses( debug, retSwagger, parseOutput )

	return retSwagger
}


