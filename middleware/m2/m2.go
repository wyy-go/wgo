package m2

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/wyy-go/wgo/middleware"
	"log"
)

func M2() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {

		return func(ctx context.Context, msg *nats.Msg) (rsp interface{}, err error) {
			log.Println("==========m2 start=========")
			rsp, err = handler(ctx, msg)
			log.Println("==========m2 start=========")
			return
		}
	}
}
