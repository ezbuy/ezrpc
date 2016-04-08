SHELL := /bin/bash

init:
	# go client
	go get github.com/nats-io/nats
	# server
	go get github.com/nats-io/gnatsd
	go get github.com/jteeuwen/go-bindata/...
	
buildtpl:
	rm tmpl/bindata.go
	go-bindata -o tmpl/bindata.go -ignore bindata.go -pkg tmpl tmpl/...
	
indev: buildtpl
	go build -o exe
	./exe gen -l csharp -i example/Product.thrift -o .
