package config

import (
	"go.uber.org/multierr"

	"github.com/wyy-go/wgo/pkg/config"
	"github.com/wyy-go/wgo/pkg/logger"
)

type Config struct {
	Mysql Mysql
	Redis Redis
}

type Mysql struct {
	DriverName     string `json:"driver"`
	DataSourceName string `json:"data_source"`
	MaxIdleConn    int    `json:"max_idle"`
	MaxOpenConn    int    `json:"max_open"`
}

type Redis struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

var c Config

func Load() {
	err := multierr.Combine(
		config.Get("mysql").Scan(&c.Mysql),
		config.Get("redis").Scan(&c.Redis),
	)
	if err != nil {
		logger.Fatal()
	}
}

func GetMysql() Mysql {
	return c.Mysql
}

func GetRedis() Redis {
	return c.Redis
}
