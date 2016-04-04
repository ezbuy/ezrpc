package ezrpc

import (
	"bytes"
	"erproduct"

	"github.com/nats-io/nats"
	"github.com/samuel/go-thrift/thrift"
)

func (s *ThriftNatsServer) onMSG(msg *nats.Msg) {
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	p := &erproduct.ProductGetProductDetailRequest{}
	res := &erproduct.ProductGetProductDetailResponse{}
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
}

type ThriftNatsServer struct {
	Server *erproduct.ProductServer
	Conn   *nats.Conn
}

func NewServer(impl erproduct.Product, conn *nats.Conn) {
	s := &erproduct.ProductServer{Implementation: impl}

	server := &ThriftNatsServer{
		Server: s,
		Conn:   conn,
	}
	server.Conn.Subscribe("Product.GetProductDetail", server.onMSG)
}
