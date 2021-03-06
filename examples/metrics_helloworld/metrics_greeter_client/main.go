package main

import (
	"fmt"
	"github.com/wyy-go/wgo"
	"log"
	"net/http"
	// This is the package containing the generated *.pb.go and *.nrpc.go
	// files.
	"github.com/wyy-go/wgo/examples/metrics_helloworld/helloworld"

	// If you've used the prometheus plugin when generating the code, you
	// can import the HTTP handler of Prometheus to serve up the metrics.
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
		log.Fatal(err)
	}

	// print
	fmt.Printf("Greeting: %s\n", resp.Message)

	// Do this block only if you generated the code with the prometheus plugin.
	fmt.Println("Check metrics at http://localhost:6061/metrics. Hit ^C to exit.")
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":6061", nil)
}
