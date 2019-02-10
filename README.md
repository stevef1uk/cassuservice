# cassuservice
Repo that contains code that will when completed autogenerate a Golang microservice from a Cassandra DDL for a table.

I originally wrote equivalent code whilst at my last employer as an exercise to learning Golang. I am no longer a professional programmer, but still enjoy the challenges of coding. I did gain permission to open source my first working version of the previous  work, but was unable to get the code officially listed before I decided to move on. Therefore, to ensure that there is no risk of my former employer causing issues for me or anyone using this library I have rewritten it in this repo. Given I am now working again this work will proceed quite slowly.

The architecture approach is simple:

1. Parse Cassandra DDL to create a swagger file
2. Use go-swagger to generate the RESTful server
3. Parse Cassandra DDL to create the Cassandra handler and wite it into the RESTful server.

My last approach for the parser was not to use a lex / yacc grammer but use regular expressions to identify the key data needed. This approach worked but produced code that was very hard to understand. Therefore, I have taken a differnet approach this time. I have written a simple FSM that contains states and is configured to look for regular exprerssions in each state and invoke functions to process the matches.

Latest Status: FSM to parse Cassandra CQL written so (1) done. Also used the output of the parser to create the swagger data is done. This time
I have minimised the use of templates and used string concatenation.
As of 9th Feb the work to generate the Cassandra handler is working in that the code generated compiles and returns the appropriate data from
Cassandra. This is very early days, but promising as the only test example used UDTs that themselves use UDTs so recursion going on.


Known issues:

The functonality of this tool to support the Cassandra MAP type is very limited. This tool can only cope with maps if they are defined as map<text,text>; this is simply becasue there seems to be no way of modelling the map type in Swagger! I have used the addionalProperties arroach and hard coded this to a string. For now if maps are used that are not of this trivial form the generated code handler will need to be modified manually.


The first test for the genration uses the following schema, so not trivial!

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
 )



