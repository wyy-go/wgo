package auth

import (
	"github.com/wyy-go/wgo/nrpc"
)

func Authorize() nrpc.Middleware {
	return func(handler nrpc.Handler) nrpc.Handler {
		return handler
	}
}
