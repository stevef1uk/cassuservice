package handler

import (
	"github.com/stevef1uk/cassuservice/parser"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

//address_list set<frozen<simple>>,
//lastUpdatedAt TIMESTAMP,

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
    myfloat float,
    events set<int>,
    mymap  map<text, text>,
    address_list set<frozen<simple>>
);

CREATE TABLE demo.employee (
    id int,
    address_set set<frozen<city>>,
    my_List list<frozen<simple>>,
    name text,
    mediate TIMESTAMP,
    second_ts TIMESTAMP,
    tevents set<int>,
    tmylist list<float>,
    tmymap  map<text, text>,
   PRIMARY KEY (id, mediate, second_ts )
 ) WITH CLUSTERING ORDER BY (mediate ASC, second_ts ASC)
`
	EXPECTED_OUTPUT_TEST1 = `// GENERATED FILE so do not edit or will be overwritten upon next generate
package data

import (
    "github.com/stevef1uk/test4/models"
    "github.com/stevef1uk/test4/restapi/operations"
    middleware "github.com/go-openapi/runtime/middleware"
    "github.com/gocql/gocql"
    "os"
    "log"
    "time"
     strfmt "github.com/go-openapi/strfmt"
)

type Simple struct {
    Dummy string `+"`"+`cql:"dummy"`+"`"+`
}

type City struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Citycode string `+"`"+`cql:"citycode"`+"`"+`
    Cityname string `+"`"+`cql:"cityname"`+"`"+`
    TestInt int `+"`"+`cql:"test_int"`+"`"+`
    Lastupdatedat time.Time `+"`"+`cql:"lastupdatedat"`+"`"+`
    Myfloat float32 `+"`"+`cql:"myfloat"`+"`"+`
    Events []int `+"`"+`cql:"events"`+"`"+`
    Mymap models.CityMymap `+"`"+`cql:"mymap"`+"`"+`
    AddressList models.CityAddressList `+"`"+`cql:"address_list"`+"`"+`
}


var cassuservice_session *gocql.Session

func SetUp() {
  var err error
  log.Println("Tring to connect to Cassandra database using ", os.Getenv("CASSANDRA_SERVICE_HOST"))
  cluster := gocql.NewCluster(os.Getenv("CASSANDRA_SERVICE_HOST"))
  cluster.Keyspace = "demo"
  cluster.Consistency = gocql.One
  cassuservice_session, err = cluster.CreateSession()
  if ( err != nil ) {
    log.Fatal("Have you remembered to set the env var $CASSANDRA_SERVICE_HOST as connection to Cannandra failed with error = ", err)
  } else {
    log.Println("Yay! Connection to Cannandra established")
  }
}

func Stop() {
    log.Println("Shutting down the service handler")
  cassuservice_session.Close()
}

