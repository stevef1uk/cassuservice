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
	parser.Setup()

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

// */

}
