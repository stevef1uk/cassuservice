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

	ret1 := CreateSwagger( true, ret, "" )

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

	ret1 := CreateSwagger( true, ret, "" )
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

func TestTable_1(t *testing.T) {

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
  /accounts4:
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
        - name: name
          in: query
          description: Primary Key field in Table
          required: true
          type: string
          format: string
        - name: time1
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
                - id
                - name
                - ascii1
                - bint1
                - blob1
                - bool1
                - counter1
                - dec1
                - double1
                - flt1
                - inet1
                - int1
                - text1
                - time1
                - time2
                - uuid1
                - varchar1
                - events
                - mylist
                - myset
                - mymap
              properties:
                 id:
                   type: integer
                 name:
                   type: string
                 ascii1:
                   type: string
                 bint1:
                   type: integer
                 blob1:
                   type: string
                 bool1:
                   type: boolean
                 counter1:
                   type: integer
                 dec1:
                   type: number
                 double1:
                   type: number
                 flt1:
                   type: number
                 inet1:
                   type: string
                 int1:
                   type: integer
                 text1:
                   type: string
                 time1:
                   type: string
                 time2:
                   type: string
                 uuid1:
                   type: string
                 varchar1:
                   type: string
                 events:
                   type: array
                   items:
                     type: integer
                 mylist:
                   type: array
                   items:
                     type: number
                 myset:
                   type: array
                   items:
                     type: string
                 mymap:
                   $ref: "#/definitions/mymap"
        400: 
          description: Record not found
        default:
          description: Sorry unexpected error
definitions:
  mymap:
      additionalProperties:
        type: string`

	ret := parser.ParseText( false, parser.Setup, parser.Reset, `
CREATE TABLE demo.accounts4 (
    id int,
    name text,
    ascii1 ascii,
    bint1 bigint,
    blob1 blob,
    bool1 boolean,
    counter1 counter,
    dec1 decimal,
    double1 double,
    flt1 float,
    inet1 inet,
    int1 int,
    text1 text,
    time1 timestamp,
    time2 timeuuid,
    uuid1 uuid,
    varchar1 varchar,
    events set<int>,
    mylist list<float>,
    myset set<text>,
    mymap  map<int, text>,
    PRIMARY KEY (id, name, time1)
) WITH CLUSTERING ORDER BY (name ASC)
` )

	ret1 := CreateSwagger( true, ret, "" )
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

func TestSimple1(t *testing.T) {

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
  /accounts:
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
        - name: name
          in: query
          description: Primary Key field in Table
          required: true
          type: string
          format: string
      responses:
        200:
          description: A list of records
          schema:
            type: array
            items:
              required:
                - id
                - name
              properties:
                 id:
                   type: integer
                 name:
                   type: string
        400: 
          description: Record not found
        default:
          description: Sorry unexpected error`

	ret := parser.ParseText( false, parser.Setup, parser.Reset, `
CREATE TABLE demo.accounts (
    id int,
    name text,
    PRIMARY KEY (id, name)
) WITH CLUSTERING ORDER BY (name ASC)
` )

	ret1 := CreateSwagger( true, ret, "" )
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

func TestAlmostSimple(t *testing.T) {

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
        - name: mediate
          in: query
          description: Primary Key field in Table
          required: true
          type: string
          format: date-time
        - name: second_ts
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
                - id
                - address_set
                - my_list
                - name
                - mediate
                - second_ts
              properties:
                 id:
                   type: integer
                 address_set:
                   $ref: "#/definitions/address_set"
                 my_list:
                   $ref: "#/definitions/my_list"
                 name:
                   type: string
                 mediate:
                   type: string
                 second_ts:
                   type: string
        400: 
          description: Record not found
        default:
          description: Sorry unexpected error
definitions:
  simple:
    properties:
       dummy:
         type: string
  city:
    properties:
       id:
         type: integer
       citycode:
         type: string
       cityname:
         type: string
       test_int:
         type: integer
       lastupdatedat:
         type: string
       myfloat:
         type: number
  address_set:
      type: array
      items:
         $ref: "#/definitions/city"
  my_list:
      type: array
      items:
         $ref: "#/definitions/simple"`

	ret := parser.ParseText( false, parser.Setup, parser.Reset, `
CREATE TYPE demo.simple (
   dummy text
);

CREATE TYPE demo.city (
    id int,
    citycode text,
    cityname text,
    test_int int,
    lastUpdatedAt TIMESTAMP,
    myfloat float
);

CREATE TABLE demo.employee (
    id int,
    address_set set<frozen<city>>,
    my_List list<frozen<simple>>,
    name text,
    mediate TIMESTAMP,
    second_ts timestamp,
   PRIMARY KEY (id, mediate, second_ts )
 ) WITH CLUSTERING ORDER BY (mediate ASC, second_ts ASC)
` )

	ret1 := CreateSwagger( true, ret, "" )
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


func TestComplexUDT(t *testing.T) {

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
  /testing:
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
        - name: mediate
          in: query
          description: Primary Key field in Table
          required: true
          type: string
          format: date-time
        - name: second_ts
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
                - id
                - address_set
                - my_list
                - name
                - mediate
                - second_ts
                - tevents
                - tmylist
                - tmymap
              properties:
                 id:
                   type: integer
                 address_set:
                   $ref: "#/definitions/address_set"
                 my_list:
                   $ref: "#/definitions/my_list"
                 name:
                   type: string
                 mediate:
                   type: string
                 second_ts:
                   type: string
                 tevents:
                   type: array
                   items:
                     type: integer
                 tmylist:
                   type: array
                   items:
                     type: number
                 tmymap:
                   $ref: "#/definitions/tmymap"
        400: 
          description: Record not found
        default:
          description: Sorry unexpected error
definitions:
  simple:
    properties:
       dummy:
         type: string
  city:
    properties:
       id:
         type: integer
       citycode:
         type: string
       cityname:
         type: string
       test_int:
         type: integer
       lastupdatedat:
         type: string
       myfloat:
         type: number
       events:
         type: array
         items:
           type: integer
       mymap:
         $ref: "#/definitions/city_mymap"
       address_list:
         $ref: "#/definitions/city_address_list"
  tmymap:
      additionalProperties:
        type: string
  city_mymap:
      additionalProperties:
        type: string
  address_set:
      type: array
      items:
         $ref: "#/definitions/city"
  my_list:
      type: array
      items:
         $ref: "#/definitions/simple"
  city_address_list:
      type: array
      items:
         $ref: "#/definitions/simple"`

	ret := parser.ParseText( false, parser.Setup, parser.Reset, `
CREATE TYPE demo.simple (
   dummy text
);

CREATE TYPE demo.city (
    id int,
    citycode text,
    cityname text,
    test_int int,
    lastUpdatedAt TIMESTAMP,
    myfloat float
    events set<int>,
    mymap  map<int, text>
    address_list set<frozen<simple>>,
);

CREATE TABLE demo.employee (
    id int,
    address_set set<frozen<city>>,
    my_List list<frozen<simple>>,
    name text,
    mediate TIMESTAMP,
    second_ts timestamp,
    tevents set<int>,
    tmylist list<float>
    tmymap  map<int, text>
   PRIMARY KEY (id, mediate, second_ts )
 ) WITH CLUSTERING ORDER BY (mediate ASC, second_ts ASC)
` )
	ret1 := CreateSwagger( true, ret, "testing" )
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