func Search(params operations.GetEmployeeParams) middleware.Responder {

    var ID int64
    _ = ID
    var AddressSet models.AddressSet
    _ = AddressSet
    var MyList models.MyList
    _ = MyList
    var Name string
    _ = Name
    var Mediate time.Time
    _ = Mediate
    var SecondTs time.Time
    _ = SecondTs
    var Tevents []int64
    _ = Tevents
    var Tmylist []float64
    _ = Tmylist
    var Tmymap []string
    _ = Tmymap
    tmp_cassuservice_tmp_time_0 := strfmt.NewDateTime().String()

    codeGenRawTableResult := map[string]interface{}{}

    Mediate,_ = time.Parse(time.RFC3339,params.Mediate.String() ) 
    SecondTs,_ = time.Parse(time.RFC3339,params.SecondTs.String() ) 
    if err := cassuservice_session.Query(`+"`"+` SELECT id, address_set, my_list, name, mediate, second_ts, tevents, tmylist, tmymap FROM employee WHERE id = ? and mediate = ? and second_ts = ? `+"`"+`,params.ID,Mediate,SecondTs).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetEmployeeBadRequest()
    }
    payLoad := operations.NewGetEmployeeOK()
    payLoad.Payload = make([]*models.GetEmployeeOKBodyItems,1)
    payLoad.Payload[0] = new(models.GetEmployeeOKBodyItems)
    retParams := payLoad.Payload[0]
    tmp_ID_1 := codeGenRawTableResult["id"].(int)
    ID = int64(tmp_ID_1)
    retParams.ID = &ID
    tmp_City_2, ok := codeGenRawTableResult["address_set"].([]map[string]interface{})
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for address_set", ok )
    }
    retParams.AddressSet = make([]*models.City, len(tmp_City_2))
    for i3, v3 := range tmp_City_2 {
    
      tmp_mymap_4 := v3["mymap"].(map[string]string)
      tmp_address_list_6:= v3["address_list"].([]map[string]interface{})
      tmp_address_list_7:= make(models.CityAddressList, len(tmp_address_list_6) )
      
          for i4, v4 := range tmp_address_list_6 {
    
          tmp_Simple_8 := &Simple{
    
                v4["dummy"].(string),
              }
                
            tmp_address_list_7[i4] = &models.Simple{}
            tmp_address_list_7[i4].Dummy = tmp_Simple_8.Dummy
            }
      tmp_City_3 := &City{
    
          v3["id"].(int),
          v3["citycode"].(string),
          v3["cityname"].(string),
          v3["test_int"].(int),
          v3["lastupdatedat"].(time.Time),
          v3["myfloat"].(float32),
          v3["events"].([]int),
          tmp_mymap_4,
          tmp_address_list_7,
        }
          
      retParams.AddressSet[i3] = &models.City{}
      tmp_ID_9 := int64(tmp_City_3.ID)
      retParams.AddressSet[i3].ID = tmp_ID_9      
      retParams.AddressSet[i3].Citycode = tmp_City_3.Citycode      
      retParams.AddressSet[i3].Cityname = tmp_City_3.Cityname      
      tmp_TestInt_10 := int64(tmp_City_3.TestInt)
      retParams.AddressSet[i3].TestInt = tmp_TestInt_10      
      tmp_cassuservice_tmp_time_0 = tmp_City_3.Lastupdatedat.String()
      tmp_Lastupdatedat_11 := tmp_cassuservice_tmp_time_0[0:10] + "T" + tmp_cassuservice_tmp_time_0[11:19] + "." + tmp_cassuservice_tmp_time_0[20:22]
      if tmp_cassuservice_tmp_time_0[22] == ' ' {
        tmp_cassuservice_tmp_time_0 = tmp_Lastupdatedat_11 + "0" + "Z" 
      } else { 
        tmp_cassuservice_tmp_time_0 = tmp_Lastupdatedat_11 + "Z"
      }
      tmp_Lastupdatedat_12, _  := strfmt.ParseDateTime(tmp_cassuservice_tmp_time_0)
      tmp_Lastupdatedat_13 := tmp_Lastupdatedat_12.String()
      retParams.AddressSet[i3].Lastupdatedat = tmp_Lastupdatedat_13      
      tmp_Myfloat_14 := float64(tmp_City_3.Myfloat)
      retParams.AddressSet[i3].Myfloat = tmp_Myfloat_14      
      retParams.AddressSet[i3].Events = make([] int64, len(tmp_City_3.Events) )
      for j := 0; j < len(tmp_City_3.Events ); j++ { 
        retParams.AddressSet[i3].Events[j] = int64(tmp_City_3.Events[j])
      }      
      retParams.AddressSet[i3].Mymap = tmp_City_3.Mymap      
      retParams.AddressSet[i3].AddressList = tmp_City_3.AddressList
      }
    tmp_Simple_15, ok := codeGenRawTableResult["my_list"].([]map[string]interface{})
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for my_list", ok )
    }
    retParams.MyList = make([]*models.Simple, len(tmp_Simple_15))
    for i6, v6 := range tmp_Simple_15 {
    
      tmp_Simple_16 := &Simple{
    
          v6["dummy"].(string),
        }
          
      retParams.MyList[i6] = &models.Simple{}
      retParams.MyList[i6].Dummy = tmp_Simple_16.Dummy
      }
    retParams.Name = &Name
    Mediate = codeGenRawTableResult["mediate"].(time.Time)
    tmp_cassuservice_tmp_time_0 = Mediate.String()
    tmp_Mediate_17 := tmp_cassuservice_tmp_time_0[0:10] + "T" + tmp_cassuservice_tmp_time_0[11:19] + "." + tmp_cassuservice_tmp_time_0[20:22]
    if tmp_cassuservice_tmp_time_0[22] == ' ' {
      tmp_cassuservice_tmp_time_0 = tmp_Mediate_17 + "0" + "Z" 
    } else { 
      tmp_cassuservice_tmp_time_0 = tmp_Mediate_17 + "Z"
    }
    tmp_Mediate_18, _  := strfmt.ParseDateTime(tmp_cassuservice_tmp_time_0)
    tmp_Mediate_19 := tmp_Mediate_18.String()
    retParams.Mediate = &tmp_Mediate_19
    SecondTs = codeGenRawTableResult["second_ts"].(time.Time)
    tmp_cassuservice_tmp_time_0 = SecondTs.String()
    tmp_SecondTs_20 := tmp_cassuservice_tmp_time_0[0:10] + "T" + tmp_cassuservice_tmp_time_0[11:19] + "." + tmp_cassuservice_tmp_time_0[20:22]
    if tmp_cassuservice_tmp_time_0[22] == ' ' {
      tmp_cassuservice_tmp_time_0 = tmp_SecondTs_20 + "0" + "Z" 
    } else { 
      tmp_cassuservice_tmp_time_0 = tmp_SecondTs_20 + "Z"
    }
    tmp_SecondTs_21, _  := strfmt.ParseDateTime(tmp_cassuservice_tmp_time_0)
    tmp_SecondTs_22 := tmp_SecondTs_21.String()
    retParams.SecondTs = &tmp_SecondTs_22
    
    tmp_Tevents_23 := codeGenRawTableResult["tevents"].([]int)
    retParams.Tevents = make([] int64, len(tmp_Tevents_23) )
    for j := 0; j < len(tmp_Tevents_23 ); j++ { 
      retParams.Tevents[j] = int64(tmp_Tevents_23[j])
    }
    
    tmp_Tmylist_24 := codeGenRawTableResult["tmylist"].([]float32)
    retParams.Tmylist = make([] float64, len(tmp_Tmylist_24) )
    for j := 0; j < len(tmp_Tmylist_24 ); j++ { 
      retParams.Tmylist[j] = float64(tmp_Tmylist_24[j])
    }
    tmp_Tmymap_25, ok := codeGenRawTableResult["tmymap"].(map[string]string)
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for tmymap", ok )
    }
    retParams.Tmymap = make(map[string]string,len(tmp_Tmymap_25))
    for i13, v := range tmp_Tmymap_25 {
      retParams.Tmymap[i13] = v
    }
    return operations.NewGetEmployeeOK().WithPayload( payLoad.Payload)
    }`
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
	CreateCode( false, "/tmp", "github.com/stevef1uk/test4", parse1,  "",  "",  0, false , false , true   )


	// Read generated file
	path := "/tmp/data/" + MAINFILE
	fileo, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Expected to read %d bytes\n", len(expected))
	byteSlice := make([]byte, 10000)
	numBytesRead, err := io.ReadFull(fileo, byteSlice)
	if err != nil {
		log.Printf("Number of bytes read: %d\n", numBytesRead)
	}
	tmpbytes := byteSlice[0:numBytesRead]
	s := string(tmpbytes[:])

	if (len(expected) != len(s) ) {
		t.Errorf("Read %d bytes expected %d bytes\n", len(s), len(expected) )
	}

	for i, _ := range s {
		if (expected[i] != s[i] ) {
			t.Errorf("Difference at %d, got %c expected %c", i, expected[i], s[i] )
		}
	}

	//log.Printf("Expected bytes *%s*\n", expected)
	//log.Printf("Read bytes *%s*\n", s)

}





func Test1(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST1, EXPECTED_OUTPUT_TEST1, t )
	/*
	path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
	ret6 :=  SpiceInHandler( false , path, "Employee", "" )
	_ = ret6
	*/
}


/*
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
    myfloat float,
    events set<int>,
    mymap  map<text, text>,
    address_list set<frozen<simple>>
);

CREATE TABLE demo.employee (
    id int,
    address_set set<frozen<city>>,
    my_List list<frozen<simple>>,
    name text,
    mediate TIMESTAMP,
    second_ts TIMESTAMP,
    tevents set<int>,
    tmylist list<float>,
    tmymap  map<text, text>,
   PRIMARY KEY (id, mediate, second_ts )
 ) WITH CLUSTERING ORDER BY (mediate ASC, second_ts ASC)
`
 */
