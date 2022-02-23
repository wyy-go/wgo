package wgo

import (
	"context"
	"github.com/nats-io/nats.go"
	"github.com/wyy-go/wgo/nrpc"
	"time"
)

type Options struct {
	Name        string
	Version     string
	Metadata    map[string]string
	DeployEnv   string
	Verbose     bool
	Trace       bool
	NatsUrl     string
	NatsTimeOut time.Duration
	Middleware  []nrpc.Middleware
	Context     context.Context
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	options := Options{
		Metadata: map[string]string{},
		Verbose:  false,
		Trace:    true,
	}

	for _, o := range opts {
		o(&options)
	}

	if options.NatsUrl == "" {
		options.NatsUrl = nats.DefaultURL
	}

	if options.NatsTimeOut == 0 {
		options.NatsTimeOut = 5 * time.Second
	}

	//if options.Registry == nil {
	//	options.Registry = mdns.NewRegistry()
	//}

	return options
}

func Name(name string) Option {
	return func(o *Options) {
		o.Name = name
	}
}

func Version(version string) Option {
	return func(o *Options) {
		o.Version = version
	}
}

func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.Metadata = md
	}
}

func Verbose(v bool) Option {
	return func(o *Options) {
		o.Verbose = v
	}
}

func Trace(t bool) Option {
	return func(o *Options) {
		o.Trace = t
	}
}

func NatsUrl(url string) Option {
	return func(o *Options) {
		o.NatsUrl = url
	}
}

func NatsTimeOut(d time.Duration) Option {
	return func(o *Options) {
		o.NatsTimeOut = d
	}
}

func Middleware(m ...nrpc.Middleware) Option {
	return func(o *Options) {
		o.Middleware = append(o.Middleware, m...)
	}
}
