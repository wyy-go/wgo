package m2

import (
	"context"
	"github.com/gogo/protobuf/proto"
	"github.com/wyy-go/wgo/nrpc"
	"log"
)

func M2() nrpc.Middleware {
	return func(handler nrpc.Handler) nrpc.Handler {

		return func(ctx context.Context) (rsp proto.Message, err error) {
			log.Println("==========m2 start=========")
			rsp, err = handler(ctx)
			if err != nil {
				return nil, err
			}
			log.Println("==========m2 start=========")
			return
		}
	}
}
