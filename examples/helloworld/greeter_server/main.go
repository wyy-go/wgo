package main

import (
	"context"
	"github.com/wyy-go/wgo"
	"log"
	// This is the package containing the generated *.pb.go and *.nrpc.go
	// files.
	"github.com/wyy-go/wgo/examples/helloworld/helloworld"
	"github.com/wyy-go/wgo/middleware/m1"
	"github.com/wyy-go/wgo/middleware/m2"
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
	log.Println("======service handler======")
	return
}

func main() {
	app := wgo.New(wgo.Middleware(m1.M1(), m2.M2()))

	// Our server implementation.
	s := &server{}

	// The NATS handler from the helloworld.nrpc.proto file.
	h := helloworld.NewGreeterHandler(s)

	// Start a NATS subscription using the handler. You can also use the
	// RegisterHandlerForLB() method for a load-balanced set of servers.
	err := app.RegisterHandler(h)
	if err != nil {
		log.Fatal(err)
	}

	// Keep running until ^C.
	app.Run()
}
