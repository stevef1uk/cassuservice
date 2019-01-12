package handler

import (
	"github.com/stevef1uk/cassuservice/parser"
	"os"
	"testing"
)


func TestReviceFieldName(t *testing.T) {


	ret := Capitiseid( false, "id", false )
	if ret != "ID" {
		t.Errorf("Expencted ID got %s", ret )
	}

	ret = Capitiseid( false, "id", true )
	if ret != "id" {
		t.Errorf("Expencted id got %s", ret )
	}

	ret = Capitiseid( false, "Id", false )
	if ret != "ID" {
		t.Errorf("Expencted id got %s", ret )
	}

	ret = Capitiseid( false, "iD", false )
	if ret != "ID" {
		t.Errorf("Expencted id got %s", ret )
	}

	ret = Capitiseid( false, "ID", false )
	if ret != "ID" {
		t.Errorf("Expencted id got %s", ret )
	}

	ret = Capitiseid( false, "aid", false )
	if ret != "aID" {
		t.Errorf("Expencted id got %s", ret )
	}
	ret = Capitiseid( true, "aid1", false )
	if ret != "aid1" {
		t.Errorf("Expencted id got %s", ret )
	}


	ret1:= CreateFile( true , "/tmp", "/tmp" )
	ret1.Close()

	field := parser.FieldDetails{ "test", "int", "", "", ""}
	ret2 := mapCassandraTypeToGoType( true, field,false,   false, false, false)
	if ret2 != "int64" {
		t.Errorf("Expected int64 got %s", ret2 )
	}

	field = parser.FieldDetails{ "test", "int", "", "", ""}
	ret2 = mapCassandraTypeToGoType( true, field,false,   true, false, false)
	if ret2 != "int" {
		t.Errorf("Expected int got %s", ret2 )
	}

	ret3 := createTempVar ( "field1" )
	if ret3 != "tmp_field1_0" {
		t.Errorf("Expected tmp_field1_0 got %s", ret3 )
	}
	ret3 = createTempVar ( "field1" )
	if ret3 != "tmp_field1_1" {
		t.Errorf("Expected tmp_field1_1 got %s", ret3 )
	}

	field = parser.FieldDetails{ "id", "TIMEUUID", "", "", ""}
	output := ""

	ret4 := setUpArrayTypes(  true , output , field,  false )
	_ = ret4


	ret5 := CapitaliseSplitFieldName( true, "id", false )
	if ret5 != "ID" {
		t.Errorf("Expected ID got %s", ret5 )
	}
	ret5 = CapitaliseSplitFieldName( true, "id", true )
	if ret5 != "id" {
		t.Errorf("Expected id got %s", ret5 )
	}
	ret5 = CapitaliseSplitFieldName( true, "steve", false )
	if ret5 != "Steve" {
		t.Errorf("Expected Steve got %s", ret5 )
	}
	ret5 = CapitaliseSplitFieldName( true, "my_id_twoid", false )
	if ret5 != "MyIDTwoID" {
		t.Errorf("Expected MyIDTwoID got :%s:", ret5 )
	}
	ret5 = CapitaliseSplitFieldName( true, "my_id_twoid", true )
	if ret5 != "my_id_twoid" {
		t.Errorf("Expected MyIDTwoID got :%s:", ret5 )
	}

	path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test3/"
	ret6 :=  SpiceInHandler( false , path, "Employee", "" )
	_ = ret6

	parse1 := parser.ParseText( false, parser.Setup, parser.Reset, `
			CREATE TABLE demo.accounts (
			id int,
			name text,
			PRIMARY KEY (id, name)
		) WITH CLUSTERING ORDER BY (name ASC)
		` )

	CreateCode( true, path, parse1,  "",  "",  0, false , false , true   )

	}

