package middlewares

import (
	"fmt"
	"time"

	"github.com/ezbuy/ezrpc/ezrpc"
	"github.com/nats-io/nats"
)

func RequestTimer() ezrpc.MsgMiddleware {

	return func(h nats.MsgHandler) nats.MsgHandler {
		return func(msg *nats.Msg) {
			start := time.Now()

			h(msg)

			fmt.Printf("REQ %s : %s\n", msg.Subject, time.Now().Sub(start))
		}
	}
}
