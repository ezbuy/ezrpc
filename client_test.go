package ezrpc

import (
	"bytes"
	"encoding/json"
	"erproduct"
	"fmt"
	"testing"
	"time"

	"github.com/samuel/go-thrift/thrift"

	"github.com/nats-io/nats"
)

func onReply(msg *nats.Msg) {
	var err error
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	p := &erproduct.TProduct{}

	err = thrift.DecodeStruct(r, p)

	if err != nil {
		println(err.Error())
		return
	}
	println(p.Location)
}

func onFoo(msg *nats.Msg) {
	println("miao")
	var err error
	r := thrift.NewCompactProtocolReader(bytes.NewReader(msg.Data))

	p := &erproduct.TProduct{}

	err = thrift.DecodeStruct(r, p)

	if err != nil {
		println(err.Error())
		return
	}
	println(p.Location)
	err = thrift.DecodeStruct(r, p)

	if err != nil {
		println(err.Error())
		return
	}
	println(p.Location)

	obj := &erproduct.SearchFilter{}
	err = thrift.DecodeStruct(r, obj)

	if err != nil {
		println(err.Error())
		return
	}
	println(obj.Name)
	err = thrift.DecodeStruct(r, obj)

	if err != nil {
		println(err.Error())
		return
	}
	println(obj.Name)

	buf := &bytes.Buffer{}
	p.Location = "xiamen"
	w := thrift.NewCompactProtocolWriter(buf)
	thrift.EncodeStruct(w, p)
	nc.Publish(msg.Reply, buf.Bytes())
}

type productServiceImplementation int

func (s *productServiceImplementation) GetProductDetail(productUrl string, purchaseSource string) (*erproduct.TProduct, error) {
	result := new(erproduct.TProduct)
	result.ProductUrl = productUrl + purchaseSource
	return result, nil
}

func TestMain(t *testing.T) {
	server := new(productServiceImplementation)
	NewServer(server)
	time.Sleep(10 * time.Millisecond)

	client := NewClient("Product")
	scr := erproduct.ProductClient{Client: client}
	product, err := scr.GetProductDetail("productUrl", "surf")
	if err != nil {
		t.Error(err)
	}

	if product.ProductUrl != "productUrlsurf" {
		t.Error("server response error")
	}

}

func MainTEst() {

	var nc *nats.Conn
	client := NewClient("Product")
	scr := erproduct.ProductClient{Client: client}
	product, _ := scr.GetProductDetail("productUrl", "surf")
	println(product.ProductUrl)
	nc, _ = nats.Connect(nats.DefaultURL)
	// nc.Request("", data, timeout)
	buf := &bytes.Buffer{}
	obj := &erproduct.SearchFilter{}
	obj.Name = "bingo"

	p := &erproduct.TProduct{}
	p.Location = "singapore"
	w := thrift.NewCompactProtocolWriter(buf)
	thrift.EncodeStruct(w, obj)

	data, _ := json.Marshal(obj)
	fmt.Printf("%x\n", buf.Bytes())
	println("===")
	fmt.Printf("%x\n", data)
	println("json")
	thrift.EncodeStruct(w, p)
	p.Location = "shanghai"

	thrift.EncodeStruct(w, p)

	obj.Name = "blah blah blah blah "
	thrift.EncodeStruct(w, obj)

	// nc.Publish("foo", buf.Bytes())
	nc.Subscribe("foo", onFoo)
	// nc.QueueSubscribe("foo", "miao", onFoo)

	// nc.QueueSubscribe("foo", "miao", func(msg *nats.Msg) {
	// 	println(string(msg.Data))
	// })

	// Simple Async Subscriber
	// nc.Subscribe("foo", func(m *nats.Msg) {
	// 	fmt.Printf("Received a message: %s\n", string(m.Data))
	// })
	// Simple Publisher

	time.Sleep(10 * time.Millisecond)
	// nc.Publish("foo", buf.Bytes())
	msg, err := nc.Request("foo", buf.Bytes(), 100*time.Millisecond)
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Millisecond)
	onReply(msg)
	time.Sleep(10 * time.Millisecond)
	// Simple Sync Subscriber
	// sub, err := nc.SubscribeSync("foo")
	// m, err := sub.NextMsg(100)
	// println(m, err)

	// // Channel Subscriber
	// ch := make(chan *nats.Msg, 64)
	// sub, err = nc.ChanSubscribe("foo", ch)
	// var msg *nats.Msg
	// msg <- ch

	// // Unsubscribe
	// sub.Unsubscribe()

	// // Requests

	// // Replies
	// nc.Subscribe("help", func(m *nats.Msg) {
	// 	nc.Publish(m.Reply, []byte("I can help!"))
	// })

	// Close connection
	nc.Close()

}
