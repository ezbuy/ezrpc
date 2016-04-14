# ezrpc

ezrpc is a `micro service` framework for server side rpc communication.

It's based on [nats](http://nats.io/) and [thrift](https://github.com/samuel/go-thrift), using code-gen approach, supporting Go & .net(C#).

# Service Definition

```thrift
service Category {
	list<string> GetIDs(1:i32 offset, 2:i32 limit),
}
```

# Usage

1 Generate language specified source files by `thrift` IDL

	* C#

	thrift --gen csharp -o ./sample/ ./sample/HelloWorld.thrift

	* Go

	generator ./sample/HelloWorld.thrift` ./sample/

2 Genrate source files which will be used for subscribing NATS messages

	* C#

	./ezrpc gen -l csharp -i ./sample/HelloWorld.thrift -o ./sample

	* Go

	./ezrpc gen -l go -i ./sample/HelloWorld.thrift -o ./sample
