package file

import (
	"encoding/json"
	"fmt"

	"github.com/wyy-go/wgo/pkg/config"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type file struct {
	v      *viper.Viper
	values *config.JSONValues
	opts   config.Options
}

func (f *file) Init(opts ...config.Option) error {
	for _, o := range opts {
		o(&f.opts)
	}

	return nil
}

func NewConfig(opts ...config.Option) (config.Config, error) {
	options := config.Options{
		Type: "yaml",
		Path: "config.yaml",
	}
	f := &file{
		v:    viper.New(),
		opts: options,
	}
	if err := f.Init(opts...); err != nil {
		return nil, err
	}

	if err := f.Load(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *file) Load() error {
	if f.opts.Type != "" {
		f.v.SetConfigType(f.opts.Type)
	}

	if f.opts.Path != "" {
		f.v.SetConfigFile(f.opts.Path)
	}

	if err := f.v.ReadInConfig(); err != nil {
		return err
	}

	m := f.v.AllSettings()
	b, _ := json.Marshal(m)
	f.values = config.NewJSONValues(b)

	f.v.OnConfigChange(func(e fsnotify.Event) {
		m := f.v.AllSettings()
		b, _ := json.Marshal(m)
		if string(f.values.Bytes()) != string(b) {
			f.values = config.NewJSONValues(b)
			fmt.Println("Config file changed:", e.Name)
		}
	})
	f.v.WatchConfig()
	return nil
}

func (f *file) Get(path string) config.Value {
	return f.values.Get(path)
}

func (f *file) Scan(val interface{}) error {
	return f.values.Scan(val)
}
