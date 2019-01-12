package handler

import (
	"github.com/stevef1uk/cassuservice/parser"
	"os"
	"testing"
)


func TestFieldName(t *testing.T) {

	ret := Capitiseid(false, "id", false)
	if ret != "ID" {
		t.Errorf("Expencted ID got %s", ret)
	}

	ret = Capitiseid(false, "id", true)
	if ret != "id" {
		t.Errorf("Expencted id got %s", ret)
	}

	ret = Capitiseid(false, "Id", false)
	if ret != "ID" {
		t.Errorf("Expencted id got %s", ret)
	}

	ret = Capitiseid(false, "iD", false)
	if ret != "ID" {
		t.Errorf("Expencted id got %s", ret)
	}

	ret = Capitiseid(false, "ID", false)
	if ret != "ID" {
		t.Errorf("Expencted id got %s", ret)
	}

	ret = Capitiseid(false, "aid", false)
	if ret != "aID" {
		t.Errorf("Expencted id got %s", ret)
	}

	ret = Capitiseid(false, "aid1", false)
	if ret != "aid1" {
		t.Errorf("Expencted aid1 got %s", ret)
	}
}

func TestFieldName2(t *testing.T) {

	ret1 := CreateFile(true, "/tmp", "/tmp")
	ret1.Close()

	field := parser.FieldDetails{"test", "int", "", "", ""}
	ret2 := mapCassandraTypeToGoType(true, field, false, false, false, false)
	if ret2 != "int64" {
		t.Errorf("Expected int64 got %s", ret2)
	}

	field = parser.FieldDetails{"test", "int", "", "", ""}
	ret2 = mapCassandraTypeToGoType(true, field, false, true, false, false)
	if ret2 != "int" {
		t.Errorf("Expected int got %s", ret2)
	}

	ret3 := createTempVar("field1")
	if ret3 != "tmp_field1_0" {
		t.Errorf("Expected tmp_field1_0 got %s", ret3)
	}
	ret3 = createTempVar("field1")
	if ret3 != "tmp_field1_1" {
		t.Errorf("Expected tmp_field1_1 got %s", ret3)
	}

	field = parser.FieldDetails{"id", "TIMEUUID", "", "", ""}
	output := ""

	ret4 := setUpArrayTypes(true, output, field, false)
	_ = ret4
}

func TestFieldName3(t *testing.T) {

	ret5 := CapitaliseSplitFieldName(false, "id", false)
	if ret5 != "ID" {
		t.Errorf("Expected ID got %s", ret5)
	}

	ret5 = CapitaliseSplitFieldName(false, "id", true)
	if ret5 != "id" {
		t.Errorf("Expected id got %s", ret5)
	}

	ret5 = CapitaliseSplitFieldName(false, "steve", false)
	if ret5 != "Steve" {
		t.Errorf("Expected Steve got %s", ret5)
	}
	ret5 = CapitaliseSplitFieldName(false, "my_id_twoid", false)
	if ret5 != "MyIDTwoID" {
		t.Errorf("Expected MyIDTwoID got :%s:", ret5)
	}
	ret5 = CapitaliseSplitFieldName(false, "my_id_twoid", true)
	if ret5 != "my_id_twoid" {
		t.Errorf("Expected MyIDTwoID got :%s:", ret5)
	}
}

func TestSplice(t *testing.T) {

	path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test3/"
	_ = path
	ret6 :=  SpiceInHandler( false , path, "Employee", "" )
	_ = ret6


}

