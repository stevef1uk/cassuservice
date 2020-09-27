#!/bin/bash
# Simple script to run on a GitPod terminal to prove the software is working
mkdir /workspace/test4
cd /workspace/test4
cat > t.cql <<EOF
CREATE TABLE demo.verysimple (
    id int PRIMARY KEY,
    message text
);
EOF
cd /workspace/cassuservice
go run main.go -file=/workspace/test4/t.cql -goPackageName=github.com/stevef1uk/test4 -dirToGenerateIn=/workspace/go/src/github.com/stevef1uk/test4
echo "Code generated in /workspace/go/src/github.com/stevef1uk/test4 "
echo "cd /workspace/go/src/github.com/stevef1uk/test4 & run go init followed by go run ... "

