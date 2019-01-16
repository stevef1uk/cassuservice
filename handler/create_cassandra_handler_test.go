package handler

import (
	"github.com/stevef1uk/cassuservice/parser"
	"io"
	"io/ioutil"
	"os"
	"log"
	"testing"
)

//address_list set<frozen<simple>>,
const (

	CSQ_TEST1 = `

    CREATE TYPE demo.simple (
       dummy text
    );

    CREATE TYPE demo.city (
    id int,
    citycode text,
    cityname text,
    test_int int,
    lastUpdatedAt TIMESTAMP,
    myfloat float
    events set<int>,
    mymap  map<int, text>
    address_list set<frozen<simple>>,
);

CREATE TABLE demo.employee (
    id int,
    address_set set<frozen<city>>,
    my_List list<frozen<simple>>,
    name text,
    mediate TIMESTAMP,
    second_ts timestamp,
    tevents set<int>,
    tmylist list<float>
    tmymap  map<int, text>
   PRIMARY KEY (id, mediate, second_ts )
 ) WITH CLUSTERING ORDER BY (mediate ASC, second_ts ASC)
`
	EXPECTED_OUTPUT_TEST1 = `// GENERATED FILE so do not edit or will be overwritten upon next generate
package data

import (
    "github.com/stevef1uk/test3/models"
    "github.com/stevef1uk/test3/restapi/operations"
    "middleware "github.com/go-openapi/runtime/middleware"
    "github.com/gocql/gocql"
    "os"
    "log"
    "time"
     strfmt "github.com/go-openapi/strfmt"
    "gopkg.in/inf.v0"
    "strconv"
)
`

)


func performCreateTest1( debug bool, test string, cql string, expected string , t *testing.T ) {

	// Mock stdin
	file := tempFile()
	defer os.Remove(file.Name())
	//Mock Stdout
	fileout := tempFile()
	defer os.Remove(fileout.Name())


	err := ioutil.WriteFile(file.Name(), []byte(cql), 0666)
	if err != nil {
		log.Fatal(err)
	}

	file.Sync()
	input(file)

	parse1 := parser.ParseText( false, parser.Setup, parser.Reset, cql )
	CreateCode( debug, "/tmp", "github.com/stevef1uk/test3", parse1,  "",  "",  0, false , false , true   )


	// Read generated file
	path := "/tmp/data/" + MAINFILE
	fileo, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	byteSlice := make([]byte, 10000)
	numBytesRead, err := io.ReadFull(fileo, byteSlice)
	if err != nil {
		log.Printf("Number of bytes read: %d\n", numBytesRead)
	}
	tmpbytes := byteSlice[0:numBytesRead]
	s := string(tmpbytes[:])

	if ( s != expected) {
		t.Errorf("Create Handler Test %s,  incorrect output read \n:%s:, want\n:%s:", test, s, expected)
	}
}


/*

path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test3/"
ret6 :=  SpiceInHandler( false , path, "Employee", "" )
_ = ret6
*/

func Test1(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST1, EXPECTED_OUTPUT_TEST1, t )
}


