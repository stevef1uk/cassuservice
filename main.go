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
}
