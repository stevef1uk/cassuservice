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
// Insert test: insert into employee ( id, mediate, second_ts, name,  my_list, address_set  ) values (1, '2018-02-17T13:01:05.000Z', '1999-12-01T23:21:59.123Z', 'steve', [{dummy:'fred'}], {{id:1, mymap:{'a':'fred'}, citycode:'Peef',lastupdatedat:'2019-02-18T14:02:06.000Z',address_list:{{dummy:'foobar'}},events:{1,2,3} }} ) ;
// curl -X GET "http://127.0.0.1:5000/v1/employee?id=1&mediate=2018-02-17T13:01:05.000Z&second_ts=1999-12-01T23:21:59.123Z"

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
//insert into accounts4 ( id, name, time1, dec1, uuid1,time2)  values (1, 'steve', '2017-02-18T13:01:06.000Z', 1.23e40, 74e61f45-bff0-11e6-b8d5-843835632426,Now() );
//curl -X GET "http://127.0.0.1:5000/v1/accounts4?id=1&name=steve&time1=2017-02-18T13:01:06.000Z"

	CSQ_TEST3 = `
CREATE TYPE demo.simple (
       id int,
       dummy text,
       mediate TIMESTAMP
    );

CREATE TABLE demo.employee1 (
    id int PRIMARY KEY,
    tSimple simple
) WITH CLUSTERING ORDER BY (name ASC)
`

	CSQ_TEST4= `
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
`

	CSQ_TEST5= `
CREATE TABLE demo.demo1 (
                        id int PRIMARY KEY,
                        testtimestamp  timestamp,
                        testbigint     bigint,
                        testblob       blob,
                        testbool       boolean,
                        testfloat      float,
                        testdouble     double,
                        testint        int,
                        testlist       list<text>,
                        testset        set<int>,
                        testmap        map<text, text> )

`

	CSQ_TEST6= `
CREATE TYPE demo.simple3 (
       id int,
       floter float
    );

CREATE TYPE demo.simples (
       id int,
       dummy text,
       mediate TIMESTAMP,
       embedded list<frozen<simple3>>
    );

CREATE TABLE demo.employee1 (
    id int PRIMARY KEY,
    tSimple frozen <simples>
)
`

	CSQ_TEST7= `
CREATE TYPE demo.simple (
    id int,
    floter float
    );


CREATE TYPE demo.simple3 (
    id int,
    floter float,
    etype frozen <simple>
    );

CREATE TYPE demo.simples (
    id int,
    dummy text,
    mediate TIMESTAMP,
    embedded list<frozen<simple3>>
    );

CREATE TABLE demo.employee11 (
   id int PRIMARY KEY,
   tSimple list<frozen <simples>>
);
`

	CSQ_TEST8= `
CREATE TYPE demo.simple (
    id int,
    floter float
    );

CREATE TABLE demo.maptest1 (
      id int PRIMARY KEY,
      mymap map<text, frozen<simple>>
);
`
	CSQ_TEST9= `
CREATE TYPE demo.simple (
    id int,
    floter float
);

CREATE TABLE demo.maptest1 (
    id int PRIMARY KEY,
    mymap map<text, frozen<simple>>
);
`
//curl -d '{"id": 1, "mymap": {"b": {"id": 3, "floter": 2.2}}}' -H "Content-Type: application/json" -v -X POST http://localhost:5000/v1/maptest1

//insert into employee11 ( id, tsimple ) values (1, [{id:1,dummy:'hi',mediate:'018-02-17T13:01:05.000Z',embedded:[ {id:3,floter:2.23,etype:{id:4,floter:8.9} }, {id:10, floter:6.98,etype:{id:5,floter:7.1}}] }] );

// insert into employee1 (id, tsimple ) VALUES (1, {id:2,dummy:'text',mediate:'2017-01-18T13:01:06.000Z',estruct:{{id:4,citycode:'Peef'}}} );
//curl -X GET "http://127.0.0.1:5000/v1/employee1?id=1"
/*
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

if err := cassuservice_session.Query(`+"`"+` SELECT id, address_set, my_list, name, mediate, second_ts, tevents, tmylist, tmymap FROM employee WHERE id = ? and mediate = ? and second_ts = ? `+"`"+`,params.ID,Mediate,SecondTs).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
}


 */
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
     
)

