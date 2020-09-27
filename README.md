# cassuservice 
[![Open in Gitpod](https://gitpod.io/button/open-in-gitpod.svg)](https://gitpod.io/#https://github.com/stevef1uk/cassuservice)
If you click the above link go-swagger has been added to the default GitPod docker image. For the impatient take a look at the run-test.sh

Repo that contains Golang code that can autogenerate a Golang service from a Cassandra DDL for a table.
This tool can genenerate a service that can read and insert into a Cassandra table.

I originally wrote equivalent code whilst at my last employer as an exercise to learning Golang. I am no longer a professional programmer, but still enjoy the challenges of coding. I did gain permission to open source my first working version of the previous  work, but was unable to get the code officially listed before I decided to move on. Therefore, I have rewritten it in this repo and done it better the second time I believe. 

The architecture approach is simple:

1. Parse Cassandra DDL to create a swagger file
2. Use go-swagger to generate the RESTful server, see: https://github.com/go-swagger/go-swagger (Tested against v0.19.0)
3. Parse Cassandra DDL to create the Cassandra handler and wire it into the RESTful server.

I have written a short blog providing a little more information: https://stevef1uk.blogspot.com/2019/02/an-easy-way-to-build-restful-micro.html

Prerequisites:
I have updated the main.go file to perform the above steps in one. However, on my machine I have run the 'swagger generate server -f t.cql' command before and then followed the instructions to run 'go get -u X' as instructed, which will need to be done once manually to load the required packages. I have also installed gocql see: https://github.com/gocql/gocql as the generated Cassandra handler uses this.

Step 1: Export the Cassandra DDL into a file e.g. t.cql. This file needs to only contain the types (if any) and table definition for a single table. The best way to do this is from cqlsh and use the describe table command.
```
A very simple example t.cql is as follows:

    CREATE TYPE demo.simple (
           id int,
           dummy text,
           mediate TIMESTAMP
        );
    
    CREATE TABLE demo.employee1 (
        id int PRIMARY KEY,
        tSimple simple
    ) WITH CLUSTERING ORDER BY (id ASC)
```
As the parser uses regular expressions please don't delete the WITH CLUSTERING string or the programme won't work.

Step 2:

Execute in the home directory of this project the main program using the following minimal set of flags:
```
go run main.go -file=/Users/stevef/Source_Code/go/src/github.com/stevef1uk/test4/t.cql \
               -goPackageName=github.com/stevef1uk/test4 \
               -dirToGenerateIn=/Users/stevef/Source_Code/go/src/github.com/stevef1uk/test4
```
The -debug-true flag will help debug any issues

The path names will need to be adjusted to where you want the generated microservice to be where you want it created on your machine

Then to run the generated microservice run:
```
cd Users/stevef/Source_Code/go/src/github.com/stevef1uk/test4

export CASSANDRA_SERVICE_HOST=127.0.0.1

epxort PORT=5000

go run cmd/simple-server/main.go 
NOTE: the latets verion of swagger generates the entry point as cmd/simple-api-server/main.go
```
In order for this command to work the environment variable CASSANDRA_SERVICE_HOST needs to set to the host name(s) of the Cassandra cluster. 
Setting the PORT environment variable will make it easier to test 
If the cassandra database has authentication enabled then also set CASSANDRA_USERNAME and CASSANDRA_PASSWORD env vars
For DataStax Astra download the credentials from their site, unzip the file and set the following two environment variables:
ASTRA_SECURE_CONNECT_PATH - to the absolute path of the extracted zip file directory
ASTRA_PORT - the port contained in the cqlshrc file contained in the directory above
Astra uses a TLS connection and also requires the ASSANDRA_USERNAME and CASSANDRA_PASSWORD env vars set

For the example above the command to test it is:
```
curl -X GET "http://127.0.0.1:5000/v1/employee?id=1"
```
Note: To insert test data: 	
```
insert into employee1 (id, tsimple) values ( 1, { id:1,dummy:'steve',mediate:'1999-12-01T23:21:59.123Z'} ) ;

or if you have built using the post flag use curl e.g. 
curl -d '{"id": 1, "mesdummy": "steve"}' -H "Content-Type: application/json" -v -X POST http://localhost:5000/v1/employee
```
Examples of other test tables & types I have used to test this are shown in the handler folder in file create_cassandra_handler_test.go


Known issues:

The functonality of this tool to support the Cassandra MAP type is very limited. This tool can only cope with maps if they are defined as map<text,x>; this is simply becasue there seems to be no way of modelling the map type in Swagger! 
I have used the addionalProperties arroach and hard coded this to a string. For now if maps are used that are not of this trivial form the generated code handler will need to be modified manually.
Also, whilst the value of the map can be a UDT as well as simple types, the UDT in the map (or UDT chain) can't contain other UDTs at present.

I am in the process of supporting UDTs in insert clauses and there are some limitations. I haven't started work on Posts for UDT.

The ability to use the -post flag to support POST operation is limited to the following basic table types + date
```
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
```
To insert records (assuming --post flag set) used the command:
```
curl -d '{"id": 1, "testfloat": 123.45, "testtimestamp": "2018-02-17T13:01:05.000Z", "testbigint": 123456789012345, "testblob": "Some long blob stuff here!", "testbool": true, "testdouble": 123.987, "testlist": ["dummy", "something"], "testset": [6,7,8], "testmap": {"a":"alpha","b":"beta"}}' -H "Content-Type: application/json" -v -X POST http://localhost:5000/v1/demo1
```
Retrieve the data with:
```
curl -X GET "http://127.0.0.1:5000/v1/demo1?id=1"
```
Not supplying all fields in a POST will reset the values not passed to null or defaults e.g. date fields will be set to 1970-01-01

BUGS:
1. Types of VARINT don't work - I have found set<VARINT> won't work
2. For Post only tables without UDTs are supported. This is because I haven't managed to use gocql to insert UDTs.
2. For Post only a subset of Cassandra data types are supported (as I can't figure out how to get gocql to insert them :-) 


A more complex example:
```
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



insert into employee ( id, mediate, second_ts, name,  my_list, address_set  ) values (1, '2018-02-17T13:01:05.000Z', '1999-12-01T23:21:59.123Z', 'steve', [{dummy:'fred'}], {{id:1, mymap:{'a':'fred'}, citycode:'Peef',lastupdatedat:'2019-02-18T14:02:06.000Z',address_list:{{dummy:'foobar'}},events:{1,2,3} }} ) ;

curl -X GET "http://127.0.0.1:5000/v1/employee?id=1&mediate=2018-02-17T13:01:05.000Z&second_ts=1999-12-01T23:21:59.123Z"
```

Random notes:

My last approach for the parser was not to use a lex / yacc grammer but use regular expressions to identify the key data needed. This approach worked but produced code that was very hard to understand. Therefore, I have taken a differnet approach this time. I have written a simple FSM that contains states and is configured to look for regular expressions in each state and invoke functions to process the matches.
