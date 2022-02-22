package main

import (
	"context"
	"github.com/wyy-go/wgo"
	"log"
	// This is the package containing the generated *.pb.go and *.nrpc.go
	// files.
	"github.com/wyy-go/wgo/examples/helloworld/helloworld"
)

// server implements the helloworld.GreeterServer interface.
type server struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello is an implementation of the SayHello method from the definition of
// the Greeter service.
func (s *server) SayHello(ctx context.Context, req *helloworld.HelloRequest) (resp *helloworld.HelloReply, err error) {
	resp = &helloworld.HelloReply{}
	resp.Message = "Hello " + req.Name
	return
}

func main() {
	app := wgo.New()

	// Our server implementation.
	s := &server{}

	// The NATS handler from the helloworld.nrpc.proto file.
	h := helloworld.NewGreeterHandler(context.TODO(), app.GetNats(), s)

	// Start a NATS subscription using the handler. You can also use the
	// QueueSubscribe() method for a load-balanced set of servers.
	err := app.Subscribe(h.Subject(), h.Handler)
	if err != nil {
		log.Fatal(err)
	}

	// Keep running until ^C.
	app.Run()
}
