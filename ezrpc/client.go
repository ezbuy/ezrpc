package ezrpc

import (
	"bytes"
	"time"

	"github.com/nats-io/nats"
	"github.com/samuel/go-thrift/thrift"
)

type Client struct {
	Conn   *nats.Conn
	Serice string
}

func NewClient(service string, conn *nats.Conn) *Client {
	return &Client{
		Serice: service,
		Conn:   conn,
	}
}

func (c *Client) Call(method string, request interface{}, response interface{}) error {
	buf := &bytes.Buffer{}
	w := thrift.NewCompactProtocolWriter(buf)
	thrift.EncodeStruct(w, request)

	if response == nil {
		return c.Conn.Publish(c.Serice+"."+method, buf.Bytes())
	}

	msg, err := c.Conn.Request(c.Serice+"."+method, buf.Bytes(), 10*time.Second)
	if err != nil {
		println(err.Error())
	}
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	return thrift.DecodeStruct(r, response)
}
