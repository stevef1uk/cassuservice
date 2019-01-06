package swagger

import (
	"github.com/stevef1uk/cassuservice/parser"
	"strings"
)

// returns the outout plus a flag indicating if definitions are required
func addDefinitions( debug bool, output string, parseOutput  parser.ParseOutput ) (string, bool) {
	ret := output
	defsRequired := false;
	tableDetails := parseOutput.TableDetails

	for i :=0;  i < tableDetails.TableFields.FieldIndex; i++ {
		if  tableDetails.TableFields.DbFieldDetails[i].DbFieldMapType != "" ||
			IsFieldTypeUDT( parseOutput, tableDetails.TableFields.DbFieldDetails[i].DbFieldType ) {
			defsRequired = true
		}

	}
	if defsRequired {
		ret = ret + `
` + "definitions:"
	}

	return ret, defsRequired
}

// Add the definition details for map types
func addMaps( debug bool, output string, parseOutput  parser.ParseOutput ) string {
	ret := output
	tableDetails := parseOutput.TableDetails

	for i :=0;  i < tableDetails.TableFields.FieldIndex; i++ {
		if  tableDetails.TableFields.DbFieldDetails[i].DbFieldMapType != "" {
			ret = ret + `
` + "  " +  strings.ToLower( tableDetails.TableFields.DbFieldDetails[i].DbFieldName) + ":" + `
` + "    additionalProperties:" + `
` + "      type: " +  mapCassandraTypeToSwaggerType(true, tableDetails.TableFields.DbFieldDetails[i].DbFieldMapType)

		}
	}

	return ret
}

// Add field details
func addFieldDetails( debug bool, spaces string, output string, tableDetails  parser.AllFieldDetails, parseOutput  parser.ParseOutput ) string {
	ret := output

	for i :=0;  i < tableDetails.FieldIndex; i++ {
		ret = ret + `
` + spaces + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldName) + ":"
		if  tableDetails.DbFieldDetails[i].DbFieldMapType != "" ||
			IsFieldTypeUDT( parseOutput, tableDetails.DbFieldDetails[i].DbFieldType ) {
			ret = ret + `
` + spaces + "  $ref: " + `"#/definitions/` + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldName) + `"`
		} else {
			if tableDetails.DbFieldDetails[i].DbFieldCollectionType != "" {
				if IsFieldTypeUDT(parseOutput, tableDetails.DbFieldDetails[i].DbFieldCollectionType)  {
					ret = ret + `
` + spaces + "  $ref: " + `"#/definitions/` + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldName) + `"`
				} else {
					ret = ret + `
` + spaces + "  type: " + mapCassandraTypeToSwaggerType(true, tableDetails.DbFieldDetails[i].DbFieldType)
					if tableDetails.DbFieldDetails[i].DbFieldCollectionType != "" {

						ret = ret + `
` + spaces + "  items:" + `
` + spaces + "    type: " + mapCassandraTypeToSwaggerType(true, tableDetails.DbFieldDetails[i].DbFieldCollectionType)
					} else {
						if IsFieldaTime(tableDetails.DbFieldDetails[i].DbFieldType) {
							ret = ret + `
` + spaces + "format: date-time"
						}
					}
				}
			} else {
				ret = ret + `
` + spaces + "  type: " + mapCassandraTypeToSwaggerType(true, tableDetails.DbFieldDetails[i].DbFieldType)
			}
		}
	}

	return ret
}


// Add the UDT details in the definitions sections
func addUDTs( debug bool, output string, parseOutput  parser.ParseOutput ) string {
	ret := output


	for i :=0;  i < parseOutput.TypeIndex; i++ {
		tableDetails := parseOutput.TypeDetails[i]
		ret = ret + `
` + "  " +  strings.ToLower( tableDetails.TypeName) + ":" + `
`
		ret = ret + "    properties:"

		ret = ret + addFieldDetails( debug, "       " , ret, tableDetails.TypeFields, parseOutput  )

	}

	return ret
}


func addParametersAndResponses( debug bool, output string, parseOutput  parser.ParseOutput) string {
	ret := output
	tableDetails := parseOutput.TableDetails

	for i :=0;  i < tableDetails.PkIndex; i++ {
		fieldDetails := findFieldByname( tableDetails.DbPKFields[i], tableDetails.TableFields.FieldIndex, tableDetails.TableFields)
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

	for i :=0;  i < tableDetails.TableFields.FieldIndex; i++ {
		ret = ret + `
` + "                - " + strings.ToLower( tableDetails.TableFields.DbFieldDetails[i].DbFieldName)
	}

	ret = ret + `
` + "              properties:"

	ret = ret + addFieldDetails( debug, "                 " , ret, tableDetails.TableFields, parseOutput  )

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
	retSwagger, haveDefs := addDefinitions( debug, retSwagger, parseOutput  )
	if haveDefs {
		retSwagger = addUDTs( debug, retSwagger, parseOutput )
		retSwagger = addMaps( debug, retSwagger, parseOutput )
	}

	return retSwagger
}


