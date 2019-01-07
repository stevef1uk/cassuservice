# cassuservice
Repo that contains code that will when completed autogenerate a Golang microservice from a Cassandra DDL for a table.

I originally wrote equivalent code whilst at my last employer as an exercise to learning Golang. I am no longer a professional programmer, but still enjoy the challenges of coding. I did gain permission to open source my first working version of the previous  work, but was unable to get the code officially listed before I decided to move on. Therefore, to ensure that there is no risk of my former employer causing issues for me or anyone using this library I have rewritten it in this repo. Given I am now working again this work will proceed quite slowly.

The architecture approach is simple:

1. Parse Cassandra DDL to create a swagger file
2. Use go-swagger to generate the RESTful server
3. Parse Cassandra DDL to create the Cassandra handler and wite it into the RESTful server.

My last approach for the parser was not to use a lex / yacc grammer but use regular expressions to identify the key data needed. This approach worked but produced code that was very hard to understand. Therefore, I have taken a differnet approach this time. I have written a simple FSM that contains states and is configured to look for regular exprerssions in each state and invoke functions to process the matches.

Latest Status: FSM to parse Cassandra CQL written so (1) done. Also used the output of the parser to create the swagger data so (2) possible. Last time I used Go templates for the swagger creation, this time I used simple string concatenation.






