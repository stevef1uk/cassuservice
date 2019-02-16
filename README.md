# cassuservice
Repo that contains Golang code that can autogenerate a Golang microservice from a Cassandra DDL for a table.

I originally wrote equivalent code whilst at my last employer as an exercise to learning Golang. I am no longer a professional programmer, but still enjoy the challenges of coding. I did gain permission to open source my first working version of the previous  work, but was unable to get the code officially listed before I decided to move on. Therefore, I have rewritten it in this repo and done it better the second time I believe. 

The architecture approach is simple:

1. Parse Cassandra DDL to create a swagger file
2. Use go-swagger to generate the RESTful server, see: https://github.com/go-swagger/go-swagger
3. Parse Cassandra DDL to create the Cassandra handler and wire it into the RESTful server.

I have updated the main.go file to perform the above steps in one.

Step 1: Export the Cassandra DDL into a file e.g. t.cql. This file needs to only contain the types (if any) and table definition for a single table. The best way to do this is from cqlsh and use the describe table command.

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

As the parser uses regular expressions please don't delete the WITH CLUSTERING string or the programme won't work.

Step 2:

Execute the main program using the following minimal set of flags:

go run main.go -file=/Users/stevef/Source_Code/go/src/github.com/stevef1uk/test4/t.cql \
               -goPackageName=github.com/stevef1uk/test4 \
               -dirToGenerateIn=/Users/stevef/Source_Code/go/src/github.com/stevef1uk/test4

The -debug-true flag will help debug any issues

The path names will need to be adjusted to where you want the generated microservice to be where you want it created on you machine

Then to run the generated microservice run:

cd Users/stevef/Source_Code/go/src/github.com/stevef1uk/test4

export CASSANDRA_SERVICE_HOST=127.0.0.1

epxort PORT=5000

go run cmd/simple-server/main.go 

In order for this command to work the environment variable CASSANDRA_SERVICE_HOST needs to set to the host name(s) of the Cassandra cluster. 
Setting the PORT environmnet variable will make it easier to test 

For the example above the command to test it is:

curl -X GET "http://127.0.0.1:5000/v1/employee?id=1"

Note: To insert test data: 	
insert into employee1 (id, tsimple) values ( 1, { id:1,dummy:'steve',mediate:'1999-12-01T23:21:59.123Z'} ) ;

Examples of other test tables & types I have used to test this are shown in the handler folder in file create_cassandra_handler_test.go

Random notes:

My last approach for the parser was not to use a lex / yacc grammer but use regular expressions to identify the key data needed. This approach worked but produced code that was very hard to understand. Therefore, I have taken a differnet approach this time. I have written a simple FSM that contains states and is configured to look for regular expressions in each state and invoke functions to process the matches.

Known issues:

The functonality of this tool to support the Cassandra MAP type is very limited. This tool can only cope with maps if they are defined as map<text,text>; this is simply becasue there seems to be no way of modelling the map type in Swagger! I have used the addionalProperties arroach and hard coded this to a string. For now if maps are used that are not of this trivial form the generated code handler will need to be modified manually.

A more complex example:

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