func parseTime ( input string) time.Time {
    var ret time.Time
    if input == "" {
        ret = time.Now()
    } else {
        ret, _ = time.Parse( time.RFC3339, input )
    }
    return ret;
}

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

    codeGenRawTableResult := map[string]interface{}{}

    Mediate,_ = time.Parse(time.RFC3339,params.Mediate.String() ) 
    SecondTs,_ = time.Parse(time.RFC3339,params.SecondTs.String() ) 
    if err := cassuservice_session.Query(`+"`"+` SELECT id, address_set, my_list, name, mediate, second_ts, tevents, tmylist, tmymap FROM employee WHERE id = ? and mediate = ? and second_ts = ? `+"`"+`,params.ID,Mediate,SecondTs).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetEmployeeBadRequest()
    }
    payLoad := operations.NewGetEmployeeOK()
    payLoad.Payload = make([]*operations.GetEmployeeOKBodyItems0,1)
    payLoad.Payload[0] = new(operations.GetEmployeeOKBodyItems0)
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
      if v3["address_list"] == nil { 
          continue
      }
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
      tmp_Lastupdatedat_11 := tmp_City_3.Lastupdatedat.String()
      retParams.AddressSet[i3].Lastupdatedat = tmp_Lastupdatedat_11      
      tmp_Myfloat_12 := float64(tmp_City_3.Myfloat)
      retParams.AddressSet[i3].Myfloat = tmp_Myfloat_12      
      retParams.AddressSet[i3].Events = make([] int64, len(tmp_City_3.Events) )
      for j := 0; j < len(tmp_City_3.Events ); j++ { 
        retParams.AddressSet[i3].Events[j] = int64(tmp_City_3.Events[j])
      }      
      retParams.AddressSet[i3].Mymap = tmp_City_3.Mymap      
      retParams.AddressSet[i3].AddressList = tmp_City_3.AddressList
      }
    tmp_Simple_13, ok := codeGenRawTableResult["my_list"].([]map[string]interface{})
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for my_list", ok )
    }
    retParams.MyList = make([]*models.Simple, len(tmp_Simple_13))
    for i6, v6 := range tmp_Simple_13 {
    
      tmp_Simple_14 := &Simple{
    
          v6["dummy"].(string),
        }
          
      retParams.MyList[i6] = &models.Simple{}
      retParams.MyList[i6].Dummy = tmp_Simple_14.Dummy
      }
    tmp_Name_15 := codeGenRawTableResult["name"].(string)
    retParams.Name = &tmp_Name_15
    Mediate = codeGenRawTableResult["mediate"].(time.Time)
    tmp_Mediate_16 := Mediate.String()
    retParams.Mediate = &tmp_Mediate_16
    SecondTs = codeGenRawTableResult["second_ts"].(time.Time)
    tmp_SecondTs_17 := SecondTs.String()
    retParams.SecondTs = &tmp_SecondTs_17
    
    tmp_Tevents_18 := codeGenRawTableResult["tevents"].([]int)
    retParams.Tevents = make([] int64, len(tmp_Tevents_18) )
    for j := 0; j < len(tmp_Tevents_18 ); j++ { 
      retParams.Tevents[j] = int64(tmp_Tevents_18[j])
    }
    
    tmp_Tmylist_19 := codeGenRawTableResult["tmylist"].([]float32)
    retParams.Tmylist = make([] float64, len(tmp_Tmylist_19) )
    for j := 0; j < len(tmp_Tmylist_19 ); j++ { 
      retParams.Tmylist[j] = float64(tmp_Tmylist_19[j])
    }
    tmp_Tmymap_20, ok := codeGenRawTableResult["tmymap"].(map[string]string)
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for tmymap", ok )
    }
    retParams.Tmymap = make(map[string]string,len(tmp_Tmymap_20))
    for i13, v := range tmp_Tmymap_20 {
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
     
	"gopkg.in/inf.v0"
	"strconv"
	
)

