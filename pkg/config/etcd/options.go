package etcd

import "github.com/wyy-go/wgo/pkg/config"

type addressKey struct{}

func WithAddress(addr string) config.Option {
	return config.SetOption(addressKey{}, addr)
}
