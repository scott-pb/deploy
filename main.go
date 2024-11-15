package main

import (
	"deploy/config"
	"deploy/log"
	"deploy/router"
	"flag"
)

var configName string

func init() {
	flag.StringVar(&configName, "conf", "./config/config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	log.SetLogFile("deploy.log")
	config.Init(configName)
	router.InitRouter()
}