func parseTime ( input string) time.Time {
    var ret time.Time
    if input == "" {
        ret = time.Now()
    } else {
        ret, _ = time.Parse( time.RFC3339, input )
    }
    return ret;
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
    var Time2 gocql.UUID
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

    codeGenRawTableResult := map[string]interface{}{}

    Time1,_ = time.Parse(time.RFC3339,params.Time1.String() ) 
    if err := cassuservice_session.Query(`+"`"+` SELECT id, name, ascii1, bint1, blob1, bool1, dec1, double1, flt1, inet1, int1, text1, time1, time2, mydate1, uuid1, varchar1, events, mylist, myset, adec FROM accounts4 WHERE id = ? and name = ? and time1 = ? `+"`"+`,params.ID,params.Name,Time1).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetAccounts4BadRequest()
    }
    payLoad := operations.NewGetAccounts4OK()
    payLoad.Payload = make([]*operations.GetAccounts4OKBodyItems0,1)
    payLoad.Payload[0] = new(operations.GetAccounts4OKBodyItems0)
    retParams := payLoad.Payload[0]
    tmp_ID_1 := codeGenRawTableResult["id"].(int)
    ID = int64(tmp_ID_1)
    retParams.ID = &ID
    tmp_Name_2 := codeGenRawTableResult["name"].(string)
    retParams.Name = &tmp_Name_2
    retParams.Ascii1 = &Ascii1
    tmp_Bint1_3 := codeGenRawTableResult["bint1"].(int64)
    retParams.Bint1 = &tmp_Bint1_3
    tmp_Blob1_4 := codeGenRawTableResult["blob1"].([]uint8)
    Blob1 = string(tmp_Blob1_4)
    retParams.Blob1 = &Blob1
    tmp_Bool1_5 := codeGenRawTableResult["bool1"].(bool)
    retParams.Bool1 = &tmp_Bool1_5
    Dec1 = codeGenRawTableResult["dec1"].(*inf.Dec)
    tmp_Dec1_6,_ := strconv.ParseFloat(Dec1.String(), 64 )
    retParams.Dec1 = &tmp_Dec1_6
    retParams.Double1 = &Double1
    tmp_Flt1_7 := codeGenRawTableResult["flt1"].(float32)
    Flt1 = float64(tmp_Flt1_7)
    retParams.Flt1 = &Flt1
    retParams.Inet1 = &Inet1
    tmp_Int1_8 := codeGenRawTableResult["int1"].(int)
    Int1 = int64(tmp_Int1_8)
    retParams.Int1 = &Int1
    tmp_Text1_9 := codeGenRawTableResult["text1"].(string)
    retParams.Text1 = &tmp_Text1_9
    Time1 = codeGenRawTableResult["time1"].(time.Time)
    tmp_Time1_10 := Time1.String()
    retParams.Time1 = &tmp_Time1_10
    Time2 = codeGenRawTableResult["time2"].(gocql.UUID)
    tmp_Time2_11 := Time2.String()
    retParams.Time2 = &tmp_Time2_11
    Mydate1 = codeGenRawTableResult["mydate1"].(time.Time)
    tmp_Mydate1_12 := Mydate1.String()
    retParams.Mydate1 = &tmp_Mydate1_12
    tmp_Uuid1_13 := codeGenRawTableResult["uuid1"].(gocql.UUID)
    Uuid1 = tmp_Uuid1_13.String()
    retParams.Uuid1 = &Uuid1
    retParams.Varchar1 = &Varchar1
    
    tmp_Events_14 := codeGenRawTableResult["events"].([]int)
    retParams.Events = make([] int64, len(tmp_Events_14) )
    for j := 0; j < len(tmp_Events_14 ); j++ { 
      retParams.Events[j] = int64(tmp_Events_14[j])
    }
    
    tmp_Mylist_15 := codeGenRawTableResult["mylist"].([]float32)
    retParams.Mylist = make([] float64, len(tmp_Mylist_15) )
    for j := 0; j < len(tmp_Mylist_15 ); j++ { 
      retParams.Mylist[j] = float64(tmp_Mylist_15[j])
    }
    
    tmp_Myset_16 := codeGenRawTableResult["myset"].([]string)
    retParams.Myset = make([] string, len(tmp_Myset_16) )
    for j := 0; j < len(tmp_Myset_16 ); j++ { 
      retParams.Myset[j] = tmp_Myset_16[j]
    }
    
    tmp_Adec_17 := codeGenRawTableResult["adec"].([]*inf.Dec)
    retParams.Adec = make([] float64, len(tmp_Adec_17) )
    for j := 0; j < len(tmp_Adec_17 ); j++ { 
      tmp_tmp_Adec_17_18,_ := strconv.ParseFloat(tmp_Adec_17[j].String(), 64 )
      retParams.Adec[j] = tmp_tmp_Adec_17_18
    }
    return operations.NewGetAccounts4OK().WithPayload( payLoad.Payload)
    }`

	EXPECTED_OUTPUT_TEST3 = `// GENERATED FILE so do not edit or will be overwritten upon next generate
package data

import (
    "github.com/stevef1uk/test4/models"
    "github.com/stevef1uk/test4/restapi/operations"
    middleware "github.com/go-openapi/runtime/middleware"
    "github.com/gocql/gocql"
    "os"
    "log"
    "time"
     
)

func parseTime ( input string) time.Time {
    var ret time.Time
    if input == "" {
        ret = time.Now()
    } else {
        ret, _ = time.Parse( time.RFC3339, input )
    }
    return ret;
}

type Simple struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Dummy string `+"`"+`cql:"dummy"`+"`"+`
    Mediate time.Time `+"`"+`cql:"mediate"`+"`"+`
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

func Search(params operations.GetEmployee1Params) middleware.Responder {

    var ID int64
    _ = ID
    Tsimple := &Simple{}
    _ = Tsimple

    codeGenRawTableResult := map[string]interface{}{}

    if err := cassuservice_session.Query(`+"`"+` SELECT id, tsimple FROM employee1 WHERE id = ? `+"`"+`,params.ID).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetEmployee1BadRequest()
    }
    payLoad := operations.NewGetEmployee1OK()
    payLoad.Payload = make([]*operations.GetEmployee1OKBodyItems0,1)
    payLoad.Payload[0] = new(operations.GetEmployee1OKBodyItems0)
    retParams := payLoad.Payload[0]
    tmp_ID_1 := codeGenRawTableResult["id"].(int)
    ID = int64(tmp_ID_1)
    retParams.ID = &ID
    tmp_TSIMPLE_2, ok := codeGenRawTableResult["tsimple"].(map[string]interface{})
    tmp_TSIMPLE_3 := &models.Simple{}
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for tsimple", ok )
    }
    
    tmp_Simple_5 := &Simple{
    
        tmp_TSIMPLE_2["id"].(int),
        tmp_TSIMPLE_2["dummy"].(string),
        tmp_TSIMPLE_2["mediate"].(time.Time),
      }
        
    tmp_ID_6 := int64(tmp_Simple_5.ID)
    tmp_TSIMPLE_3.ID = tmp_ID_6    
    tmp_TSIMPLE_3.Dummy = tmp_Simple_5.Dummy    
    tmp_Mediate_7 := tmp_Simple_5.Mediate.String()
    tmp_TSIMPLE_3.Mediate = tmp_Mediate_7
    retParams.Tsimple = tmp_TSIMPLE_3
    return operations.NewGetEmployee1OK().WithPayload( payLoad.Payload)
    }`

	EXPECTED_OUTPUT_TEST4 = `// GENERATED FILE so do not edit or will be overwritten upon next generate
