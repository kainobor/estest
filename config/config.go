package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type (
	Config struct {
		Server  Server
		Elastic Elastic
	}

	Server struct {
		Port string
	}

	Elastic struct {
		IP   string
		Port string
	}
)

func New(confPath string) (*Config, error) {
	c := new(Config)
	if _, err := toml.DecodeFile(confPath, c); err != nil {
		return nil, fmt.Errorf("can't decode config file `%s`: %v", confPath, err)
	}

	return c, nil
}
