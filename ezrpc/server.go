package ezrpc

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats"
)

type MsgMiddleware func(h nats.MsgHandler) nats.MsgHandler

type Daemon struct {
	mutex sync.RWMutex
	msgWg sync.WaitGroup

	timeout time.Duration

	conn          *nats.Conn
	subscriptions []*nats.Subscription
	middelwares   []MsgMiddleware

	exit chan bool
}

func NewDaemon(opts nats.Options, middlewares ...MsgMiddleware) (*Daemon, error) {
	conn, err := opts.Connect()
	if err != nil {
		return nil, err
	}

	return NewDaemonWithConn(conn, middlewares...), nil
}

func NewDaemonWithConn(conn *nats.Conn, middlewares ...MsgMiddleware) *Daemon {
	server := &Daemon{
		timeout:       30 * time.Second,
		conn:          conn,
		subscriptions: make([]*nats.Subscription, 0),
		middelwares:   middlewares,
		exit:          make(chan bool, 1),
	}

	return server
}

func (this *Daemon) Use(middlewares ...MsgMiddleware) {
	this.mutex.Lock()
	this.middelwares = append(this.middelwares, middlewares...)
	this.mutex.Unlock()
}

func (this *Daemon) Subscribe(subject string, h nats.MsgHandler) error {
	h = this.buildMsgHandler(h)

	sub, err := this.conn.Subscribe(subject, func(msg *nats.Msg) {
		go h(msg)
	})

	if err != nil {
		return err
	}

	this.addSubscription(sub)

	return nil
}

func (this *Daemon) QueueSubscribe(subject, queue string, h nats.MsgHandler) error {
	h = this.buildMsgHandler(h)

	sub, err := this.conn.QueueSubscribe(subject, queue, func(msg *nats.Msg) {
		go h(msg)
	})

	if err != nil {
		return err
	}

	this.addSubscription(sub)

	return nil
}

func (this *Daemon) Run() {
	fmt.Fprintln(os.Stdout, "nats msg server running")
	go this.handleSignal()

	<-this.exit

	this.conn.Close()

	fmt.Fprintln(os.Stdout, "nats msg server stopped")
}

func (this *Daemon) Stop() {
	this.unsubscribeAll()

	go this.waitForHandlers()
	go this.waitForTimeout()

}

func (this *Daemon) waitForHandlers() {
	this.msgWg.Wait()

	fmt.Fprintln(os.Stdout, "nats msg server: all handlers finished")
	this.exit <- true
}

func (this *Daemon) waitForTimeout() {
	timer := time.NewTimer(this.timeout)
	<-timer.C

	fmt.Fprintln(os.Stderr, "nats msg server: timeout before all handlers finish")
	this.exit <- true
}

func (this *Daemon) handleSignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM)

	<-c
	this.Stop()
}

func (this *Daemon) unsubscribeAll() {
	this.mutex.Lock()

	subs := this.subscriptions
	this.subscriptions = []*nats.Subscription{}

	this.mutex.Unlock()

	for _, one := range subs {
		one.Unsubscribe()
	}
}

func (this *Daemon) addSubscription(sub *nats.Subscription) {
	fmt.Printf("SUB %q\n", sub.Subject)

	this.mutex.Lock()

	this.subscriptions = append(this.subscriptions, sub)

	this.mutex.Unlock()
}

func (this *Daemon) buildMsgHandler(h nats.MsgHandler) nats.MsgHandler {
	handler := func(msg *nats.Msg) {
		this.msgWg.Add(1)
		defer this.msgWg.Done()

		h(msg)
	}

	this.mutex.RLock()
	mws := this.middelwares
	this.mutex.RUnlock()

	if len(mws) == 0 {
		return handler
	}

	for i := len(mws) - 1; i >= 0; i -= 1 {
		handler = mws[i](handler)
	}

	return handler
}
