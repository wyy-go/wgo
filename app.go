package wgo

import (
	"flag"
	"github.com/natefinch/lumberjack"
	"github.com/nats-io/nats.go"
	"github.com/wyy-go/wgo/nrpc"
	"github.com/wyy-go/wgo/pkg/config"
	"github.com/wyy-go/wgo/pkg/config/file"
	"github.com/wyy-go/wgo/pkg/logger"
	"github.com/wyy-go/wgo/pkg/logger/zap"
	"github.com/wyy-go/wgo/pkg/uuid"
	"io"
	"os"
	"os/signal"
	"syscall"
)

var cfgFile string

func init() {
	flag.StringVar(&cfgFile, "config", "config.yaml", "config file")
}

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
	h.SetMiddleware(a.opts.Middleware)
	sub, err := a.nc.Subscribe(h.Subject(), h.Handler)
	a.sub = sub
	return err
}

// RegisterHandlerForLB for a load-balanced set of servers
func (a *App) RegisterHandlerForLB(h nrpc.H) error {
	queue := uuid.New().String()
	h.SetNats(a.nc)
	h.SetMiddleware(a.opts.Middleware)
	sub, err := a.nc.QueueSubscribe(h.Subject(), queue, h.Handler)
	a.sub = sub
	return err
}

func (a *App) Run() error {

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, shutdown()...)

	logger.Infof("Name: %s Version: %s", a.opts.Name, a.opts.Version)
	logger.Info("server is running, ^C quits.")
	logger.Info("received signal %s", <-ch)
	if a.nc != nil {
		a.nc.Close()
	}
	if a.sub != nil {
		a.sub.Unsubscribe()
	}

	close(ch)
	return nil
}

func NewNats(opts ...Option) *nats.Conn {
	// app
	options := newOptions(opts...)

	// Connect to the NATS server.
	nc, err := nats.Connect(options.NatsUrl, nats.Timeout(options.NatsTimeOut))
	if err != nil {
		logger.Fatal(err)
	}

	return nc
}

func New(opts ...Option) *App {

	flag.Parse()
	_, err := os.Stat(cfgFile)
	if os.IsNotExist(err) {
		logger.Fatal("config file not exists")
	}

	config.DefaultConfig, _ = file.NewConfig(config.Path(cfgFile))
	conf := Config{}
	if err := config.Scan(&conf); err != nil {
		logger.Fatal(err)
	}

	// app
	options := newOptions(opts...)
	app := &App{}
	options.Name = conf.Service.Name
	options.Version = conf.Service.Version
	options.Verbose = conf.Service.Verbose
	options.DeployEnv = conf.Service.DeployEnv

	app.opts = options

	// logger
	w := &lumberjack.Logger{
		Filename:   conf.Logger.Filename,
		MaxSize:    conf.Logger.MaxSize,
		MaxBackups: conf.Logger.MaxBackups,
		MaxAge:     conf.Logger.MaxAge,
		Compress:   conf.Logger.Compress,
	}

	lvl, err := logger.GetLevel(conf.Logger.Level)
	if err != nil {
		logger.Fatal(err)
	}

	l, err := zap.NewLogger(logger.WithLevel(lvl), logger.WithWriter([]io.Writer{os.Stderr, w}))
	if err != nil {
		logger.Fatal(err)
	}
	logger.DefaultLogger = l

	// Connect to the NATS server.
	nc, err := nats.Connect(app.opts.NatsUrl, nats.Timeout(app.opts.NatsTimeOut))
	if err != nil {
		logger.Fatal(err)
	}
	app.nc = nc
	return app
}