package data

import (
    "github.com/stevef1uk/test4/models"
    "github.com/stevef1uk/test4/restapi/operations"
    middleware "github.com/go-openapi/runtime/middleware"
    "github.com/gocql/gocql"
    "os"
    "log"
    "time"
     
)

func parseTime ( input string) time.Time {
    var ret time.Time
    if input == "" {
        ret = time.Now()
    } else {
        ret, _ = time.Parse( time.RFC3339, input )
    }
    return ret;
}

type Simple1 struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Citycode string `+"`"+`cql:"citycode"`+"`"+`
}

type Simple struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Dummy string `+"`"+`cql:"dummy"`+"`"+`
    Mediate time.Time `+"`"+`cql:"mediate"`+"`"+`
    Estruct models.SimpleEstruct `+"`"+`cql:"estruct"`+"`"+`
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

func Search(params operations.GetEmployee1Params) middleware.Responder {

    var ID int64
    _ = ID
    Tsimple := &Simple{}
    _ = Tsimple

    codeGenRawTableResult := map[string]interface{}{}

    if err := cassuservice_session.Query(`+"`"+` SELECT id, tsimple FROM employee1 WHERE id = ? `+"`"+`,params.ID).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetEmployee1BadRequest()
    }
    payLoad := operations.NewGetEmployee1OK()
    payLoad.Payload = make([]*operations.GetEmployee1OKBodyItems0,1)
    payLoad.Payload[0] = new(operations.GetEmployee1OKBodyItems0)
    retParams := payLoad.Payload[0]
    tmp_ID_1 := codeGenRawTableResult["id"].(int)
    ID = int64(tmp_ID_1)
    retParams.ID = &ID
    tmp_TSIMPLE_2, ok := codeGenRawTableResult["tsimple"].(map[string]interface{})
    tmp_TSIMPLE_3 := &models.Simple{}
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for tsimple", ok )
    }
    
    tmp_estruct_6:= tmp_TSIMPLE_2["estruct"].([]map[string]interface{})
    tmp_estruct_7:= make(models.SimpleEstruct, len(tmp_estruct_6) )
    
        for i3, v3 := range tmp_estruct_6 {
    
        tmp_Simple1_8 := &Simple1{
    
              v3["id"].(int),
                v3["citycode"].(string),
              }
                
            tmp_estruct_7[i3] = &models.Simple1{}
            tmp_ID_9 := int64(tmp_Simple1_8.ID)
            tmp_estruct_7[i3].ID = tmp_ID_9            
            tmp_estruct_7[i3].Citycode = tmp_Simple1_8.Citycode
            }
    tmp_Simple_5 := &Simple{
    
        tmp_TSIMPLE_2["id"].(int),
        tmp_TSIMPLE_2["dummy"].(string),
        tmp_TSIMPLE_2["mediate"].(time.Time),
        tmp_estruct_7,
      }
        
    tmp_ID_10 := int64(tmp_Simple_5.ID)
    tmp_TSIMPLE_3.ID = tmp_ID_10    
    tmp_TSIMPLE_3.Dummy = tmp_Simple_5.Dummy    
    tmp_Mediate_11 := tmp_Simple_5.Mediate.String()
    tmp_TSIMPLE_3.Mediate = tmp_Mediate_11    
    tmp_TSIMPLE_3.Estruct = tmp_Simple_5.Estruct
    retParams.Tsimple = tmp_TSIMPLE_3
    return operations.NewGetEmployee1OK().WithPayload( payLoad.Payload)
    }`

	EXPECTED_OUTPUT_TEST5 = `// GENERATED FILE so do not edit or will be overwritten upon next generate
package data

import (
    "github.com/stevef1uk/test4/models"
    "github.com/stevef1uk/test4/restapi/operations"
    middleware "github.com/go-openapi/runtime/middleware"
    "github.com/gocql/gocql"
    "os"
    "log"
    "time"
     
    "github.com/stevef1uk/test4/restapi/operations/demo1"
     
    "fmt"
    "strconv"
     
)

