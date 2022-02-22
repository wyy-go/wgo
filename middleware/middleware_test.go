package middleware

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"testing"

	"github.com/stretchr/testify/assert"
)

var i int

func TestChain(t *testing.T) {
	next := func(ctx context.Context, msg *nats.Msg) (interface{}, error) {
		t.Log(msg.Subject)
		i += 10
		return "reply", nil
	}

	got, err := Chain(test1Middleware, test2Middleware, test3Middleware)(next)(context.Background(), &nats.Msg{})
	assert.Nil(t, err)
	assert.Equal(t, got, "reply")
	assert.Equal(t, i, 16)
}

func test1Middleware(handler Handler) Handler {
	return func(ctx context.Context, msg *nats.Msg) (reply interface{}, err error) {
		fmt.Println("test1 before")
		i++
		reply, err = handler(ctx, msg)
		fmt.Println("test1 after")
		return
	}
}

func test2Middleware(handler Handler) Handler {
	return func(ctx context.Context, msg *nats.Msg) (reply interface{}, err error) {
		fmt.Println("test2 before")
		i += 2
		reply, err = handler(ctx, msg)
		fmt.Println("test2 after")
		return
	}
}

func test3Middleware(handler Handler) Handler {
	return func(ctx context.Context, msg *nats.Msg) (reply interface{}, err error) {
		fmt.Println("test3 before")
		i += 3
		reply, err = handler(ctx, msg)
		fmt.Println("test3 after")
		return
	}
}
