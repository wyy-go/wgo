package wgo

import (
	"context"
	"github.com/nats-io/nats.go"
	"time"
)

type Options struct {
	name        string
	version     string
	metadata    map[string]string
	deployEnv   string
	verbose     bool
	trace       bool
	natsUrl     string
	natsTimeOut time.Duration
	context     context.Context
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	options := Options{
		metadata: map[string]string{},
		verbose:  false,
		trace:    true,
	}

	for _, o := range opts {
		o(&options)
	}

	if options.natsUrl == "" {
		options.natsUrl = nats.DefaultURL
	}

	if options.natsTimeOut == 0 {
		options.natsTimeOut = 5 * time.Second
	}

	//if options.Registry == nil {
	//	options.Registry = mdns.NewRegistry()
	//}

	return options
}

func Name(name string) Option {
	return func(o *Options) {
		o.name = name
	}
}

func Version(version string) Option {
	return func(o *Options) {
		o.version = version
	}
}

func Metadata(md map[string]string) Option {
	return func(o *Options) {
		o.metadata = md
	}
}

func Verbose(v bool) Option {
	return func(o *Options) {
		o.verbose = v
	}
}

func Trace(t bool) Option {
	return func(o *Options) {
		o.trace = t
	}
}

func NatsUrl(url string) Option {
	return func(o *Options) {
		o.natsUrl = url
	}
}

func NatsTimeOut(d time.Duration) Option {
	return func(o *Options) {
		o.natsTimeOut = d
	}
}

//func Registry(r registry.Registry) Option {
//	return func(o *Options) {
//		o.Registry = r
//	}
//}