func parseTime ( input string) time.Time {
    var ret time.Time
    if input == "" {
        ret = time.Now()
    } else {
        ret, _ = time.Parse( time.RFC3339, input )
    }
    return ret;
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

func Search(params operations.GetDemo1Params) middleware.Responder {

    var ID int64
    _ = ID
    var Testtimestamp time.Time
    _ = Testtimestamp
    var Testbigint int64
    _ = Testbigint
    var Testblob string
    _ = Testblob
    var Testbool bool
    _ = Testbool
    var Testfloat float64
    _ = Testfloat
    var Testdouble float64
    _ = Testdouble
    var Testint int64
    _ = Testint
    var Testlist []string
    _ = Testlist
    var Testset []int64
    _ = Testset
    var Testmap []string
    _ = Testmap
    _ = models.Demo1{}

    codeGenRawTableResult := map[string]interface{}{}

    if err := cassuservice_session.Query(`+"`"+` SELECT id, testtimestamp, testbigint, testblob, testbool, testfloat, testdouble, testint, testlist, testset, testmap FROM demo1 WHERE id = ? `+"`"+`,params.ID).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetDemo1BadRequest()
    }
    payLoad := operations.NewGetDemo1OK()
    payLoad.Payload = make([]*operations.GetDemo1OKBodyItems0,1)
    payLoad.Payload[0] = new(operations.GetDemo1OKBodyItems0)
    retParams := payLoad.Payload[0]
    tmp_ID_1 := codeGenRawTableResult["id"].(int)
    ID = int64(tmp_ID_1)
    retParams.ID = &ID
    Testtimestamp = codeGenRawTableResult["testtimestamp"].(time.Time)
    tmp_Testtimestamp_2 := Testtimestamp.String()
    retParams.Testtimestamp = &tmp_Testtimestamp_2
    tmp_Testbigint_3 := codeGenRawTableResult["testbigint"].(int64)
    retParams.Testbigint = &tmp_Testbigint_3
    tmp_Testblob_4 := codeGenRawTableResult["testblob"].([]uint8)
    Testblob = string(tmp_Testblob_4)
    retParams.Testblob = &Testblob
    tmp_Testbool_5 := codeGenRawTableResult["testbool"].(bool)
    retParams.Testbool = &tmp_Testbool_5
    tmp_Testfloat_6 := codeGenRawTableResult["testfloat"].(float32)
    Testfloat = float64(tmp_Testfloat_6)
    retParams.Testfloat = &Testfloat
    retParams.Testdouble = &Testdouble
    tmp_Testint_7 := codeGenRawTableResult["testint"].(int)
    Testint = int64(tmp_Testint_7)
    retParams.Testint = &Testint
    
    tmp_Testlist_8 := codeGenRawTableResult["testlist"].([]string)
    retParams.Testlist = make([] string, len(tmp_Testlist_8) )
    for j := 0; j < len(tmp_Testlist_8 ); j++ { 
      retParams.Testlist[j] = tmp_Testlist_8[j]
    }
    
    tmp_Testset_9 := codeGenRawTableResult["testset"].([]int)
    retParams.Testset = make([] int64, len(tmp_Testset_9) )
    for j := 0; j < len(tmp_Testset_9 ); j++ { 
      retParams.Testset[j] = int64(tmp_Testset_9[j])
    }
    tmp_Testmap_10, ok := codeGenRawTableResult["testmap"].(map[string]string)
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for testmap", ok )
    }
    retParams.Testmap = make(map[string]string,len(tmp_Testmap_10))
    for i12, v := range tmp_Testmap_10 {
      retParams.Testmap[i12] = v
    }
    return operations.NewGetDemo1OK().WithPayload( payLoad.Payload)
    }

func Insert(params demo1.AddDemo1Params) middleware.Responder {

    m := make(map[string]interface{})
    
    
    m["id"] = params.Body.ID
    m["testtimestamp"] = parseTime(params.Body.Testtimestamp)
    m["testbigint"] = params.Body.Testbigint
    m["testblob"] = params.Body.Testblob
    m["testbool"] = params.Body.Testbool
    tmp_testfloat_11:= fmt.Sprintf("%f",params.Body.Testfloat)
    tmp_testfloat_12,_ := strconv.ParseFloat(tmp_testfloat_11,32)
    m["testfloat"] = float32(tmp_testfloat_12)
    m["testdouble"] = params.Body.Testdouble
    m["testint"] = params.Body.Testint
    m["testlist"] = params.Body.Testlist
    m["testset"] = params.Body.Testset
    m["testmap"] = params.Body.Testmap
    if err := cassuservice_session.Query(`+"`"+` INSERT INTO demo1(id, testtimestamp, testbigint, testblob, testbool, testfloat, testdouble, testint, testlist, testset, testmap) VALUES (?,?,?,?,?,?,?,?,?,?,?)`+"`"+`,m["id"],m["testtimestamp"],m["testbigint"],m["testblob"],m["testbool"],m["testfloat"],m["testdouble"],m["testint"],m["testlist"],m["testset"],m["testmap"]).Consistency(gocql.One).Exec(); err != nil {
      return demo1.NewAddDemo1MethodNotAllowed()
    }
    return demo1.NewAddDemo1Created()
}`

	EXPECTED_OUTPUT_TEST6 = `// GENERATED FILE so do not edit or will be overwritten upon next generate
package data

import (
    "github.com/stevef1uk/test4/models"
    "github.com/stevef1uk/test4/restapi/operations"
    middleware "github.com/go-openapi/runtime/middleware"
    "github.com/gocql/gocql"
    "os"
    "log"
    "time"
     
)

func parseTime ( input string) time.Time {
    var ret time.Time
    if input == "" {
        ret = time.Now()
    } else {
        ret, _ = time.Parse( time.RFC3339, input )
    }
    return ret;
}

type Simple3 struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Floter float32 `+"`"+`cql:"floter"`+"`"+`
}

