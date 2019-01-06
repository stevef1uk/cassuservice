package swagger

import (
	"github.com/stevef1uk/cassuservice/parser"
	"testing"
)


func TestUDT_Colections(t *testing.T) {

	expectedOutput := `swagger: '2.0'
info:
  version: 1.0.0
  title: Simple API
  description: A generated file representing a Cassandra Table definition
schemes:
  - http
host: localhost
basePath: /v1
paths:
  /employee:
    get: 
      summary: Retrieve some records from the Cassandra table 
      description: Returns rows from the Cassandra table
      parameters:
        - name: id
          in: query
          description: Primary Key field in Table
          required: true
          type: integer
          format: int32
      responses:
        200:
          description: A list of records
          schema:
            type: array
            items:
              required:
                - id
                - address_map
                - address_list
                - address_set
                - name
              properties:
                 id:
                   type: integer
                 address_map:
                   $ref: "#/definitions/address_map"
                 address_list:
                   $ref: "#/definitions/address_list"
                 address_set:
                   $ref: "#/definitions/address_set"
                 name:
                   type: string
        400: 
          description: Record not found
        default:
          description: Sorry unexpected error
definitions:
  city:
    properties:
       id:
         type: integer
       citycode:
         type: string
       cityname:
         type: string
  address_map:
      additionalProperties:
         $ref: "#/definitions/city"
  address_list:
      type: array
      items:
         $ref: "#/definitions/city"
  address_set:
      type: array
      items:
         $ref: "#/definitions/city"`

	ret := parser.ParseText( false, parser.Setup, parser.Reset, `
       CREATE TYPE demo.city (
	       id int,
	       citycode text,
	       cityname text
	   );

	   		CREATE TABLE demo.employee (
	       id int PRIMARY KEY,
	       address_map map<text, frozen <city>>,
	       address_list list<frozen<city>>,
	       address_set set<frozen<city>>,
	       name text
	   ) WITH CLUSTERING ORDER BY (name DESC);
	   	` )

	ret1 := CreateSwagger( true, ret )

	if expectedOutput != ret1 {

		if len(expectedOutput) != len(ret1) {
			t.Errorf("Expected length %d, actual, %d", len(expectedOutput), len(ret1) )
		}
		for i, _ := range expectedOutput {
			if (expectedOutput[i] != ret1[i] ) {
				t.Errorf("Difference at %d, got %c expected %c", i, expectedOutput[i], ret1[i] )
			}
		}
		t.Errorf("Swagger output wrong got:%s: want:%s:", ret1, expectedOutput )
	}
}