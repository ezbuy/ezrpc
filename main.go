package main

import "github.com/ezbuy/ezrpc/cmd"

//go:generate go get github.com/jteeuwen/go-bindata/...
//go:generate go-bindata -o ./tmpl/bindata.go -ignore bindata.go -pkg tmpl tmpl/golang
func main() {
	cmd.Execute()
}
