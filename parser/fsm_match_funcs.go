package parser

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
	if debug { println("Processing Table Field") }
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
	theFSM.state = regRow.nextState // Force searching for other fields
	return ret
}


func processSimpleFrozenField(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	if debug { println("Processing Simple Frozen Field") }
	for i, v := range p {
		if debug { println(i, v) }
	}
	theFSM.state = tableField // Force searching for other fields
	return ret
}

func processMapFrozenField(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	if debug { println("Processing Map Frozen Field") }
	for i, v := range p {
		if debug { println(i, v) }
	}
	theFSM.state = tableField // Force searching for other fields
	return ret
}