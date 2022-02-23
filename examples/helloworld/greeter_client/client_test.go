package main

import (
	"github.com/wyy-go/wgo"
	"github.com/wyy-go/wgo/examples/helloworld/helloworld"
	"log"
	"testing"
)

func BenchmarkClient(b *testing.B) {
	app := wgo.New()
	nc := app.GetNats()
	defer nc.Close()
	// This is our generated client.
	cli := helloworld.NewGreeterClient(nc)
	for i := 0; i < b.N; i++ {
		_, err := cli.SayHello(helloworld.HelloRequest{Name: "world"})
		if err != nil {
			log.Fatal(err)
		}
	}
}