type Simples struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Dummy string `+"`"+`cql:"dummy"`+"`"+`
    Mediate time.Time `+"`"+`cql:"mediate"`+"`"+`
    Embedded models.SimplesEmbedded `+"`"+`cql:"embedded"`+"`"+`
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

func Search(params operations.GetEmployee1Params) middleware.Responder {

    var ID int64
    _ = ID
    Tsimple := &Simples{}
    _ = Tsimple

    codeGenRawTableResult := map[string]interface{}{}

    if err := cassuservice_session.Query(`+"`"+` SELECT id, tsimple FROM employee1 WHERE id = ? `+"`"+`,params.ID).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetEmployee1BadRequest()
    }
    payLoad := operations.NewGetEmployee1OK()
    payLoad.Payload = make([]*operations.GetEmployee1OKBodyItems0,1)
    payLoad.Payload[0] = new(operations.GetEmployee1OKBodyItems0)
    retParams := payLoad.Payload[0]
    tmp_ID_1 := codeGenRawTableResult["id"].(int)
    ID = int64(tmp_ID_1)
    retParams.ID = &ID
    tmp_TSIMPLE_2, ok := codeGenRawTableResult["tsimple"].(map[string]interface{})
    tmp_TSIMPLE_3 := &models.Simples{}
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for tsimple", ok )
    }
    
    tmp_embedded_6:= tmp_TSIMPLE_2["embedded"].([]map[string]interface{})
    tmp_embedded_7:= make(models.SimplesEmbedded, len(tmp_embedded_6) )
    
        for i3, v3 := range tmp_embedded_6 {
    
        tmp_Simple3_8 := &Simple3{
    
              v3["id"].(int),
                v3["floter"].(float32),
              }
                
            tmp_embedded_7[i3] = &models.Simple3{}
            tmp_ID_9 := int64(tmp_Simple3_8.ID)
            tmp_embedded_7[i3].ID = tmp_ID_9            
            tmp_Floter_10 := float64(tmp_Simple3_8.Floter)
            tmp_embedded_7[i3].Floter = tmp_Floter_10
            }
    tmp_Simples_5 := &Simples{
    
        tmp_TSIMPLE_2["id"].(int),
        tmp_TSIMPLE_2["dummy"].(string),
        tmp_TSIMPLE_2["mediate"].(time.Time),
        tmp_embedded_7,
      }
        
    tmp_ID_11 := int64(tmp_Simples_5.ID)
    tmp_TSIMPLE_3.ID = tmp_ID_11    
    tmp_TSIMPLE_3.Dummy = tmp_Simples_5.Dummy    
    tmp_Mediate_12 := tmp_Simples_5.Mediate.String()
    tmp_TSIMPLE_3.Mediate = tmp_Mediate_12    
    tmp_TSIMPLE_3.Embedded = tmp_Simples_5.Embedded
    retParams.Tsimple = tmp_TSIMPLE_3
    return operations.NewGetEmployee1OK().WithPayload( payLoad.Payload)
    }`

	EXPECTED_OUTPUT_TEST7 = `// GENERATED FILE so do not edit or will be overwritten upon next generate
package data

import (
    "github.com/stevef1uk/test4/models"
    "github.com/stevef1uk/test4/restapi/operations"
    middleware "github.com/go-openapi/runtime/middleware"
    "github.com/gocql/gocql"
    "os"
    "log"
    "time"
     
)

func parseTime ( input string) time.Time {
    var ret time.Time
    if input == "" {
        ret = time.Now()
    } else {
        ret, _ = time.Parse( time.RFC3339, input )
    }
    return ret;
}

type Simple struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Floter float32 `+"`"+`cql:"floter"`+"`"+`
}

type Simple3 struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Floter float32 `+"`"+`cql:"floter"`+"`"+`
    Etype Simple `+"`"+`cql:"etype"`+"`"+`
}

type Simples struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Dummy string `+"`"+`cql:"dummy"`+"`"+`
    Mediate time.Time `+"`"+`cql:"mediate"`+"`"+`
    Embedded models.SimplesEmbedded `+"`"+`cql:"embedded"`+"`"+`
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

