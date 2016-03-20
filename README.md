# ezrpc

ezrpc is a `micro service` framework for server side rpc communication.

It's based on [nats](http://nats.io/) and [thrift](https://github.com/samuel/go-thrift), using code-gen approach, supporting Go & .net(C#).

# Service Definition

```thrift
service Category {
	list<string> GetIDs(1:i32 offset, 2:i32 limit),
}
```
