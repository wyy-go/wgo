package config

import "time"

var (
	DefaultConfig Config
)

type Config interface {
	Get(string) Value
	Scan(interface{}) error
	Load() error
}

type Value interface {
	Exists() bool
	Bool(def bool) bool
	Int(def int) int
	String(def string) string
	Float64(def float64) float64
	Duration(def time.Duration) time.Duration
	StringSlice(def []string) []string
	StringMap(def map[string]string) map[string]string
	Scan(val interface{}) error
	Bytes() []byte
}

func Get(path string) Value {
	return DefaultConfig.Get(path)
}

func Scan(val interface{}) error {
	return DefaultConfig.Scan(val)
}

func Load() error {
	return DefaultConfig.Load()
}
