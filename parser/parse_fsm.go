package parser

import (
	"regexp"
	"strings"
)

const (
	start      = "start"
	tableField = "tableField"
	primaryKey = "primaryKey"


	PRIMARY_STRING = "PRIMARY"
	FROZEN         = "FROZEN"
)

// State, Parser String, Next State
type fsmRow struct {
	expression string
	proc       func(bool, []string, [] string, fsmRow) bool
	nextState  string
	index      int
	reg        *regexp.Regexp
}

type fsm struct {
	rows        map[string][]fsmRow
	state       string
	breakString string
}

var theFSM fsm

//var theRegs []*regexp.Regexp

// Setup This function needs to be called first to initialise the FSM
// Note: I used https://regex-golang.appspot.com/assets/html/index.html to test the regular expressions
func Setup( debug bool ) {
	parseOutput = ParseOutput{}
	parseOutput.TableDetails = TableDetails{}
	//parseOutput.typeIndex = -1 // Set to -1 as always incrememeted when Type found
	parseOutput.TypeIndex = 0
	theFSM.rows = map[string][]fsmRow {
		start:{
			{`\s*CREATE TABLE\s*(\w+).(\w+)?`, processTable, tableField, 0, new(regexp.Regexp)},
			{`\s*CREATE TYPE\s*(\w+).(\w+)?`, processType, tableField, 0, new(regexp.Regexp)},
			  },
		tableField:{
			{`\s*PRIMARY\s+`, notePrimary, primaryKey, 0, new(regexp.Regexp)},
			// To handle text like: address_set list<frozen<city>>,
			{`\s*(\w+)\s+(\w+)\s*<\s*(\w+)\s*<\s*(\w+)\s*>>,?`, processSimpleFrozenField, tableField, 0, new(regexp.Regexp)},
			// TO handle a simple UTD field that can only occur in a table e.g. tSimple  frozen <simple>
			{`\s*(\w+)\s+(\w+)\s*<\s*(\w+)\s*>,?`, processSimpleFrozenField, tableField, 0, new(regexp.Regexp)},
			// To handle text like: address_set map<text, frozen <city>>,
			{`\s*(\w+)\s+(\w+)\s*<\s*(\w+),\s*\w+\s*\w+\s*<\s*(\w+)\s*>>,?`, processMapFrozenField, tableField, 0, new(regexp.Regexp)},
			// To handle normal table / type fields
			{`\s*(\w+)\s+(\w+)<?(\w+)?,?\s?(\w+)?`, processTableField, tableField, 0, new(regexp.Regexp)},
			// To terminate a Type definition
			{`\s*\)\s*;`, procNil, start, 0, new(regexp.Regexp)},
			  },
		primaryKey: {
			{`\s*PRIMARY\s+KEY\s*\(?\s*(\w+)\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*\)`, processPrimary, primaryKey, 0, new(regexp.Regexp)},
			{`\s*(\w+)\s+(\w+)\s+PRIMARY`, processPrimaryInLine, primaryKey, 0, new(regexp.Regexp)},
			},
	}

	theFSM.breakString = "WITH"
	theFSM.state = start

	index := 0
	for i, v := range theFSM.rows {
		for j, k := range v {
			if debug { println(i, k.expression)  }
			tableRe, err := regexp.Compile(k.expression)
			if err == nil {
				*theFSM.rows[i][j].reg = *tableRe
				_ = tableRe
			} else {
				println("Failed to compile expression %s", k.expression)
			}
		}
		index++
	}
}

func Reset() {
	theFSM = fsm{}
}


func parseLine(debug bool, text string, origText string) bool {

	ret := false
	var fields []string

	// Find RegEx to use based upon FSM state
	var rows [] fsmRow = theFSM.rows[theFSM.state]

	for _, j := range rows {
		result := j.reg.FindStringSubmatch(text)
		if result != nil {
			indexes := j.reg.FindStringSubmatchIndex(text)
			for i :=0; i  < len(indexes ); i++ {
				if ( i % 2 == 0 ) ||indexes[i] < 0 {
					continue
				}
				fields = append( fields, origText[indexes[i-1]:indexes[i]])
			}

			if j.proc != nil && j.proc(debug, result, fields, j) { parseLine( debug, text, origText ) } // Recurse if proc returns true
			break // Only allow one RegEx match within an FSM state
		}
	}

	return ret
}

// ParseText is called to process the Cassandra CQL definitions. setup and reset functions allow this function to do different things
func ParseText(debug bool, setUp func( bool), reset func(),  text string) ParseOutput {


	setUp( debug )

	lines := strings.SplitAfter(text, "\n")
	for _, v := range lines {
		if debug { println("Line:", v, "::") }
		if strings.Contains(v, theFSM.breakString) {
			if debug { println("I am out of here!") }
			break
		}
		parseLine(debug, strings.ToUpper(v), v)
	}
	reset()
	if debug { println("Finished ParseText") }
	return parseOutput
}
