package wgo

import (
	"fmt"
	"github.com/nats-io/nats.go"
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

func (a *App) Subscribe(subj string, cb nats.MsgHandler) error {
	sub, err := a.nc.Subscribe(subj, cb)
	a.sub = sub
	return err
}

// QueueSubscribe for a load-balanced set of servers
func (a *App) QueueSubscribe(subj string, queue string, cb nats.MsgHandler) error {
	sub, err := a.nc.QueueSubscribe(subj, queue, cb)
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

	// Connect to the NATS server.
	nc, err := nats.Connect(app.opts.natsUrl, nats.Timeout(app.opts.natsTimeOut))
	if err != nil {
		log.Fatal(err)
	}
	app.nc = nc

	return app
}
