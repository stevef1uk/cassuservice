package parser

import (
	"regexp"
	"strings"
)

const (
	start      = "start"
	table      = "table"
	tableField = "tableField"
	primaryKey = "primaryKey"

	primaryString = "PRIMARY"
)

// State, Parser String, Next State
type fsmRow struct {
	expression string
	proc       func([]string) bool
	nextState  string
	index      int
	reg        *regexp.Regexp
}

type fsm struct {
	rows        map[string][]fsmRow
	state       string
	breakString string
}

// Null process function
func processTable(p []string) bool {
	ret := false
	for i, v := range p {
		println(i, v)
	}
	theFSM.state = table
	return ret
}

func processTableField(p []string) bool {
	ret := false
	for i, v := range p {
		if strings.TrimSpace(v) == primaryString {
			println("Found Primary Key )")
		}
		println(i, v)
	}
	theFSM.state = table
	return ret
}

var theFSM fsm

//var theRegs []*regexp.Regexp

// Setup THis function needs to be called first
func Setup() {

	//theFSM  := new(  fsm  )
	//theRegs := new( regexs )
	//var theRegs []*regexp.Regexp

	theFSM.rows = map[string][]fsmRow {
		start: []fsmRow{{`CREATE TABLE (\w+).(\w+)?`, processTable, tableField, 0, new(regexp.Regexp)}},
		table: []fsmRow{{`\s*(\w+)\s+(\w+)<?(\w+)?,?\s?(\w+)?`, processTableField, tableField, 0, new(regexp.Regexp)}},
	}
	//theFSM.rows[start] = fsmRow{`CREATE TABLE (\w+).(\w+)?`, processTable, tableField, 0, new(regexp.Regexp)}
	//theFSM.rows[table] = fsmRow{`\s*(\w+)\s+(\w+)<?(\w+)?,?\s?(\w+)?`, processTableField, tableField, 0, new(regexp.Regexp)}
	//theFSM.rows[table] = fsmRow{ `\s*(\w+)\s+(\w+)<?(\w+)?,?\s?(\w+)?`, processTableField, tableField, 0, new(regexp.Regexp ) }

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

func ParseLine(debug bool, text string) bool {

	ret := false

	// Find RegEx to use based upon FSM state
	var row fsmRow = theFSM.rows[theFSM.state][0]
	//println("FSM =", theFSM.rows[theFSM.state].reg )

	result := row.reg.FindStringSubmatch(text)
	if result != nil {
		row.proc(result)
	}
	return ret
}

func ParseText(debug bool, text string) {

	lines := strings.SplitAfter(text, "\n")
	for _, v := range lines {
		if strings.Contains(v, theFSM.breakString) {
			if debug {
				println("I am out of here!")
			}
			break
		}
		ParseLine(debug, v)
	}
	if debug {
		println("Finished ParseText")
	}
}
