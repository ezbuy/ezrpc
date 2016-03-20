package main

import (
	"time"

	"github.com/nats-io/nats"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)
	nc.Request("", data, timeout)
	nc.Publish("foo", []byte("Hello World1"))
	nc.QueueSubscribe("foo", "miao", func(msg *nats.Msg) {
		println(string(msg.Data))
	})

	nc.QueueSubscribe("foo", "miao", func(msg *nats.Msg) {
		println(string(msg.Data))
	})

	// Simple Async Subscriber
	// nc.Subscribe("foo", func(m *nats.Msg) {
	// 	fmt.Printf("Received a message: %s\n", string(m.Data))
	// })
	// Simple Publisher
	nc.Publish("foo", []byte("Hello World"))

	time.Sleep(2 * time.Second)
	// Simple Sync Subscriber
	sub, err := nc.SubscribeSync("foo")
	m, err := sub.NextMsg(100)
	println(m, err)

	// // Channel Subscriber
	// ch := make(chan *nats.Msg, 64)
	// sub, err = nc.ChanSubscribe("foo", ch)
	// var msg *nats.Msg
	// msg <- ch

	// // Unsubscribe
	// sub.Unsubscribe()

	// // Requests
	// msg, err = nc.Request("help", []byte("help me"), 10*time.Millisecond)

	// // Replies
	// nc.Subscribe("help", func(m *nats.Msg) {
	// 	nc.Publish(m.Reply, []byte("I can help!"))
	// })

	// Close connection
	nc.Close()

}
