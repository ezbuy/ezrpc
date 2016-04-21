package product

import (
	"bytes"

	"github.com/nats-io/nats"
	"github.com/samuel/go-thrift/thrift"
	"github.com/ezbuy/statsd"
)

type ThriftNatsProductServer struct {
	Server    *ProductServer
	Conn      *nats.Conn
	DirectKey string
}

func (s *ThriftNatsProductServer) onMsg(msg *nats.Msg) {
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	switch msg.Subject {
	case "Product.GetProductDetail":
		t1 := statsd.Now()

		p := &ProductGetProductDetailRequest{}
		res := &ProductGetProductDetailResponse{}
		err := thrift.DecodeStruct(r, p)
		if err != nil {
			println(err)
		}
		err = s.Server.GetProductDetail(p, res)
		if err != nil {
			println(err)
		}

		buf := &bytes.Buffer{}
		w := thrift.NewCompactProtocolWriter(buf)
		thrift.EncodeStruct(w, res)
		s.Conn.Publish(msg.Reply, buf.Bytes())

		t2 := statsd.Now()
		statsd.Timing("Product.GetProductDetail.timing", t1, t2)
	case "Product.Ping":
		p := &ProductPingRequest{}
		err := thrift.DecodeStruct(r, p)
		if err != nil {
			println(err)
		}
		err = s.Server.Ping(p)
		if err != nil {
			println(err)
		}

		statsd.Incr("Product.Ping.count")
	}
}

func (s *ThriftNatsProductServer) onDirect(msg *nats.Msg) {
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))
	keylen := len(s.DirectKey)
	switch msg.Subject[keylen+1:] {
	case "Product.DirectGetProductDetail":
		t1 := statsd.Now()

		p := &ProductDirectGetProductDetailRequest{}
		res := &ProductDirectGetProductDetailResponse{}
		err := thrift.DecodeStruct(r, p)
		if err != nil {
			println(err)
		}
		err = s.Server.DirectGetProductDetail(p, res)
		if err != nil {
			println(err)
		}

		buf := &bytes.Buffer{}
		w := thrift.NewCompactProtocolWriter(buf)
		thrift.EncodeStruct(w, res)
		s.Conn.Publish(msg.Reply, buf.Bytes())

		t2 := statsd.Now()
		statsd.Timing("Product.DirectGetProductDetail.timing", t1, t2)
	case "Product.DirectOnCacheEvict":
		p := &ProductDirectOnCacheEvictRequest{}
		err := thrift.DecodeStruct(r, p)
		if err != nil {
			println(err)
		}
		err = s.Server.DirectOnCacheEvict(p)
		if err != nil {
			println(err)
		}

		statsd.Incr("Product.DirectOnCacheEvict.count")
	}
}

func (s *ThriftNatsProductServer) onBroadcast(msg *nats.Msg) {
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	switch msg.Subject {
	case "On.Product.OnCacheEvict":
		p := &ProductOnCacheEvictRequest{}
		err := thrift.DecodeStruct(r, p)
		if err != nil {
			println(err)
		}
		err = s.Server.OnCacheEvict(p)
		if err != nil {
			println(err)
		}

		statsd.Incr("On.Product.OnCacheEvict.count")
	case "On.Product.OnExchangeUpdate":
		p := &ProductOnExchangeUpdateRequest{}
		err := thrift.DecodeStruct(r, p)
		if err != nil {
			println(err)
		}
		err = s.Server.OnExchangeUpdate(p)
		if err != nil {
			println(err)
		}

		statsd.Incr("On.Product.OnExchangeUpdate.count")
	}
}

func (s *ThriftNatsProductServer) SetDirectKey(key string) {
	s.DirectKey = key
	s.Conn.Subscribe(key + ".Product.*", s.onDirect)
}

func NewProductServer(impl Product, conn *nats.Conn) *ThriftNatsProductServer {
	s := &ProductServer{Implementation: impl}

	server := &ThriftNatsProductServer{
		Server: s,
		Conn:   conn,
	}

	// all broadcast messages should be under namespace 'On'
	server.Conn.Subscribe("On.Product.*", server.onBroadcast)
	server.Conn.QueueSubscribe("Product.*", "ezrpc", server.onMsg)
	return server
}