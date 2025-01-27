package main

// Generating swagger doc spec
//go:generate swag init --parseDependency -g ../../pkg/api/api.go -d ../../pkg/api -o ../../pkg/api/docs

import (
	"flag"

	"github.com/gin-contrib/cors"

	"hms/gateway/pkg/api"
	_ "hms/gateway/pkg/api/docs"
	"hms/gateway/pkg/config"
)

func main() {

	var (
		cfgPath = flag.String("config", "./config.json", "config file path")
	)
	flag.Parse()

	cfg, err := config.New(*cfgPath)
	if err != nil {
		panic(err)
	}

	a := api.New(cfg).Build()

	//TODO complete CORS config
	a.Use(cors.Default())

	if err = a.Run(cfg.Host); err != nil {
		panic(err)
	}
}
