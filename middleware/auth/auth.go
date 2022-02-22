package auth

import (
	"github.com/wyy-go/wgo/middleware"
)

func Authorize() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return handler
	}
}
