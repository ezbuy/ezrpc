package ezrpc

import (
	"bytes"
	"errors"
	"strings"
	"time"

	"github.com/Wuvist/go-thrift/thrift"
	"github.com/nats-io/nats"
)

type OnewayRequest interface {
	Oneway() bool
}

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

	// 认为客户端 UNTIL 类请求是超长超时的请求
	timeout := 10 * time.Second
	if strings.HasPrefix(method, "UNTIL") {
		timeout = time.Hour
		method = method[5:]
	}

	var subject string
	if strings.HasPrefix(method, "Direct") {
		if c.DirectKey == "" {
			return errors.New("client DirectKey is empty")
		}
		subject = c.DirectKey + "." + c.Service + "." + method
	} else {
		subject = c.Service + "." + method

		if strings.HasPrefix(method, "On") {
			if onewayReq, ok := request.(OnewayRequest); ok && onewayReq != nil && onewayReq.Oneway() {
				subject = "On." + subject
			}
		}
	}

	// 认为客户端实现过程中 broadcast 类的请求始终 reponse == nil
	if response == nil {
		return c.Conn.Publish(subject, buf.Bytes())
	}

	msg, err := c.Conn.Request(subject, buf.Bytes(), timeout)
	if err != nil {
		println(err.Error())
		return err
	}

	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	return thrift.DecodeStruct(r, response)
}
