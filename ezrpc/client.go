package ezrpc

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"github.com/Wuvist/go-thrift/thrift"
	"github.com/nats-io/nats"
)

type Client struct {
	Conn      *nats.Conn
	Service   string
	DirectKey string
}

func NewClient(service string, conn *nats.Conn) *Client {
	return &Client{
		Service: service,
		Conn:    conn,
	}
}

func (c *Client) Call(method string, request interface{}, response interface{}) error {
	buf := &bytes.Buffer{}
	w := thrift.NewCompactProtocolWriter(buf)
	thrift.EncodeStruct(w, request)

	var subject string
	if strings.HasPrefix(method, "Direct") {
		if c.DirectKey == "" {
			return errors.New("client DirectKey is empty")
		}
		subject = c.DirectKey + "." + c.Service + "." + method
	} else {
		subject = c.Service + "." + method
	}

	if response == nil {
		return c.Conn.Publish(subject, buf.Bytes())
	}

	msg, err := c.Conn.Request(subject, buf.Bytes(), 10*time.Second)
	if err != nil {
		return err
	}
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	return thrift.DecodeStruct(r, response)
}
