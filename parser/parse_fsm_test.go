package parser

import (
	"testing"
)

func TestSimpleTable(t *testing.T) {

	expected := ParseText( false, Setup, Reset, `
			CREATE TABLE demo.accounts (
			id int,
			name text,
			PRIMARY KEY (id, name)
		) WITH CLUSTERING ORDER BY (name ASC)
		` )

	if expected.TableSpace != "DEMO" {
		t.Errorf("Tablespace incorrect, got: %s, want: %s.", expected.TableSpace, "DEMO")
	}
	if expected.TypeIndex != 0 {
		t.Errorf("TypeIndex incorrect, got: %d, want: %d.", expected.TypeIndex, 0)
	}
	if expected.TableDetails.TableName != "ACCOUNTS" {
		t.Errorf("TableName incorrect, got: %s, want: %s.", expected.TableDetails.TableName, "ACCOUNTS")
	}
	if expected.TableDetails.PkIndex != 2 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.PkIndex, 2)
	}
	if expected.TableDetails.DbPKFields[0] != "ID" {
		t.Errorf("DbPKFields[0] incorrect, got: %s, want: %s.", expected.TableDetails.DbPKFields[0], "ID")
	}
	if expected.TableDetails.DbPKFields[1] != "NAME" {
		t.Errorf("DbPKFields[1] incorrect, got: %s, want: %s.", expected.TableDetails.DbPKFields[0], "NAME")
	}
	if expected.TableDetails.TableFields.FieldIndex != 2 {
		t.Errorf("FieldIndex incorrect, got: %d, want: %d.", expected.TableDetails.TableFields.FieldIndex, 2)
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldName != "ID" {
		t.Errorf("DbFieldName incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldName, "ID" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldType != "INT" {
		t.Errorf("DbFieldType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldType, "INT")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldCollectionType != "" {
		t.Errorf("DbFieldCollectionType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldCollectionType, "")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldMapType != "" {
		t.Errorf("DbFieldMapType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldMapType, "")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldName != "NAME" {
		t.Errorf("DbFieldName incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldName, "NAME" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldType != "TEXT" {
		t.Errorf("DbFieldType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldType, "TEXT")
	}


}

func TestComplexTable(t *testing.T) {

	expected := ParseText( false, Setup, Reset, `
    CREATE TABLE Space1.accounts4 (
    id int,
    name text,
    ascii1 ascii,
    bint1 bigint,
    blob1 blob,
    bool1 boolean,
    counter1 counter,
    dec1 decimal,
    double1 double,
    flt1 float,
    inet1 inet,
    int1 int,
    text1 text,
    time1 timestamp,
    time2 timeuuid,
    uuid1 uuid,
    varchar1 varchar,
    events set<int>,
    mylist list<float>,
    myset set<text>,
    mymap  map<int, text>,
    PRIMARY KEY (id, name, time1)
) WITH CLUSTERING ORDER BY (name ASC)
		` )

	if expected.TableSpace != "SPACE1" {
		t.Errorf("Tablespace incorrect, got: %s, want: %s.", expected.TableSpace, "SPACE1")
	}
	if expected.TypeIndex != 0 {
		t.Errorf("TypeIndex incorrect, got: %d, want: %d.", expected.TypeIndex, 0)
	}
	if expected.TableDetails.TableName != "ACCOUNTS4" {
		t.Errorf("TableName incorrect, got: %s, want: %s.", expected.TableDetails.TableName, "ACCOUNTS4")
	}
	if expected.TableDetails.PkIndex != 3 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.PkIndex, 3)
	}
	if expected.TableDetails.DbPKFields[0] != "ID" {
		t.Errorf("DbPKFields[0] incorrect, got: %s, want: %s.", expected.TableDetails.DbPKFields[0], "ID")
	}
	if expected.TableDetails.DbPKFields[1] != "NAME" {
		t.Errorf("DbPKFields[1] incorrect, got: %s, want: %s.", expected.TableDetails.DbPKFields[0], "NAME")
	}
	if expected.TableDetails.DbPKFields[2] != "TIME1" {
		t.Errorf("DbPKFields[2] incorrect, got: %s, want: %s.", expected.TableDetails.DbPKFields[0], "TIME1")
	}
	if expected.TableDetails.TableFields.FieldIndex != 21 {
		t.Errorf("FieldIndex incorrect, got: %d, want: %d.", expected.TableDetails.TableFields.FieldIndex, 21)
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldName != "ID" {
		t.Errorf("DbFieldName incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldName, "ID" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldType != "INT" {
		t.Errorf("DbFieldType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldType, "INT")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldName != "NAME" {
		t.Errorf("DbFieldName incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldName, "NAME" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldType != "TEXT" {
		t.Errorf("DbFieldType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldType, "TEXT")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[14].DbFieldName != "TIME2" {
		t.Errorf("DbFieldName incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[14].DbFieldName, "TIME2" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[14].DbFieldType != "TIMEUUID" {
		t.Errorf("DbFieldType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[14].DbFieldType, "TIMEUUID")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[17].DbFieldName != "EVENTS" {
		t.Errorf("DbFieldName incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[17].DbFieldName, "EVENTS" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[17].DbFieldType != "SET" {
		t.Errorf("DbFieldType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[17].DbFieldType, "SET")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[17].DbFieldCollectionType != "INT" {
		t.Errorf("DbFieldCollectionType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[17].DbFieldCollectionType, "INT")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[17].DbFieldMapType != "" {
		t.Errorf("DbFieldMapType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[17].DbFieldMapType, "")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[20].DbFieldName != "MYMAP" {
		t.Errorf("DbFieldName incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[20].DbFieldName, "MYMAP" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[20].DbFieldType != "MAP" {
		t.Errorf("DbFieldType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[20].DbFieldType, "MAP")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[20].DbFieldCollectionType != "INT" {
		t.Errorf("DbFieldCollectionType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[20].DbFieldCollectionType, "INT")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[20].DbFieldMapType != "TEXT" {
		t.Errorf("DbFieldMapType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[20].DbFieldMapType, "TEXT")
	}
}

func TestTypeandTable(t *testing.T) {

	expected := ParseText(false, Setup, Reset, `
       CREATE TYPE TEST.city (
	       id int,
	       citycode text,
	       cityname text
	   );

       CREATE TABLE TEST.employee (
	       id int PRIMARY KEY,
	       address_map map<text, frozen <city>>,
	       address_list list<frozen<city>>,
	       address_set set<frozen<city>>,
	       name text
	   ) WITH CLUSTERING ORDER BY (name DESC);
		`)

	if expected.TableSpace != "TEST" {
		t.Errorf("Tablespace incorrect, got: %s, want: %s.", expected.TableSpace, "TEST")
	}
	if expected.TypeIndex != 1 {
		t.Errorf("TypeIndex incorrect, got: %d, want: %d.", expected.TypeIndex, 1)
	}
	if expected.TypeDetails[0].TypeFields.FieldIndex != 3 {
		t.Errorf("FieldIndex incorrect, got: %d, want: %d.", expected.TypeDetails[0].TypeFields.FieldIndex, 3)
	}
	if expected.TypeDetails[0].TypeName != "CITY" {
		t.Errorf("TypeName incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeName, "CITY")
	}
	if expected.TypeDetails[0].TypeFields.DbFieldDetails[0].DbFieldName != "ID" {
		t.Errorf("DbFieldName incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeFields.DbFieldDetails[0].DbFieldName, "ID")
	}
	if expected.TypeDetails[0].TypeFields.DbFieldDetails[0].DbFieldType != "INT" {
		t.Errorf("DbFieldType incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeFields.DbFieldDetails[0].DbFieldType, "INT")
	}
	if expected.TypeDetails[0].TypeFields.DbFieldDetails[1].DbFieldName != "CITYCODE" {
		t.Errorf("DbFieldName incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeFields.DbFieldDetails[1].DbFieldName, "CITYCODE")
	}
	if expected.TypeDetails[0].TypeFields.DbFieldDetails[1].DbFieldType != "TEXT" {
		t.Errorf("DbFieldType incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeFields.DbFieldDetails[1].DbFieldType, "INT")
	}
	if expected.TypeDetails[0].TypeFields.DbFieldDetails[2].DbFieldName != "CITYNAME" {
		t.Errorf("DbFieldName incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeFields.DbFieldDetails[2].DbFieldName, "CITYNAME")
	}
	if expected.TypeDetails[0].TypeFields.DbFieldDetails[2].DbFieldType != "TEXT" {
		t.Errorf("DbFieldType incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeFields.DbFieldDetails[2].DbFieldType, "INT")
	}
	if expected.TypeDetails[0].TypeFields.DbFieldDetails[2].DbFieldCollectionType != "" {
		t.Errorf("DbFieldCollectionType incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeFields.DbFieldDetails[2].DbFieldCollectionType, "")
	}
	if expected.TypeDetails[0].TypeFields.DbFieldDetails[2].DbFieldMapType != "" {
		t.Errorf("DbFieldMapType incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeFields.DbFieldDetails[2].DbFieldMapType, "")
	}

	if expected.TableDetails.TableName != "EMPLOYEE" {
		t.Errorf("TableName incorrect, got: %s, want: %s.", expected.TableDetails.TableName, "EMPLOYEE")
	}
	if expected.TableDetails.PkIndex != 1 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.PkIndex, 1)
	}

	if expected.TableDetails.TableFields.FieldIndex != 5 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.TableFields.FieldIndex, 5)
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldName != "ID" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldName, "ID" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldType != "INT" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldType, "INT")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldCollectionType != "" {
		t.Errorf("DbFieldCollectionType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldCollectionType, "")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldMapType != "" {
		t.Errorf("DbFieldMapType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldMapType, "")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldName != "ADDRESS_MAP" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldName, "ADDRESS_MAP" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldType != "MAP" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldType, "MAP")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldCollectionType != "TEXT" {
		t.Errorf("DbFieldCollectionType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldCollectionType, "TEXT")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldMapType != "CITY" {
		t.Errorf("DbFieldMapType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldMapType, "CITY")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[2].DbFieldName != "ADDRESS_LIST" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[2].DbFieldName, "ADDRESS_LIST" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[2].DbFieldType != "LIST" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[2].DbFieldType, "LIST")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[2].DbFieldCollectionType != "CITY" {
		t.Errorf("DbFieldCollectionType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[2].DbFieldCollectionType, "CITY")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[2].DbFieldMapType != "" {
		t.Errorf("DbFieldMapType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[2].DbFieldMapType, "")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[3].DbFieldName != "ADDRESS_SET" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[3].DbFieldName, "ADDRESS_SET" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[3].DbFieldType != "SET" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[3].DbFieldType, "SET")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[3].DbFieldCollectionType != "CITY" {
		t.Errorf("DbFieldCollectionType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[3].DbFieldCollectionType, "CITY")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[3].DbFieldMapType != "" {
		t.Errorf("DbFieldMapType incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[3].DbFieldMapType, "")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[4].DbFieldName != "NAME" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[4].DbFieldName, "NAME" )
	}
	if expected.TableDetails.TableFields.DbFieldDetails[4].DbFieldType != "TEXT" {
		t.Errorf("TypeIndex incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[4].DbFieldType, "TEXT")
	}
}

func TestMultipleTypes(t *testing.T) {

	expected := ParseText(false, Setup, Reset, `
CREATE TYPE demo.debtor_agent (
  schemeName text,
  identification text
);


CREATE TYPE demo.debtor_account (
  schemeName text,
  identification text,
  name text,
  secondaryIdentification text
);


CREATE TYPE demo.creditor_agent (
  schemeName text,
  identification text
);

CREATE TABLE demo.pisp_submissions_per_id (
	submissionId uuid,
	timeBucket text,
	debtorAgent debtor_agent,
	debtorAccount debtor_account,
	creditorAgent creditor_agent,
	lastUpdatedAt timestamp,
	PRIMARY KEY (submissionId, lastUpdatedAt)
) WITH CLUSTERING ORDER BY (lastUpdatedAt DESC)
	`)

	if expected.TableSpace != "DEMO" {
		t.Errorf("Tablespace incorrect, got: %s, want: %s.", expected.TableSpace, "DEMO")
	}
	if expected.TypeIndex != 3 {
		t.Errorf("TypeIndex incorrect, got: %d, want: %d.", expected.TypeIndex, 3)
	}
	if expected.TypeDetails[0].TypeName != "DEBTOR_AGENT" {
		t.Errorf("TypeName incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeName, "DEBTOR_AGENT")
	}
	if expected.TypeDetails[0].TypeFields.FieldIndex != 2 {
		t.Errorf("FieldIndex incorrect, got: %d, want: %d.", expected.TypeDetails[0].TypeFields.FieldIndex, 2)
	}
	if expected.TypeDetails[1].TypeName != "DEBTOR_ACCOUNT" {
		t.Errorf("TypeName incorrect, got: %s, want: %s.", expected.TypeDetails[1].TypeName, "DEBTOR_ACCOUNT")
	}
	if expected.TypeDetails[1].TypeFields.FieldIndex != 4 {
		t.Errorf("FieldIndex incorrect, got: %d, want: %d.", expected.TypeDetails[1].TypeFields.FieldIndex, 4)
	}
	if expected.TypeDetails[2].TypeName != "CREDITOR_AGENT" {
		t.Errorf("TypeName incorrect, got: %s, want: %s.", expected.TypeDetails[2].TypeName, "CREDITOR_AGENT")
	}
	if expected.TypeDetails[2].TypeFields.FieldIndex != 2 {
		t.Errorf("FieldIndex incorrect, got: %d, want: %d.", expected.TypeDetails[2].TypeFields.FieldIndex, 2)
	}
	if expected.TableDetails.TableName != "PISP_SUBMISSIONS_PER_ID" {
		t.Errorf("TableName incorrect, got: %s, want: %s.", expected.TableDetails.TableName, "PISP_SUBMISSIONS_PER_ID")
	}
	if expected.TableDetails.PkIndex != 2 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.PkIndex, 2)
	}

	if expected.TableDetails.TableFields.FieldIndex != 6 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.TableFields.FieldIndex, 6)
	}
	if expected.TableDetails.PkIndex != 2 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.PkIndex, 2)
	}
}

func TestSimpleFrozen(t *testing.T) {

	expected := ParseText(false, Setup, Reset, `
CREATE TYPE demo.simple1(
    id int,
    citycode text
);


CREATE TYPE demo.simple (
       id int,
       dummy text,
       mediate TIMESTAMP,
       eStruct  set <frozen <simple1>>,
    );

CREATE TABLE demo.employee1 (
    id int PRIMARY KEY,
    tSimple  frozen <simple>
) WITH CLUSTERING ORDER BY (name ASC) ;
	`)

	if expected.TableSpace != "DEMO" {
		t.Errorf("Tablespace incorrect, got: %s, want: %s.", expected.TableSpace, "DEMO")
	}
	if expected.TypeIndex != 2 {
		t.Errorf("TypeIndex incorrect, got: %d, want: %d.", expected.TypeIndex, 2)
	}
	if expected.TypeDetails[0].TypeName != "SIMPLE1" {
		t.Errorf("TypeName incorrect, got: %s, want: %s.", expected.TypeDetails[0].TypeName, "SIMPLE")
	}
	if expected.TypeDetails[0].TypeFields.FieldIndex != 2 {
		t.Errorf("FieldIndex incorrect, got: %d, want: %d.", expected.TypeDetails[0].TypeFields.FieldIndex, 2)
	}
	if expected.TypeDetails[1].TypeName != "SIMPLE" {
		t.Errorf("TypeName incorrect, got: %s, want: %s.", expected.TypeDetails[1].TypeName, "SIMPLE")
	}
	if expected.TypeDetails[1].TypeFields.FieldIndex != 4 {
		t.Errorf("FieldIndex incorrect, got: %d, want: %d.", expected.TypeDetails[1].TypeFields.FieldIndex, 4)
	}
	if expected.TableDetails.TableName != "EMPLOYEE1" {
		t.Errorf("TableName incorrect, got: %s, want: %s.", expected.TableDetails.TableName, "EMPLOYEE1")
	}
	if expected.TableDetails.PkIndex != 1 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.PkIndex, 1)
	}
	if expected.TableDetails.TableFields.FieldIndex != 2 {
		t.Errorf("PkIndex incorrect, got: %d, want: %d.", expected.TableDetails.TableFields.FieldIndex, 2)
	}
	if expected.TableDetails.DbPKFields[0] != "ID" {
		t.Errorf("DbPKFields[0] incorrect, got: %s, want: %s.", expected.TableDetails.DbPKFields[0], "ID")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldName != "ID" {
		t.Errorf(".DbFieldDetails[0].DbFieldName incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[0].DbFieldName, "ID")
	}
	if expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldName != "TSIMPLE" {
		t.Errorf("DbFieldDetails[1].DbFieldName  incorrect, got: %s, want: %s.", expected.TableDetails.TableFields.DbFieldDetails[1].DbFieldName, "TSIMPLE")
	}
}