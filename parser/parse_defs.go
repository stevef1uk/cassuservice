package parser

const (
	// Assume no more than 25 fields in a table
	MAXFIELDS = 25
)




type FieldDetails struct {
	DbFieldName string
	DbFieldType string
	OrigFieldName string
	DbFieldCollectionType string
	DbFieldMapType string
}


type AllFieldDetails struct {
	DbFieldDetails [MAXFIELDS] FieldDetails
	FieldIndex int
}


type TableDetails struct {
	TableName    string
	TableFields  AllFieldDetails
	DbPKFields   [MAXFIELDS] string
	PkIndex	int
}

type TypeDetails struct {
	TypeName    string
	TypeFields  AllFieldDetails
}

type ParseOutput struct {
	TableSpace   string
	TableDetails TableDetails
	TypeDetails [MAXFIELDS] TypeDetails
	inTable bool
	TypeIndex int
}

var parseOutput ParseOutput

