package middleware

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/wyy-go/wgo/nrpc"
)

func Wrap(h nrpc.Handler, m ...Middleware) nats.MsgHandler {
	chain := Chain(m...)
	return func(msg *nats.Msg) {
		next := func(ctx context.Context, msg *nats.Msg) (interface{}, error) {
			h.SetContext(ctx)
			h.Handler(msg)
			return nil, nil
		}
		next = chain(next)
		next(context.TODO(), msg)
	}
}
