package ezrpc

import (
	"bytes"
	"time"

	"github.com/nats-io/nats"
	"github.com/samuel/go-thrift/thrift"
)

type Client struct {
	Serice string
}

var nc *nats.Conn

func init() {
	nc, _ = nats.Connect(nats.DefaultURL)
}

func NewClient(service string) *Client {
	return &Client{
		Serice: service,
	}
}

func (c *Client) Call(method string, request interface{}, response interface{}) error {
	buf := &bytes.Buffer{}
	w := thrift.NewCompactProtocolWriter(buf)
	thrift.EncodeStruct(w, request)

	msg, err := nc.Request(c.Serice+"."+method, buf.Bytes(), 10*time.Second)
	if err != nil {
		println(err.Error())
	}
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	return thrift.DecodeStruct(r, response)
}
