package nrpc

import (
	"context"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"testing"

	"github.com/stretchr/testify/assert"
)

var i int

func TestChain(t *testing.T) {
	next := func(ctx context.Context) (proto.Message, error) {
		t.Log(ctx)
		i += 10
		return nil, nil
	}

	got, err := Chain(test1Middleware, test2Middleware, test3Middleware)(next)(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, got, "reply")
	assert.Equal(t, i, 16)
}

func test1Middleware(handler Handler) Handler {
	return func(ctx context.Context) (reply proto.Message, err error) {
		fmt.Println("test1 before")
		i++
		reply, err = handler(ctx)
		fmt.Println("test1 after")
		return
	}
}

func test2Middleware(handler Handler) Handler {
	return func(ctx context.Context) (reply proto.Message, err error) {
		fmt.Println("test2 before")
		i += 2
		reply, err = handler(ctx)
		fmt.Println("test2 after")
		return
	}
}

func test3Middleware(handler Handler) Handler {
	return func(ctx context.Context) (reply proto.Message, err error) {
		fmt.Println("test3 before")
		i += 3
		reply, err = handler(ctx)
		fmt.Println("test3 after")
		return
	}
}
