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
			(tableDetails.TableFields.DbFieldDetails[i].DbFieldCollectionType != "" && IsFieldTypeUDT(parseOutput, tableDetails.TableFields.DbFieldDetails[i].DbFieldCollectionType) ) ||
			IsFieldTypeUDT( parseOutput, tableDetails.TableFields.DbFieldDetails[i].DbFieldType ) {
			defsRequired = true
		}

	}
	if ! defsRequired {
		for i :=0;  i < parseOutput.TypeIndex;  i++ {
			thisType := parseOutput.TypeDetails[i]
			for j :=0;  j < thisType.TypeFields.FieldIndex; j++ {
				v := thisType.TypeFields.DbFieldDetails
				if v[j].DbFieldMapType != "" || v[j].DbFieldCollectionType != "" || IsFieldTypeUDT(parseOutput, v[j].DbFieldType) {
					defsRequired = true
				}
			}

		}
	}
	if defsRequired {
		ret = ret + `
` + "definitions:"
	}

	return ret, defsRequired
}


// Add the definition details for map types, ensure that these are delineated between tables & UDTs
func addMap( debug bool, parserOutput  parser.ParseOutput, fieldDetails parser.FieldDetails, inType bool, typeIndex int) string {
	ret := ""
	if  fieldDetails.DbFieldMapType != "" {
		if inType {
			ret = ret + `
` + "  " + strings.ToLower(parserOutput.TypeDetails[typeIndex].TypeName) + "_" + strings.ToLower(fieldDetails.DbFieldName)
		} else {
			ret = ret + `
` + "  " + strings.ToLower(fieldDetails.DbFieldName)
		}
		ret = ret +  ":" + `
` + "      additionalProperties:" + `
`

		if IsFieldTypeUDT(parserOutput, fieldDetails.DbFieldMapType) {
			ret = ret + "         $ref: " +  `"#/definitions/` + strings.ToLower(fieldDetails.DbFieldMapType ) + `"`
		} else {
			ret = ret + "        type: " + mapCassandraTypeToSwaggerType(true, fieldDetails.DbFieldMapType)
		}
	}
	return ret
}


// Add the definition details for map types
func addMaps( debug bool, output string, parserOutput  parser.ParseOutput ) string {
	ret := output
	tableDetails := parserOutput.TableDetails

	for i :=0;  i < tableDetails.TableFields.FieldIndex; i++ {
		ret = ret + addMap( debug, parserOutput, tableDetails.TableFields.DbFieldDetails[i], false, 0 )
	}

	for i := 0; i < parserOutput.TypeIndex; i++ {
		v := parserOutput.TypeDetails[i]
		for j := 0; j < v.TypeFields.FieldIndex ; j++ {
			ret = ret + addMap( debug, parserOutput, v.TypeFields.DbFieldDetails[j], true, i )
		}
	}

	return ret
}


// Add the definition details for list & set types when they need adding to the definitions only
func addCollectionType( debug bool, parseOutput  parser.ParseOutput, fieldDetails parser.FieldDetails, inType bool, typeIndex int ) string {
	ret := `
` + "  "
	if inType {
		ret = ret + strings.ToLower(parseOutput.TypeDetails[typeIndex].TypeName) + "_" +  strings.ToLower(fieldDetails.DbFieldName)
	} else
	{
		ret = ret + strings.ToLower(fieldDetails.DbFieldName)
	}
	ret = ret + ":" + `
` + "      type: array" + `
` + "      items:" + `
` + "         $ref: " + `"#/definitions/` + strings.ToLower(fieldDetails.DbFieldCollectionType) + `"`

	return ret
}


// Add the definition details for list & set types when they need adding to the definitions only
func addCollectionTypes( debug bool, output string, parserOutput  parser.ParseOutput ) string {
	ret := output
	tableDetails := parserOutput.TableDetails

	for i := 0;  i < tableDetails.TableFields.FieldIndex; i++ {
		if  tableDetails.TableFields.DbFieldDetails[i].DbFieldCollectionType != "" && tableDetails.TableFields.DbFieldDetails[i].DbFieldMapType == "" {
			if IsFieldTypeUDT(parserOutput, tableDetails.TableFields.DbFieldDetails[i].DbFieldCollectionType) {
				ret = ret + addCollectionType(debug, parserOutput, tableDetails.TableFields.DbFieldDetails[i], false, 0)
			}
		}
	}

	for i := 0; i < parserOutput.TypeIndex; i++ {
		v := parserOutput.TypeDetails[i]
		for j := 0; j < v.TypeFields.FieldIndex ; j++ {
			if v.TypeFields.DbFieldDetails[j].DbFieldCollectionType != "" && v.TypeFields.DbFieldDetails[j].DbFieldMapType == "" {
				if IsFieldTypeUDT(parserOutput, v.TypeFields.DbFieldDetails[j].DbFieldCollectionType ) {
					ret = ret + addCollectionType(debug, parserOutput, v.TypeFields.DbFieldDetails[j], true, i)
				}
			}
		}

	}

	return ret
}

