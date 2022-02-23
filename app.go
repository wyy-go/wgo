package wgo

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"github.com/nats-io/nats.go"
	"github.com/wyy-go/wgo/nrpc"
	"github.com/wyy-go/wgo/pkg/logger"
	"github.com/wyy-go/wgo/pkg/logger/zap"
	"github.com/wyy-go/wgo/pkg/uuid"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type App struct {
	opts Options

	nc  *nats.Conn
	sub *nats.Subscription
}

func shutdown() []os.Signal {
	return []os.Signal{
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL,
	}
}

func (a *App) GetNats() *nats.Conn {
	return a.nc
}

func (a *App) RegisterHandler(h nrpc.H) error {
	h.SetNats(a.nc)
	h.SetMiddleware(a.opts.middleware)
	sub, err := a.nc.Subscribe(h.Subject(), h.Handler)
	a.sub = sub
	return err
}

// RegisterHandlerForLB for a load-balanced set of servers
func (a *App) RegisterHandlerForLB(h nrpc.H) error {
	queue := uuid.New().String()
	h.SetNats(a.nc)
	h.SetMiddleware(a.opts.middleware)
	sub, err := a.nc.QueueSubscribe(h.Subject(), queue, h.Handler)
	a.sub = sub
	return err
}

func (a *App) Run() error {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, shutdown()...)

	fmt.Println("server is running, ^C quits.")
	fmt.Printf("received signal %s", <-ch)
	if a.nc != nil {
		a.nc.Close()
	}
	if a.sub != nil {
		a.sub.Unsubscribe()
	}

	close(ch)
	return nil
}

func New(opts ...Option) *App {
	// app
	options := newOptions(opts...)
	app := &App{}
	app.opts = options

	// logger:
	//  level: "debug"
	//  filename: "log/app.log"
	//  max_size: 100
	//  max_backups: 10
	//  max_age: 7
	//  compress: true

	// logger
	w := &lumberjack.Logger{
		Filename:   "log/app.log",
		MaxSize:    100,
		MaxBackups: 10,
		MaxAge:     7,
		Compress:   true,
	}

	lvl, err := logger.GetLevel("debug")
	if err != nil {
		log.Fatal(err)
	}

	l, err := zap.NewLogger(logger.WithLevel(lvl), logger.WithWriter([]io.Writer{os.Stderr, w}))
	if err != nil {
		logger.Fatal(err)
	}
	logger.DefaultLogger = l

	// Connect to the NATS server.
	nc, err := nats.Connect(app.opts.natsUrl, nats.Timeout(app.opts.natsTimeOut))
	if err != nil {
		log.Fatal(err)
	}
	app.nc = nc
	return app
}
