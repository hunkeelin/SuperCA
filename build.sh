#!/bin/bash
go get -u github.com/hunkeelin/mtls/server 
go get -u github.com/hunkeelin/klinenv 
go get -u github.com/hunkeelin/pki 
go get -u github.com/hunkeelin/SuperCA/utils 
go get -u github.com/hunkeelin/SuperCA/server 
go build -o SuperCA *go
mv SuperCA superca_1.0-1/usr/bin/
