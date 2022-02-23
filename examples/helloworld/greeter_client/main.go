package main

import (
	"github.com/wyy-go/wgo"
	"log"

	// This is the package containing the generated *.pb.go and *.nrpc.go
	// files.
	"github.com/wyy-go/wgo/examples/helloworld/helloworld"
)

func main() {

	nc := wgo.NewNats()
	defer nc.Close()

	// This is our generated client.
	cli := helloworld.NewGreeterClient(nc)

	// Contact the server and print out its response.
	resp, err := cli.SayHello(helloworld.HelloRequest{Name: "world"})
	if err != nil {
		log.Fatal(err)
	}

	// print
	log.Printf("Greeting: %s\n", resp.GetMessage())
}
