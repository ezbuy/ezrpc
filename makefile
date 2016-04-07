SHELL := /bin/bash

init:
	go get github.com/jteeuwen/go-bindata/...
	
buildtpl:
	rm tmpl/bindata.go
	go-bindata -o tmpl/bindata.go -ignore bindata.go -pkg tmpl tmpl/...