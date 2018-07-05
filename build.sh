#!/bin/bash
go get -u github.com/hunkeelin/SuperCA/server 
go get -u github.com/hunkeelin/SuperCA/utils 
go get -u github.com/hunkeelin/mtls/server 
go get -u github.com/hunkeelin/klinenv 
go get -u github.com/hunkeelin/pki 
go build -o SuperCA *go
cp SuperCA $HOME/rpmbuild/SOURCES/superca-1.0/usr/bin
ssh rpmbuild "cd $HOME/rpmbuild/SOURCES; tar -zvcf superca-1.0.tar.gz superca-1.0"
ssh rpmbuild "rpmbuild -ba $HOME/rpmbuild/SPECS/SuperCA.spec"
cp $HOME/rpmbuild/RPMS/x86_64/superca-1.0-1.x86_64.rpm /tmp