func Search(params operations.GetEmployee11Params) middleware.Responder {

    var ID int64
    _ = ID
    var Tsimple models.Tsimple
    _ = Tsimple

    codeGenRawTableResult := map[string]interface{}{}

    if err := cassuservice_session.Query(`+"`"+` SELECT id, tsimple FROM employee11 WHERE id = ? `+"`"+`,params.ID).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetEmployee11BadRequest()
    }
    payLoad := operations.NewGetEmployee11OK()
    payLoad.Payload = make([]*operations.GetEmployee11OKBodyItems0,1)
    payLoad.Payload[0] = new(operations.GetEmployee11OKBodyItems0)
    retParams := payLoad.Payload[0]
    tmp_ID_1 := codeGenRawTableResult["id"].(int)
    ID = int64(tmp_ID_1)
    retParams.ID = &ID
    tmp_Simples_2, ok := codeGenRawTableResult["tsimple"].([]map[string]interface{})
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for tsimple", ok )
    }
    retParams.Tsimple = make([]*models.Simples, len(tmp_Simples_2))
    for i3, v3 := range tmp_Simples_2 {
    
      if v3["embedded"] == nil { 
          continue
      }
      tmp_embedded_4:= v3["embedded"].([]map[string]interface{})
      tmp_embedded_5:= make(models.SimplesEmbedded, len(tmp_embedded_4) )
      
          for i4, v4 := range tmp_embedded_4 {
    
              tmp_etype_7 := v4["etype"].(map[string]interface{})
              tmp_etype_8 := &models.Simple{}
                  
    
                  tmp_Simple_9 := Simple{
    
                        tmp_etype_7["id"].(int),
                          tmp_etype_7["floter"].(float32),
                        }
                          
                      tmp_ID_10 := int64(tmp_Simple_9.ID)
                      tmp_etype_8.ID = tmp_ID_10                      
                      tmp_Floter_11 := float64(tmp_Simple_9.Floter)
                      tmp_etype_8.Floter = tmp_Floter_11
          tmp_Simple3_6 := &Simple3{
    
                v4["id"].(int),
                  v4["floter"].(float32),
                  tmp_Simple_9,
                }
                  
              tmp_embedded_5[i4] = &models.Simple3{}
              tmp_ID_12 := int64(tmp_Simple3_6.ID)
              tmp_embedded_5[i4].ID = tmp_ID_12              
              tmp_Floter_13 := float64(tmp_Simple3_6.Floter)
              tmp_embedded_5[i4].Floter = tmp_Floter_13              
                tmp_embedded_5[i4].Etype = &models.Simple{}
                tmp_embedded_5[i4].Etype.ID = int64(tmp_Simple3_6.ID)
                tmp_embedded_5[i4].Etype.Floter = float64(tmp_Simple3_6.Floter)
              }
      tmp_Simples_3 := &Simples{
    
          v3["id"].(int),
          v3["dummy"].(string),
          v3["mediate"].(time.Time),
          tmp_embedded_5,
        }
          
      retParams.Tsimple[i3] = &models.Simples{}
      tmp_ID_14 := int64(tmp_Simples_3.ID)
      retParams.Tsimple[i3].ID = tmp_ID_14      
      retParams.Tsimple[i3].Dummy = tmp_Simples_3.Dummy      
      tmp_Mediate_15 := tmp_Simples_3.Mediate.String()
      retParams.Tsimple[i3].Mediate = tmp_Mediate_15      
      retParams.Tsimple[i3].Embedded = tmp_Simples_3.Embedded
      }
    return operations.NewGetEmployee11OK().WithPayload( payLoad.Payload)
    }`

	EXPECTED_OUTPUT_TEST8=`// GENERATED FILE so do not edit or will be overwritten upon next generate
package data

import (
    "github.com/stevef1uk/test4/models"
    "github.com/stevef1uk/test4/restapi/operations"
    middleware "github.com/go-openapi/runtime/middleware"
    "github.com/gocql/gocql"
    "os"
    "log"
)

type Simple struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Floter float32 `+"`"+`cql:"floter"`+"`"+`
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

func Search(params operations.GetMaptest1Params) middleware.Responder {

    var ID int64
    _ = ID
    var Mymap []string
    _ = Mymap

    codeGenRawTableResult := map[string]interface{}{}

    if err := cassuservice_session.Query(`+"`"+` SELECT id, mymap FROM maptest1 WHERE id = ? `+"`"+`,params.ID).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetMaptest1BadRequest()
    }
    payLoad := operations.NewGetMaptest1OK()
    payLoad.Payload = make([]*operations.GetMaptest1OKBodyItems0,1)
    payLoad.Payload[0] = new(operations.GetMaptest1OKBodyItems0)
    retParams := payLoad.Payload[0]
    tmp_ID_0 := codeGenRawTableResult["id"].(int)
    ID = int64(tmp_ID_0)
    retParams.ID = &ID
    tmp_Mymap_1, ok := codeGenRawTableResult["mymap"].(map[string]map[string]interface{})
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for mymap", ok )
    }
    retParams.Mymap = make(map[string]models.Simple,len(tmp_Mymap_1))
    for i3, v := range tmp_Mymap_1 {
    
      tmp_Mymap_2 := Simple{}
      tmp_Mymap_2.ID = v["id"].(int)
      tmp_Mymap_2.Floter = v["floter"].(float32)
      tmp_Mymap_3 := models.Simple{}
      tmp_Mymap_3.ID = int64(tmp_Mymap_2.ID)
      tmp_Mymap_3.Floter = float64(tmp_Mymap_2.Floter)
      retParams.Mymap[i3] = tmp_Mymap_3
    }
    return operations.NewGetMaptest1OK().WithPayload( payLoad.Payload)
    }`

	EXPECTED_OUTPUT_TEST9=`// GENERATED FILE so do not edit or will be overwritten upon next generate
package data

import (
    "github.com/stevef1uk/test4/models"
    "github.com/stevef1uk/test4/restapi/operations"
    middleware "github.com/go-openapi/runtime/middleware"
    "github.com/gocql/gocql"
    "os"
    "log"
    "github.com/stevef1uk/test4/restapi/operations/maptest1"
     
)

