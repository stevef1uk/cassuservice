package parser

import "strings"

// The following set of functions are called by the FSM processing logic when a regex match is made
func processTable(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	if debug { println( "Parsing new Table" ) }
	parseOutput.TableSpace = p[1]
	parseOutput.TableDetails.TableName = p[2]
	parseOutput.inTable = true
	theFSM.state = regRow.nextState
	return ret
}

func processType(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	if debug { println( "Parsing new Type" ) }
	for i, v := range p {
		if debug { println(i, v) }
	}
	parseOutput.inTable = false
	//parseOutput.typeIndex = parseOutput.typeIndex + 1
	parseOutput.TypeDetails[parseOutput.typeIndex].TypeName = p[2]
	theFSM.state = regRow.nextState
	return ret
}

func processTableField(debug bool, p []string, regRow fsmRow ) bool {
	ret := false
	var fieldDetails *FieldDetails
	var index int
	if debug {
		println("Processing Table Field")
		for i, v := range p {
			println(i, v)
		}
	}

	if ( parseOutput.inTable ) {
		index = parseOutput.TableDetails.FieldIndex
		parseOutput.TableDetails.FieldIndex = parseOutput.TableDetails.FieldIndex + 1
		fieldDetails = &parseOutput.TableDetails.TableFields.DbFieldDetails[index]
	} else {
		index = parseOutput.TypeDetails[parseOutput.typeIndex].FieldIndex
		parseOutput.TypeDetails[parseOutput.typeIndex].FieldIndex = parseOutput.TypeDetails[parseOutput.typeIndex].FieldIndex + 1
		fieldDetails = &parseOutput.TypeDetails[parseOutput.typeIndex].TypeFields.DbFieldDetails[index]
	}

	fieldDetails.DbFieldName = p[1]
	fieldDetails.DbFieldType = p[2]
	if p[3] != "" {
		fieldDetails.DbFieldCollectionType = p[3]
	}
	if p[4] != "" {
		fieldDetails.DbFieldMapType = p[4]
	}

	theFSM.state = regRow.nextState
	return ret
}

// As primary key string identified return true so that the real function to process a primary key will be called
func notePrimary(debug bool, p []string, regRow fsmRow) bool {
	ret := true
	if debug { println("Found Primary Key") }
	theFSM.state = regRow.nextState
	return ret
}

// This function handles the PRIMARY KEY (id, name) form
func processPrimary(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	if debug { println("Parsing normal PRIMARY KEY line") }

	for i, v := range p[1:] { // Element 0 contains the matched line
		if debug { println(i, v) }
		if v == "" { break }
		parseOutput.TableDetails.DbPKFields[parseOutput.TableDetails.PkIndex] = v
		parseOutput.TableDetails.PkIndex = parseOutput.TableDetails.PkIndex + 1
	}

	theFSM.state = regRow.nextState
	return ret
}

// This function processes the id int PRIMARY KEY form
func processPrimaryInLine(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	var j int = 0
	fields := make( []string, len(p) + 4 )
	if debug { println("Parsing field PRIMARY KEY annotation") }
	if debug { println(p[0], " - Storing", p[1]) }
	parseOutput.TableDetails.DbPKFields[parseOutput.TableDetails.PkIndex] = p[1]
	parseOutput.TableDetails.PkIndex = parseOutput.TableDetails.PkIndex + 1
	for i, v := range p {
		if debug { println(i, v) }
		if strings.ToUpper(strings.TrimSpace(v)) == PRIMARY_STRING {
			break
		}
		fields[j] = v
		j = j + 1
	}
	processTableField( debug, fields, regRow)
	// @TODO ensure PK stored as a field!
	theFSM.state = tableField // Force searching for other fields
	return ret
}

func procNil(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	parseOutput.typeIndex = parseOutput.typeIndex + 1
	theFSM.state = regRow.nextState // Force searching for other fields
	return ret
}


func processSimpleFrozenField(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	var j int = 0
	if debug { println("Processing Simple Frozen Field") }
	fields := make( []string, len(p) + 4 )
	for i, v := range p {
		if debug {
			println(i, v)
		}
		if strings.ToUpper(strings.TrimSpace(v)) == FROZEN {
			continue
		}
		fields[j] = v
		j = j + 1
	}
	processTableField( debug, fields, regRow)
	theFSM.state = tableField // Force searching for other fields
	return ret
}

func processMapFrozenField(debug bool, p []string, regRow fsmRow) bool {
	ret := false
	if debug { println("Processing Map Frozen Field") }

	processSimpleFrozenField( debug, p, regRow)
	theFSM.state = tableField // Force searching for other fields
	return ret
}