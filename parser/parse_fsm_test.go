package parser

import (
	"testing"
)

func TestSimple(t *testing.T) {

	expected := ParseText( false, `
			CREATE TABLE demo.accounts (
			id int,
			name text,
			PRIMARY KEY (id, name)
		) WITH CLUSTERING ORDER BY (name ASC)
		` )

	if expected.TableSpace != "DEMO" {
		t.Errorf("Tablespace incorrect, got: %s, want: %s.", expected.TableSpace, "DEMO")
	}
	if expected.typeIndex != 0 {
		t.Errorf("TypeIndex incorrect, got: %d, want: %d.", expected.typeIndex, 0)
	}
	if expected.TableDetails.TableName != "ACCOUNTS" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableName, "ACCOUNTS")
	}
	if expected.TableDetails.PkIndex != 2 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.PkIndex, 2)
	}
	if expected.TableDetails.DbPKFields[0] != "ID" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.DbPKFields[0], "ID")
	}
	if expected.TableDetails.DbPKFields[1] != "NAME" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.DbPKFields[0], "NAME")
	}
	if expected.TableDetails.FieldIndex != 2 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.FieldIndex, 2)
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldName != "ID" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldName, "ID" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldType != "INT" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldType, "INT")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldName != "NAME" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldName, "NAME" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldType != "TEXT" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldType, "TEXT")
	}
}