package ezrpc

import (
	"bytes"
	"erproduct"

	"github.com/nats-io/nats"
	"github.com/samuel/go-thrift/thrift"
)

func (s *ThriftNatsServer) onMSG(msg *nats.Msg) {
	println("onMSG")
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
	nc.Publish(msg.Reply, buf.Bytes())
}

type ThriftNatsServer struct {
	Server *erproduct.ProductServer
}

func NewServer(impl erproduct.Product) {
	s := &erproduct.ProductServer{Implementation: impl}

	server := &ThriftNatsServer{
		Server: s,
	}
	nc.Subscribe("Product.GetProductDetail", server.onMSG)
}
