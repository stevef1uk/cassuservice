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

func TestUDT_1(t *testing.T) {

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
  /pisp_submissions_per_id:
    get: 
      summary: Retrieve some records from the Cassandra table 
      description: Returns rows from the Cassandra table
      parameters:
        - name: submissionid
          in: query
          description: Primary Key field in Table
          required: true
          type: string
          format: string
        - name: lastupdatedat
          in: query
          description: Primary Key field in Table
          required: true
          type: string
          format: date-time
      responses:
        200:
          description: A list of records
          schema:
            type: array
            items:
              required:
                - submissionid
                - timebucket
                - debtoragent
                - debtoraccount
                - creditoragent
                - lastupdatedat
              properties:
                 submissionid:
                   type: string
                 timebucket:
                   type: string
                 debtoragent:
                   $ref: "#/definitions/debtor_agent"
                 debtoraccount:
                   $ref: "#/definitions/debtor_account"
                 creditoragent:
                   $ref: "#/definitions/creditor_agent"
                 lastupdatedat:
                   type: string
        400: 
          description: Record not found
        default:
          description: Sorry unexpected error
definitions:
  debtor_agent:
    properties:
       schemename:
         type: string
       identification:
         type: string
  debtor_account:
    properties:
       schemename:
         type: string
       identification:
         type: string
       name:
         type: string
       secondaryidentification:
         type: string
  creditor_agent:
    properties:
       schemename:
         type: string
       identification:
         type: string`


	ret := parser.ParseText( false, parser.Setup, parser.Reset, `
CREATE TYPE demo.debtor_agent (
  schemeName text,
  identification text
);


CREATE TYPE demo.debtor_account (
  schemeName text,
  identification text,
  name text,
  secondaryIdentification text
);


CREATE TYPE demo.creditor_agent (
  schemeName text,
  identification text
);
CREATE TABLE demo.pisp_submissions_per_id (
	submissionId uuid,
	timeBucket text,
	debtorAgent debtor_agent,
	debtorAccount debtor_account,
	creditorAgent creditor_agent,
	lastUpdatedAt timestamp,
	PRIMARY KEY (submissionId, lastUpdatedAt)
) WITH CLUSTERING ORDER BY (lastUpdatedAt DESC)
`)

	ret1 := CreateSwagger( true, ret )
	if expectedOutput != ret1 {

		if len(expectedOutput) != len(ret1) {
			t.Errorf("Expected length %d, actual, %d", len(expectedOutput), len(ret1) )
		}
		/*for i, _ := range expectedOutput {
			if (expectedOutput[i] != ret1[i] ) {
				t.Errorf("Difference at %d, got %c expected %c", i, expectedOutput[i], ret1[i] )
			}
		}*/
		t.Errorf("Swagger output wrong got:%s: want:%s:", ret1, expectedOutput )
	}

}