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

	CSQ_TEST2 = `

    CREATE TABLE demo.accounts4 (
    id int,
    name text,
    ascii1 ascii,
    bint1 bigint,
    blob1 blob,
    bool1 boolean,
    dec1 decimal,
    double1 double,
    flt1 float,
    inet1 inet,
    int1 int,
    text1 text,
    time1 timestamp,
    time2 timeuuid,
    mydate1 date,
    uuid1 uuid,
    varchar1 varchar,
    events set<int>,
    mylist list<float>,
    myset set<text>,
    adec list<decimal>,
    PRIMARY KEY (id, name, time1)
) WITH CLUSTERING ORDER BY (name ASC)
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
    tmp_Name_17 := codeGenRawTableResult["name"].(string)
    retParams.Name = &tmp_Name_17
    Mediate = codeGenRawTableResult["mediate"].(time.Time)
    tmp_cassuservice_tmp_time_0 = Mediate.String()
    tmp_Mediate_18 := tmp_cassuservice_tmp_time_0[0:10] + "T" + tmp_cassuservice_tmp_time_0[11:19] + "." + tmp_cassuservice_tmp_time_0[20:22]
    if tmp_cassuservice_tmp_time_0[22] == ' ' {
      tmp_cassuservice_tmp_time_0 = tmp_Mediate_18 + "0" + "Z" 
    } else { 
      tmp_cassuservice_tmp_time_0 = tmp_Mediate_18 + "Z"
    }
    tmp_Mediate_19, _  := strfmt.ParseDateTime(tmp_cassuservice_tmp_time_0)
    tmp_Mediate_20 := tmp_Mediate_19.String()
    retParams.Mediate = &tmp_Mediate_20
    SecondTs = codeGenRawTableResult["second_ts"].(time.Time)
    tmp_cassuservice_tmp_time_0 = SecondTs.String()
    tmp_SecondTs_21 := tmp_cassuservice_tmp_time_0[0:10] + "T" + tmp_cassuservice_tmp_time_0[11:19] + "." + tmp_cassuservice_tmp_time_0[20:22]
    if tmp_cassuservice_tmp_time_0[22] == ' ' {
      tmp_cassuservice_tmp_time_0 = tmp_SecondTs_21 + "0" + "Z" 
    } else { 
      tmp_cassuservice_tmp_time_0 = tmp_SecondTs_21 + "Z"
    }
    tmp_SecondTs_22, _  := strfmt.ParseDateTime(tmp_cassuservice_tmp_time_0)
    tmp_SecondTs_23 := tmp_SecondTs_22.String()
    retParams.SecondTs = &tmp_SecondTs_23
    
    tmp_Tevents_24 := codeGenRawTableResult["tevents"].([]int)
    retParams.Tevents = make([] int64, len(tmp_Tevents_24) )
    for j := 0; j < len(tmp_Tevents_24 ); j++ { 
      retParams.Tevents[j] = int64(tmp_Tevents_24[j])
    }
    
    tmp_Tmylist_25 := codeGenRawTableResult["tmylist"].([]float32)
    retParams.Tmylist = make([] float64, len(tmp_Tmylist_25) )
    for j := 0; j < len(tmp_Tmylist_25 ); j++ { 
      retParams.Tmylist[j] = float64(tmp_Tmylist_25[j])
    }
    tmp_Tmymap_26, ok := codeGenRawTableResult["tmymap"].(map[string]string)
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for tmymap", ok )
    }
    retParams.Tmymap = make(map[string]string,len(tmp_Tmymap_26))
    for i13, v := range tmp_Tmymap_26 {
      retParams.Tmymap[i13] = v
    }
    return operations.NewGetEmployeeOK().WithPayload( payLoad.Payload)
    }`

	EXPECTED_OUTPUT_TEST2 = `// GENERATED FILE so do not edit or will be overwritten upon next generate
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
    "gopkg.in/inf.v0"
    "strconv"
)


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

func Search(params operations.GetAccounts4Params) middleware.Responder {

    var ID int64
    _ = ID
    var Name string
    _ = Name
    var Ascii1 string
    _ = Ascii1
    var Bint1 int64
    _ = Bint1
    var Blob1 string
    _ = Blob1
    var Bool1 bool
    _ = Bool1
    var Dec1 *inf.Dec
    _ = Dec1
    var Double1 float64
    _ = Double1
    var Flt1 float64
    _ = Flt1
    var Inet1 string
    _ = Inet1
    var Int1 int64
    _ = Int1
    var Text1 string
    _ = Text1
    var Time1 time.Time
    _ = Time1
    var Time2 time.Time
    _ = Time2
    var Mydate1 time.Time
    _ = Mydate1
    var Uuid1 string
    _ = Uuid1
    var Varchar1 string
    _ = Varchar1
    var Events []int64
    _ = Events
    var Mylist []float64
    _ = Mylist
    var Myset []string
    _ = Myset
    var Adec []*inf.Dec
    _ = Adec
    tmp_cassuservice_tmp_time_26 := strfmt.NewDateTime().String()

    codeGenRawTableResult := map[string]interface{}{}

    Time1,_ = time.Parse(time.RFC3339,params.Time1.String() ) 
    if err := cassuservice_session.Query(`+"`"+` SELECT id, name, ascii1, bint1, blob1, bool1, dec1, double1, flt1, inet1, int1, text1, time1, time2, mydate1, uuid1, varchar1, events, mylist, myset, adec FROM accounts4 WHERE id = ? and name = ? and time1 = ? `+"`"+`,params.ID,params.Name,Time1).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetAccounts4BadRequest()
    }
    payLoad := operations.NewGetAccounts4OK()
    payLoad.Payload = make([]*models.GetAccounts4OKBodyItems,1)
    payLoad.Payload[0] = new(models.GetAccounts4OKBodyItems)
    retParams := payLoad.Payload[0]
    tmp_ID_27 := codeGenRawTableResult["id"].(int)
    ID = int64(tmp_ID_27)
    retParams.ID = &ID
    retParams.Name = &Name
    retParams.Ascii1 = &Ascii1
    retParams.Bint1 = &Bint1
    retParams.Blob1 = &Blob1
    retParams.Bool1 = &Bool1
    retParams.Dec1 = &Dec1
    retParams.Double1 = &Double1
    tmp_Flt1_28 := codeGenRawTableResult["flt1"].(float32)
    Flt1 = float64(tmp_Flt1_28)
    retParams.Flt1 = &Flt1
    retParams.Inet1 = &Inet1
    tmp_Int1_29 := codeGenRawTableResult["int1"].(int)
    Int1 = int64(tmp_Int1_29)
    retParams.Int1 = &Int1
    retParams.Text1 = &Text1
    Time1 = codeGenRawTableResult["time1"].(time.Time)
    tmp_cassuservice_tmp_time_26 = Time1.String()
    tmp_Time1_30 := tmp_cassuservice_tmp_time_26[0:10] + "T" + tmp_cassuservice_tmp_time_26[11:19] + "." + tmp_cassuservice_tmp_time_26[20:22]
    if tmp_cassuservice_tmp_time_26[22] == ' ' {
      tmp_cassuservice_tmp_time_26 = tmp_Time1_30 + "0" + "Z" 
    } else { 
      tmp_cassuservice_tmp_time_26 = tmp_Time1_30 + "Z"
    }
    tmp_Time1_31, _  := strfmt.ParseDateTime(tmp_cassuservice_tmp_time_26)
    tmp_Time1_32 := tmp_Time1_31.String()
    retParams.Time1 = &tmp_Time1_32
    Time2 = codeGenRawTableResult["time2"].(time.Time)
    tmp_cassuservice_tmp_time_26 = Time2.String()
    tmp_Time2_33 := tmp_cassuservice_tmp_time_26[0:10] + "T" + tmp_cassuservice_tmp_time_26[11:19] + "." + tmp_cassuservice_tmp_time_26[20:22]
    if tmp_cassuservice_tmp_time_26[22] == ' ' {
      tmp_cassuservice_tmp_time_26 = tmp_Time2_33 + "0" + "Z" 
    } else { 
      tmp_cassuservice_tmp_time_26 = tmp_Time2_33 + "Z"
    }
    tmp_Time2_34, _  := strfmt.ParseDateTime(tmp_cassuservice_tmp_time_26)
    tmp_Time2_35 := tmp_Time2_34.String()
    retParams.Time2 = &tmp_Time2_35
    Mydate1 = codeGenRawTableResult["mydate1"].(time.Time)
    tmp_cassuservice_tmp_time_26 = Mydate1.String()
    tmp_Mydate1_36 := tmp_cassuservice_tmp_time_26[0:10] + "T" + tmp_cassuservice_tmp_time_26[11:19] + "." + tmp_cassuservice_tmp_time_26[20:22]
    if tmp_cassuservice_tmp_time_26[22] == ' ' {
      tmp_cassuservice_tmp_time_26 = tmp_Mydate1_36 + "0" + "Z" 
    } else { 
      tmp_cassuservice_tmp_time_26 = tmp_Mydate1_36 + "Z"
    }
    tmp_Mydate1_37, _  := strfmt.ParseDateTime(tmp_cassuservice_tmp_time_26)
    tmp_Mydate1_38 := tmp_Mydate1_37.String()
    retParams.Mydate1 = &tmp_Mydate1_38
    retParams.Uuid1 = &Uuid1
    retParams.Varchar1 = &Varchar1
    
    tmp_Events_39 := codeGenRawTableResult["events"].([]int)
    retParams.Events = make([] int64, len(tmp_Events_39) )
    for j := 0; j < len(tmp_Events_39 ); j++ { 
      retParams.Events[j] = int64(tmp_Events_39[j])
    }
    
    tmp_Mylist_40 := codeGenRawTableResult["mylist"].([]float32)
    retParams.Mylist = make([] float64, len(tmp_Mylist_40) )
    for j := 0; j < len(tmp_Mylist_40 ); j++ { 
      retParams.Mylist[j] = float64(tmp_Mylist_40[j])
    }
    
    tmp_Myset_41 := codeGenRawTableResult["myset"].([]string)
    retParams.Myset = make([] string, len(tmp_Myset_41) )
    for j := 0; j < len(tmp_Myset_41 ); j++ { 
    
    tmp_Adec_42 := codeGenRawTableResult["adec"].([]*inf.Dec)
    retParams.Adec = make([] *inf.Dec, len(tmp_Adec_42) )
    for j := 0; j < len(tmp_Adec_42 ); j++ { 
    return operations.NewGetAccounts4OK().WithPayload( payLoad.Payload)
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

	path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
	ret6 :=  SpiceInHandler( false , path, "Employee", "" )
	_ = ret6

}


func Test2(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST2, EXPECTED_OUTPUT_TEST2, t )

	path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
	ret6 :=  SpiceInHandler( false , path, "Accounts4", "" )
	_ = ret6
	
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
