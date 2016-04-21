package ezrpc

import (
	"bytes"
	"strings"
	"time"

	"github.com/nats-io/nats"
	"github.com/samuel/go-thrift/thrift"
)

type Client struct {
	Conn      *nats.Conn
	Serice    string
	DirectKey string
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

	var subject string
	if strings.HasPrefix(method, "Direct") {
		subject = c.DirectKey + "." + c.Serice + "." + method
	} else {
		subject = c.Serice + "." + method
	}

	if response == nil {
		return c.Conn.Publish(subject, buf.Bytes())
	}

	msg, err := c.Conn.Request(subject, buf.Bytes(), 10*time.Second)
	if err != nil {
		println(err.Error())
	}
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	return thrift.DecodeStruct(r, response)
}
