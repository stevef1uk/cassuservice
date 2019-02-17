package handler

import (
	"testing"
)


func TestFieldName(t *testing.T) {


	ret := Capitiseid(false, "i", false)
	if ret != "I" {
		t.Errorf("Expected I got %s", ret)
	}


	ret = Capitiseid(false, "id", false)
	if ret != "ID" {
		t.Errorf("Expected ID got %s", ret)
	}

	ret = Capitiseid(false, "id", true)
	if ret != "Id" {
		t.Errorf("Expected id got %s", ret)
	}

	ret = Capitiseid(false, "Mid", true)
	if ret != "Mid" {
		t.Errorf("Expected Mid got %s", ret)
	}

	ret = Capitiseid(false, "Mid", false)
	if ret != "Mid" {
		t.Errorf("Expected Mid got %s", ret)
	}

	ret = Capitiseid(false, "Id", false)
	if ret != "ID" {
		t.Errorf("Expected id got %s", ret)
	}

	ret = Capitiseid(false, "iD", false)
	if ret != "ID" {
		t.Errorf("Expected id got %s", ret)
	}

	ret = Capitiseid(false, "ID", false)
	if ret != "ID" {
		t.Errorf("Expected id got %s", ret)
	}

	ret = Capitiseid(false, "aid", false)
	if ret != "Aid" {
		t.Errorf("Expected Aid got %s", ret)
	}

	ret = Capitiseid(false, "aid", true)
	if ret != "Aid" {
		t.Errorf("Expected Aid got %s", ret)
	}


	ret = Capitiseid(false, "aid_id", false)
	if ret != "Aid_id" {
		t.Errorf("Expected Aid_id got %s", ret)
	}
}

func TestFieldName2(t *testing.T) {

	ret1 := CreateFile(true, "/tmp", "/tmp", MAINFILE)
	ret1.Close()

	/*
	//field := parser.FieldDetails{"test", "int", "", "", ""}
	ret2 := mapCassandraTypeToGoType(true, "test","int", "test", false, false, false )
	if ret2 != "int64" {
		t.Errorf("Expected int64 got %s", ret2)
	}

	//field = parser.FieldDetails{"test", "int", "", "", ""}
	ret2 = mapCassandraTypeToGoType(true, "test", "int", "test", false, true, false )
	if ret2 != "int" {
		t.Errorf("Expected int got %s", ret2)
	}
*/
	ret3 := createTempVar("field1")
	if ret3 != "tmp_field1_0" {
		t.Errorf("Expected tmp_field1_0 got %s", ret3)
	}
	ret3 = createTempVar("field1")
	if ret3 != "tmp_field1_1" {
		t.Errorf("Expected tmp_field1_1 got %s", ret3)
	}

	//field = parser.FieldDetails{"id", "TIMEUUID", "", "", ""}

	//ret4 := setUpArrayTypes(true, "TIMEUUID", false)
	//_ = ret4
}

func TestFieldName3(t *testing.T) {

	ret5 := CapitaliseSplitFieldName(false, "id", false)
	if ret5 != "ID" {
		t.Errorf("Expected ID got %s", ret5)
	}

	ret5 = CapitaliseSplitFieldName(false, "id", true)
	if ret5 != "Id" {
		t.Errorf("Expected ID got %s", ret5)
	}

	ret5 = CapitaliseSplitFieldName(false, "steve", false)
	if ret5 != "Steve" {
		t.Errorf("Expected Steve got %s", ret5)
	}
	ret5 = CapitaliseSplitFieldName(false, "my_id_two_id", false)
	if ret5 != "MyIDTwoID" {
		t.Errorf("Expected MyIDTwoid got :%s:", ret5)
	}
	ret5 = CapitaliseSplitFieldName(false, "my_id_two_id", true)
	if ret5 != "MyIdTwoId" {
		t.Errorf("Expected MyIdTwoId got :%s:", ret5)
	}
}

/*
func TestSplice(t *testing.T) {

	path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
	_ = path
	ret6 :=  SpiceInHandler( false , path, "Employee", "" )
	_ = ret6


}

*/