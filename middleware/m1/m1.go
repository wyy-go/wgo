package m1

import (
	"context"
	"errors"
	"github.com/nats-io/nats.go"
	"github.com/wyy-go/wgo/middleware"
	"log"
)

func M1() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {

		return func(ctx context.Context, msg *nats.Msg) (rsp interface{}, err error) {
			log.Println("==========m1 start=========")
			err = errors.New("===============")
			return
			rsp, err = handler(ctx, msg)
			log.Println("==========m1 start=========")
			return
		}
	}
}
