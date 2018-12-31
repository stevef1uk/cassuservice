package main

import (
	"flag"
	"github.com/stevef1uk/cassuservice/parser"
)

// Set debug flag to true to output logging information to stdout
var debug = false

func main() {

	debugPtr := flag.Bool("debug", false, "Set flag to true to write trace information to stdout")

	flag.Parse()
	debug = *debugPtr
	//parser.Setup()

/*
		parser.ParseText( debug, `
			CREATE TABLE demo.accounts (
			id int,
			name text,
			PRIMARY KEY (id, name)
		) WITH CLUSTERING ORDER BY (name ASC)
		` )

*/

	 parser.ParseText( debug, `
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


/*
	parser.ParseText(debug, `
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

*/

}