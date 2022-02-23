package etcd

import (
	"encoding/json"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"github.com/wyy-go/wgo/pkg/config"
)

type etcd struct {
	v      *viper.Viper
	values *config.JSONValues
	opts   config.Options
}

func (e *etcd) Init(opts ...config.Option) error {
	for _, o := range opts {
		o(&e.opts)
	}

	return nil
}

func NewConfig(opts ...config.Option) (config.Config, error) {
	options := config.Options{
		Type: "yaml",
		Path: "zgo/config.yaml",
	}
	e := &etcd{
		v:    viper.New(),
		opts: options,
	}
	if err := e.Init(opts...); err != nil {
		return nil, err
	}

	if err := e.Load(); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *etcd) Load() error {
	//if e.opts.Type != "" {
	//	e.v.SetConfigType(e.opts.Type)
	//}

	//if e.opts.Context == nil {
	//	return errors.New("address can not be empty")
	//}
	//addr := e.opts.Context.Value(addressKey{}).(string)
	//e.v.AddRemoteProvider("etcd", addr, e.opts.Path)
	e.v.AddRemoteProvider("etcd", "http://127.0.0.1:2379", e.opts.Path)
	if e.opts.Type != "" {
		e.v.SetConfigType(e.opts.Type)
	}

	if err := e.v.ReadRemoteConfig(); err != nil {
		return err
	}

	m := e.v.AllSettings()
	b, _ := json.Marshal(m)
	e.values = config.NewJSONValues(b)

	e.v.OnConfigChange(func(evt fsnotify.Event) {
		m := e.v.AllSettings()
		b, _ := json.Marshal(m)
		if string(e.values.Bytes()) != string(b) {
			e.values = config.NewJSONValues(b)
			fmt.Println("Config file changed:", evt.Name)
		}
	})

	_ = e.v.WatchRemoteConfigOnChannel()

	return nil
}

func (e *etcd) Get(path string) config.Value {
	return e.values.Get(path)
}

func (e *etcd) Scan(val interface{}) error {
	return e.values.Scan(val)
}
