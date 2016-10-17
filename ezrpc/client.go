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
	cfg       *Config
	Conn      *nats.Conn
	Service   string
	DirectKey string
}

func NewClient(service string, conn *nats.Conn) *Client {
	return NewClientEx(&Config{}, service, conn)
}

type Config struct {
	Timeout time.Duration
	Reties  int
}

func (cfg *Config) init() {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 10 * time.Second
	}
}

func NewFastRetryClient(service string, conn *nats.Conn) *Client {
	return NewClientEx(&Config{
		Timeout: 2 * time.Second,
		Reties:  3,
	}, service, conn)
}

func NewClientEx(cfg *Config, service string, conn *nats.Conn) *Client {
	cfg.init()
	return &Client{
		Service: service,
		Conn:    conn,
		cfg:     cfg,
	}

}

func NewClientTimeout(service string, timeout time.Duration, conn *nats.Conn) *Client {
	return NewClientEx(&Config{
		Timeout: timeout,
	}, service, conn)
}

func (c *Client) Call(method string, request interface{}, response interface{}) error {
	buf := &bytes.Buffer{}
	w := thrift.NewCompactProtocolWriter(buf)
	thrift.EncodeStruct(w, request)

	// 认为客户端 UNTIL 类请求是超长超时的请求
	timeout := c.cfg.Timeout
	if strings.HasPrefix(method, "UNTIL") {
		if c.cfg.Timeout < time.Hour {
			timeout = time.Hour
		}

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

	retryTime := c.cfg.Reties
retry:
	msg, err := c.Conn.Request(subject, buf.Bytes(), timeout)
	if err == nats.ErrTimeout {
		if retryTime > 0 {
			time.Sleep(100 * time.Millisecond)
			retryTime--
			goto retry
		}
	}
	if err != nil {
		return err
	}

	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	return thrift.DecodeStruct(r, response)
}
