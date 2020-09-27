#!/bin/bash
# Simple script to run on a GitPod terminal to prove the software is working. Important assumption that the following table exists in your Astra Cassandra instance ;-)
mkdir /workspace/test4
cd /workspace/test4
cat > t.cql <<EOF
CREATE TABLE demo.verysimple (
    id int PRIMARY KEY,
    message text,
    WITH CLUSTERING ORDER BY (id ASC)
);
EOF
cd /workspace/cassuservice
go run main.go -file=/workspace/test4/t.cql -goPackageName=github.com/stevef1uk/test4 -dirToGenerateIn=/workspace/go/src/github.com/stevef1uk/test4 --post=true -consistency=gocql.LocalQuorum
echo "Code generated in /workspace/go/src/github.com/stevef1uk/test4 "
echo "You will need to donload the Astra secure connect bundle and upload them to this GitPod and set the appropriate values in env.sh before this script will work"
echo "In another terminal window type: curl -X GET 'http://127.0.0.1:5000/v1/verysimple?id=1'"
echo 'To store a record use: curl -d ' "'" '{"id": 1, "message": "steve"}' "'"  ' -H "Content-Type: application/json" -v -X POST  "http://127.0.0.1:5000/v1/verysimple?id=1"'
cd /workspace/go/src/github.com/stevef1uk/test4 
go mod init github.com/stevef1uk/test4
. /workspace/cassuservice/env.sh
go run cmd/simple-api-server/main.go 
