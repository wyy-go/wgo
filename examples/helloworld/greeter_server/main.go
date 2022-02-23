package main

import (
	"context"
	"github.com/wyy-go/wgo"
	appconfig "github.com/wyy-go/wgo/examples/helloworld/config"
	"github.com/wyy-go/wgo/pkg/config"
	"github.com/wyy-go/wgo/pkg/env"
	"github.com/wyy-go/wgo/pkg/logger"
	// This is the package containing the generated *.pb.go and *.nrpc.go
	// files.
	"github.com/wyy-go/wgo/examples/helloworld/helloworld"
	"github.com/wyy-go/wgo/examples/helloworld/middleware/m1"
	"github.com/wyy-go/wgo/examples/helloworld/middleware/m2"
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
	logger.Debug("======service handler======")
	return
}

func main() {
	app := wgo.New(wgo.Middleware(m1.M1(), m2.M2()))

	// deploy env
	deployEnv := config.Get("service.deploy_env").String("")
	env.SetDeploy(env.ToDeploy(deployEnv))
	if env.GetDeploy() == env.DeployUnknown {
		logger.Fatal("未设置发布模式")
	}

	appconfig.Load()
	//	logger.Debug(appconfig.GetMysql())
	//	logger.Debug(appconfig.GetRedis())
	
	// Our server implementation.
	s := &server{}

	// The NATS handler from the helloworld.nrpc.proto file.
	h := helloworld.NewGreeterHandler(s)

	// Start a NATS subscription using the handler. You can also use the
	// RegisterHandlerForLB() method for a load-balanced set of servers.
	err := app.RegisterHandlerForLB(h)
	if err != nil {
		logger.Fatal(err)
	}

	// Keep running until ^C.
	app.Run()
}
