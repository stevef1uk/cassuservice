package main

import (
	"flag"
	"github.com/stevef1uk/cassuservice/handler"
	"github.com/stevef1uk/cassuservice/parser"
	"github.com/stevef1uk/cassuservice/swagger"
)

// Set debug flag to true to output logging information to stdout
var debug = false

func main() {

	debugPtr := flag.Bool("debug", false, "Set flag to true to write trace information to stdout")

	flag.Parse()
	debug = *debugPtr
	//parser.Setup()

// //mylist list<float>,
		ret := parser.ParseText( debug, parser.Setup, parser.Reset,`
CREATE TABLE demo.accounts4 (

    id int,
    name text,
    ascii1 ascii,
    bint1 bigint,
    blob1 blob,
    bool1 boolean,
    dec1 decimal,
    double1 double,
    flt1 float,
    inet1 inet,
    int1 int,
    text1 text,
    time1 timestamp,
    time2 timeuuid,
    mydate1 date,
    uuid1 uuid,
    varchar1 varchar,
    events set<int>,
    mylist list<float>,
    myset set<text>,
    adec list<decimal>,
    PRIMARY KEY (id, name, time1)
)WITH CLUSTERING ORDER BY (name ASC)`)

	ret1 := swagger.CreateSwagger( false, ret )
	println("Swagger=\n")
	println(ret1)

	ret2 := handler.Capitiseid( false, "Id", false )
	println(ret2)
	_ = ret2;
/*
	 ret := parser.ParseText( debug, parser.Setup, parser.Reset, `
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

	ret1 := swagger.CreateSwagger( true, ret )
	println("Swagger=\n")
	println(ret1)
*/

/*
	ret := parser.ParseText( debug, parser.Setup, parser.Reset, `
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

` )
	ret1 := swagger.CreateSwagger( true, ret )
	println("Swagger=\n")
	println(ret1)

*/

/*
	ret :=  parser.ParseText(debug, parser.Setup, parser.Reset,`
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
) WITH CLUSTERING ORDER BY (name ASC)`)

	ret1 := swagger.CreateSwagger( true, ret )
	println("Swagger=\n")
	println(ret1)

*/

/*
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
    mymap  map<text, text>
    address_list set<frozen<simple>>,
);

CREATE TABLE demo.employee (
    id int,
    address_set set<frozen<city>>,
    my_List list<frozen<simple>>,
    name text,
    mediate TIMESTAMP,
    second_ts date,
    tevents set<int>,
    tmylist list<float>
    tmymap  map<text, text>
   PRIMARY KEY (id, mediate, second_ts )
 ) WITH CLUSTERING ORDER BY (mediate ASC, second_ts ASC)
 */
}

/*
CREATE TABLE demo.accounts4 (

    id int,
    name text,
    ascii1 ascii,
    bint1 bigint,
    blob1 blob,
    bool1 boolean,
    dec1 decimal,
    double1 double,
    flt1 float,
    inet1 inet,
    int1 int,
    text1 text,
    time1 timestamp,
    time2 timeuuid,
    mydate1 date,
    uuid1 uuid,
    varchar1 varchar,
    events set<int>,
    mylist list<float>,
    myset set<text>,
    adec list<decimal>,
    PRIMARY KEY (id, name, time1)
)WITH CLUSTERING ORDER BY (name ASC)
 */