package main

import (
	"context"
	"flag"

	"github.com/kainobor/estest/app/lib/elsearch"
	"github.com/kainobor/estest/app/lib/logger"
	"github.com/kainobor/estest/app/lib/session"
	"github.com/kainobor/estest/app/srv"
	"github.com/kainobor/estest/config"
)

func main() {
	ctx := context.Background()

	var confPath string
	var isProd bool
	flag.StringVar(&confPath, "conf", "", "Config path")
	flag.BoolVar(&isProd, "prod", false, "Is current environment production")
	flag.Parse()
	if confPath == "" {
		panic("Flag `-conf` for config path is required")
	}
	logger.Init(isProd)
	log := logger.New(ctx)

	conf, err := config.New(confPath)
	if err != nil {
		log.Fatalw("Can't read config", "error", err)
	}

	es, err := elsearch.New(conf.Elastic)
	if err != nil {
		log.Fatalw("Can't connect to elasticsearch", "error", err)
	}

	session.Init()

	server := srv.New(conf.Server)
	server.Start(ctx, es)
}