// Add field details
func addFieldDetails( debug bool, spaces string, output string, parseOutput  parser.ParseOutput, tableDetails  parser.AllFieldDetails, inType bool, typeIndex int ) string {
	ret := output

	for i :=0;  i < tableDetails.FieldIndex; i++ {
		ret = ret + `
` + spaces + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldName) + ":"
		if  tableDetails.DbFieldDetails[i].DbFieldMapType != "" ||
			IsFieldTypeUDT( parseOutput, tableDetails.DbFieldDetails[i].DbFieldType ) {
			if tableDetails.DbFieldDetails[i].DbFieldMapType != "" {
				if inType {
					ret = ret + `
` + spaces + "  $ref: " + `"#/definitions/` + strings.ToLower(parseOutput.TypeDetails[typeIndex].TypeName) + "_" + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldName) + `"`
				} else {
					ret = ret + `
` + spaces + "  $ref: " + `"#/definitions/` + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldName) + `"`
				}

			} else {
				if inType {
					if tableDetails.DbFieldDetails[i].DbFieldCollectionType != "" {
						ret = ret + `
` + spaces + "  $ref: " + `"#/definitions/` + strings.ToLower(parseOutput.TypeDetails[typeIndex].TypeName) + "_" + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldType) + `"`
					} else {
						ret = ret + `
` + spaces + "  $ref: " + `"#/definitions/` + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldType) + `"`
					}

				} else {
					ret = ret + `
` + spaces + "  $ref: " + `"#/definitions/` + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldType) + `"`
				}
			}

		} else {
			if tableDetails.DbFieldDetails[i].DbFieldCollectionType != "" {
				if IsFieldTypeUDT(parseOutput, tableDetails.DbFieldDetails[i].DbFieldCollectionType)  {
					if inType {
						ret = ret + `
` + spaces + "  $ref: " + `"#/definitions/` + strings.ToLower(parseOutput.TypeDetails[typeIndex].TypeName) + "_" + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldName) + `"`
					} else {
						ret = ret + `
` + spaces + "  $ref: " + `"#/definitions/` + strings.ToLower( tableDetails.DbFieldDetails[i].DbFieldName) + `"`
					}
				} else {
					ret = ret + `
` + spaces + "  type: " + mapCassandraTypeToSwaggerType(true, tableDetails.DbFieldDetails[i].DbFieldType)
					if tableDetails.DbFieldDetails[i].DbFieldCollectionType != "" {

						ret = ret + `
` + spaces + "  items:" + `
` + spaces + "    type: " + mapCassandraTypeToSwaggerType(true, tableDetails.DbFieldDetails[i].DbFieldCollectionType)
					} else {
						if IsFieldTypeATime(tableDetails.DbFieldDetails[i].DbFieldType) {
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
		ret = addFieldDetails( debug, "       " , ret, parseOutput, tableDetails.TypeFields,true, i   )
	}

	return ret
}


func addParametersAndResponses( debug bool, output string, parseOutput  parser.ParseOutput) string {
	ret := output
	tableDetails := parseOutput.TableDetails

	for i :=0;  i < tableDetails.PkIndex; i++ {
		fieldDetails := FindFieldByname( tableDetails.DbPKFields[i], tableDetails.TableFields.FieldIndex, tableDetails.TableFields)
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

	ret = addFieldDetails( debug, "                 " , ret, parseOutput, tableDetails.TableFields, false, 0  )

	ret = ret + `
` + "        400: " + `
` + "          description: Record not found" + `
` + "        default:" + `
` + "          description: Sorry unexpected error"


	return ret
}

// FUnction to add the schema for the Post record
func addPostFields( debug bool, haveDefs bool , parseOutput  parser.ParseOutput, tableName string  ) string {
	ret, ret1 := "", ""
	tableDetails := parseOutput.TableDetails

	if ! haveDefs {
		ret = "definitions:"
	}

    ret = ret + `
  ` + tableName + `:
    properties:`
	for i :=0;  i < tableDetails.PkIndex; i++ {
		ret1 = addFieldDetails(debug, "      ", "", parseOutput, tableDetails.TableFields, false, 0)
	}
	return ret + ret1
}


// Function to create the swagger string for a POST operation to support inserts
func CreateSwaggerPost( debug bool, output string, parseOutput parser.ParseOutput ) (string, string)  {
	retSwagger := output

	tableName := CapitaliseSplitTableName(debug, parseOutput.TableDetails.TableName)

	//Add tablename & get string
	retSwagger = output + `
    post: 
      tags:
      - "` + tableName + `" 
      summary: Add a new record to the Cassandra table 
      description: Adds or updates a row in the Cassandra table
      operationId: add` + tableName + `
      consumes:
        - application/json
      produces:
        - application/json
      parameters:
        - in: body
          name: body
          description: The fields of the table that needs to be populates in JSON form
          required: true
          schema:
            $ref: '#/definitions/` + tableName + `'
      responses:
        '201':
          description: Record created
        '405':
          description: Invalid input
`
	return retSwagger, tableName
}

// Main function to generate a string containing a swagger file

func CreateSwagger( debug bool, parseOutput parser.ParseOutput, endPointOveride string, addPost bool  ) string {
	retSwagger := HEADER

	tableName := strings.ToLower(parseOutput.TableDetails.TableName)
	if endPointOveride != "" {
		tableName = strings.ToLower(endPointOveride)
	}

	//Add tablename & get string
	retSwagger = retSwagger + `
  ` + "/" + tableName + ":" + `
    get: 
      summary: Retrieve some records from the Cassandra table 
      description: Returns rows from the Cassandra table
      parameters:`

	// Add the parameters
	retSwagger = addParametersAndResponses( debug, retSwagger, parseOutput )
	if addPost {
		retSwagger, tableName = CreateSwaggerPost(debug, retSwagger, parseOutput)
	}

	retSwagger, haveDefs := addDefinitions( debug, retSwagger, parseOutput  )
	if haveDefs {
		retSwagger = addUDTs( debug, retSwagger, parseOutput )
		retSwagger = addMaps( debug, retSwagger, parseOutput )
		retSwagger = addCollectionTypes( debug, retSwagger, parseOutput )
	}
	if addPost {
		retSwagger = retSwagger + addPostFields(debug, haveDefs, parseOutput, tableName)
	}

	return retSwagger
}



