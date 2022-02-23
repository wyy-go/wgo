package main

import (
	"github.com/wyy-go/wgo"
	"github.com/wyy-go/wgo/pkg/logger"
	// This is the package containing the generated *.pb.go and *.nrpc.go
	// files.
	"github.com/wyy-go/wgo/examples/helloworld/helloworld"
)

func main() {
	app := wgo.New()
	nc := app.GetNats()
	defer nc.Close()

	// This is our generated client.
	cli := helloworld.NewGreeterClient(nc)

	// Contact the server and print out its response.
	resp, err := cli.SayHello(helloworld.HelloRequest{Name: "world"})
	if err != nil {
		logger.Fatal(err)
	}
	
	// print
	logger.Debugf("Greeting: %s\n", resp.GetMessage())
}
