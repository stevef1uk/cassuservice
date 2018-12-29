package parser

import (
	"regexp"
	"strings"
)

const (
	start      = "start"
	tableField = "tableField"
	primaryKey = "primaryKey"


	primaryString = "PRIMARY"
)

// State, Parser String, Next State
type fsmRow struct {
	expression string
	proc       func(bool, []string, fsmRow) bool
	nextState  string
	index      int
	reg        *regexp.Regexp
}

type fsm struct {
	rows        map[string][]fsmRow
	state       string
	breakString string
}

// The following set of functions are called by the FSM processing logic when a regex match is made
func processTable(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	if debug { println( "Parsing new Table" ) }
	for i, v := range p {
		println(i, v)
	}
	theFSM.state = regRow.nextState
	return ret
}

func processType(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	if debug { println( "Parsing new Type" ) }
	for i, v := range p {
		if debug { println(i, v) }
	}
	theFSM.state = regRow.nextState
	return ret
}

func processTableField(debug bool, p []string, regRow fsmRow ) bool {
	ret := false
	for i, v := range p {
		if debug { println(i, v) }
	}
	theFSM.state = regRow.nextState
	return ret
}

// As primary key string identified return true so that the real function to process a primary key will be called
func notePrimary(debug bool, p []string, regRow fsmRow) bool {
	ret := true
	for i, v := range p {
		if debug { println(i, v) }
	}
	if debug { println("Found Primary Key") }
	theFSM.state = regRow.nextState
	return ret
}

//
func processPrimary(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	for i, v := range p {
		if debug { println(i, v) }
	}
	theFSM.state = regRow.nextState
	return ret
}

func processPrimaryInLine(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	for i, v := range p {
		if debug { println(i, v) }
	}
	theFSM.state = tableField // Force searching for other fields
	return ret
}

func procNil(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	for i, v := range p {
		if debug { println(i, v) }
	}
	theFSM.state = regRow.nextState // Force searching for other fields
	return ret
}

// End of processing logic & start of main FSM logic

var theFSM fsm

//var theRegs []*regexp.Regexp

// Setup This function needs to be called first to initialise the FSM
func Setup() {
	theFSM.rows = map[string][]fsmRow {
		start: []fsmRow{
						{`\s*CREATE TABLE\s*(\w+).(\w+)?`, processTable, tableField, 0, new(regexp.Regexp)},
						{`\s*CREATE TYPE\s*(\w+).(\w+)?`, processType, tableField, 0, new(regexp.Regexp)},
			},
		tableField: []fsmRow{
		                {`\s*PRIMARY\s+`, notePrimary, primaryKey, 0, new(regexp.Regexp)},
			            {`\s*(\w+)\s+(\w+)<?(\w+)?,?\s?(\w+)?`, processTableField, tableField, 0, new(regexp.Regexp)},
			            {`\s*\)\s*;`, procNil, start, 0, new(regexp.Regexp)},
			},
		primaryKey: {{`\s*PRIMARY\s+KEY\s*\(?\s*(\w+)\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*,?\s*(\w+)*\s*\)x`, processPrimary, primaryKey, 0, new(regexp.Regexp)},
			         {`\s*(\w+)\s+(\w+)\s+PRIMARY`, processPrimaryInLine, primaryKey, 0, new(regexp.Regexp)},
			},
	}

	theFSM.breakString = "WITH"
	theFSM.state = start

	index := 0
	for i, v := range theFSM.rows {
		for j, k := range v {
			println(i, k.expression)
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

func parseLine(debug bool, text string) bool {

	ret := false

	// Find RegEx to use based upon FSM state
	var rows [] fsmRow = theFSM.rows[theFSM.state]

	for _, j := range rows {
		result := j.reg.FindStringSubmatch(text)
		if result != nil {
			if j.proc != nil && j.proc(debug, result, j) { parseLine( debug, text ) }
			break
		}
	}

	return ret
}

// ParseText needs to be called after FMS has been initialised
func ParseText(debug bool, text string) {

	lines := strings.SplitAfter(text, "\n")
	for _, v := range lines {
		println("Line:", v, "::")
		if strings.Contains(v, theFSM.breakString) {
			if debug {
				println("I am out of here!")
			}
			break
		}
		parseLine(debug, v)
	}
	if debug {
		println("Finished ParseText")
	}
}
