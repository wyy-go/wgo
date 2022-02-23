package m1

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/wyy-go/wgo/nrpc"

	"log"
)

func M1() nrpc.Middleware {
	return func(handler nrpc.Handler) nrpc.Handler {
		return func(ctx context.Context) (rsp proto.Message, err error) {
			log.Println("==========m1 start=========")
			rsp, err = handler(ctx)
			if err != nil {
				return nil, err
			}
			log.Println("==========m1 start=========")
			return
		}
	}
}
