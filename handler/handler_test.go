package handler

import (
	"github.com/stevef1uk/cassuservice/parser"
	"testing"
)


func TestReviceFieldName(t *testing.T) {


	ret := ReviseFieldName( false, "id", false )
	if ret != "ID" {
		t.Errorf("Expencted ID got %s", ret )
	}

	ret = ReviseFieldName( false, "id", true )
	if ret != "id" {
		t.Errorf("Expencted id got %s", ret )
	}

	ret = ReviseFieldName( false, "Id", false )
	if ret != "ID" {
		t.Errorf("Expencted id got %s", ret )
	}

	ret = ReviseFieldName( false, "iD", false )
	if ret != "ID" {
		t.Errorf("Expencted id got %s", ret )
	}

	ret = ReviseFieldName( false, "ID", false )
	if ret != "ID" {
		t.Errorf("Expencted id got %s", ret )
	}

	ret = ReviseFieldName( false, "aid", false )
	if ret != "aID" {
		t.Errorf("Expencted id got %s", ret )
	}
	ret = ReviseFieldName( true, "aid1", false )
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
}