type Simple struct {
    ID int `+"`"+`cql:"id"`+"`"+`
    Floter float32 `+"`"+`cql:"floter"`+"`"+`
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

func Search(params operations.GetMaptest1Params) middleware.Responder {

    var ID int64
    _ = ID
    var Mymap []string
    _ = Mymap
    _ = models.Maptest1{}

    codeGenRawTableResult := map[string]interface{}{}

    if err := cassuservice_session.Query(`+"`"+` SELECT id, mymap FROM maptest1 WHERE id = ? `+"`"+`,params.ID).Consistency(gocql.One).MapScan(codeGenRawTableResult); err != nil {
      log.Println("No data? ", err)
      return operations.NewGetMaptest1BadRequest()
    }
    payLoad := operations.NewGetMaptest1OK()
    payLoad.Payload = make([]*operations.GetMaptest1OKBodyItems0,1)
    payLoad.Payload[0] = new(operations.GetMaptest1OKBodyItems0)
    retParams := payLoad.Payload[0]
    tmp_ID_0 := codeGenRawTableResult["id"].(int)
    ID = int64(tmp_ID_0)
    retParams.ID = &ID
    tmp_Mymap_1, ok := codeGenRawTableResult["mymap"].(map[string]map[string]interface{})
    if ! ok {
      log.Fatal("handleReturnedVar() - failed to find entry for mymap", ok )
    }
    retParams.Mymap = make(map[string]models.Simple,len(tmp_Mymap_1))
    for i3, v := range tmp_Mymap_1 {
    
      tmp_Mymap_2 := Simple{}
      tmp_Mymap_2.ID = v["id"].(int)
      tmp_Mymap_2.Floter = v["floter"].(float32)
      tmp_Mymap_3 := models.Simple{}
      tmp_Mymap_3.ID = int64(tmp_Mymap_2.ID)
      tmp_Mymap_3.Floter = float64(tmp_Mymap_2.Floter)
      retParams.Mymap[i3] = tmp_Mymap_3
    }
    return operations.NewGetMaptest1OK().WithPayload( payLoad.Payload)
    }

func Insert(params maptest1.AddMaptest1Params) middleware.Responder {

    m := make(map[string]interface{})
    
    
    m["id"] = params.Body.ID
    tmp_Mymap_4 := make( map[string]Simple, len(params.Body.Mymap) )
    m["mymap"] = tmp_Mymap_4
    for imymap,vmymap := range params.Body.Mymap{
    
        tmp_Simple_5 := Simple{
            int(vmymap.ID),
            float32(vmymap.Floter),
        }
        tmp_Mymap_4[imymap] = tmp_Simple_5
    }
    if err := cassuservice_session.Query(`+"`"+` INSERT INTO maptest1(id, mymap) VALUES (?,?)`+"`"+`,m["id"],m["mymap"]).Consistency(gocql.One).Exec(); err != nil {
      return maptest1.NewAddMaptest1MethodNotAllowed()
    }
    return maptest1.NewAddMaptest1Created()
}`

)

func performCreateTest1( debug bool, test string, cql string, expected string , t *testing.T, addPost bool) {

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
	CreateCode( false, "/tmp", "github.com/stevef1uk/test4", parse1,  "",  "",  0, false , true, addPost   )


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
	performCreateTest1(true, "Test1", CSQ_TEST1, EXPECTED_OUTPUT_TEST1, t, false )
/*
	path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
	ret6 :=  SpiceInHandler( false , path, "Employee", "" )
	_ = ret6
*/
}




func Test2(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST2, EXPECTED_OUTPUT_TEST2, t, false )
/*
	path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
	ret6 :=  SpiceInHandler( false , path, "Accounts4", "" )
	_ = ret6
*/
}

func Test3(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST3, EXPECTED_OUTPUT_TEST3, t, false )
/*
		path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
		ret6 :=  SpiceInHandler( false , path, "Employee1", "" )
		_ = ret6
*/
}

func Test4(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST4, EXPECTED_OUTPUT_TEST4, t, false )
/*
			path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
			ret6 :=  SpiceInHandler( false , path, "Employee1", "" )
			_ = ret6
*/
}


func Test5(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST5, EXPECTED_OUTPUT_TEST5, t, true )
	/*
				path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
				ret6 :=  SpiceInHandler( false , path, "Employee1", "" )
				_ = ret6
	*/
}

func Test6(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST6, EXPECTED_OUTPUT_TEST6, t, false )
	/*
		path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
		ret6 :=  SpiceInHandler( false , path, "Employee1", "" )
		_ = ret6
	*/
}

func Test7(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST7, EXPECTED_OUTPUT_TEST7, t, false )
	/*
		path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
		ret6 :=  SpiceInHandler( false , path, "Employee1", "" )
		_ = ret6
	*/
}

func Test8(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST8, EXPECTED_OUTPUT_TEST8, t, false )
	/*
		path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
		ret6 :=  SpiceInHandler( false , path, "Employee1", "" )
		_ = ret6
	*/
}

func Test9(t *testing.T) {
	performCreateTest1(true, "Test1", CSQ_TEST9, EXPECTED_OUTPUT_TEST9, t, true )
	/*
		path := os.Getenv("GOPATH")  + "/src/github.com/stevef1uk/test4/"
		ret6 :=  SpiceInHandler( false , path, "Employee1", "" )
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
