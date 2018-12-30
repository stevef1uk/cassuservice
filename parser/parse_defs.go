package parser

const (
	// Assume no more than 25 fields in a table
	MAXFIELDS = 25
)




type FieldDetails struct {
	DbFieldName string
	DbFieldType string
	FieldFormat string
	DbFieldCollectionType string
	DbFieldMapType string
}


type AllFieldDetails struct {
	DbFieldDetails [MAXFIELDS] FieldDetails
}


type TableDetails struct {
	TableName    string
	TableFields  AllFieldDetails
	DbPKFields   [MAXFIELDS] string
	FieldIndex	int
}

type TypeDetails struct {
	TypeName    string
	TypeFields  AllFieldDetails
	FieldIndex	int
}

type ParseOutput struct {
	TableSpace   string
	TableDetails TableDetails
	TypeDetails [MAXFIELDS] TypeDetails
	inTable bool
	typeIndex int
}

var parseOutput ParseOutput

