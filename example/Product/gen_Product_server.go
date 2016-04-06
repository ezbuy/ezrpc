package product

import (
	"bytes"

	"github.com/nats-io/nats"
	"github.com/samuel/go-thrift/thrift"
)

func (s *ThriftNatsProductServer) onMsg(msg *nats.Msg) {
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	switch msg.Subject {
	case "Product.GetProductDetail":
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
	}
}

type ThriftNatsProductServer struct {
	Server *ProductServer
	Conn   *nats.Conn
}

func NewProductServer(impl Product, conn *nats.Conn) {
	s := &ProductServer{Implementation: impl}

	server := &ThriftNatsProductServer{
		Server: s,
		Conn:   conn,
	}
	server.Conn.QueueSubscribe("Product.*", "ezrpc", server.onMsg)
}
