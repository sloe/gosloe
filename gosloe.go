package main

import (
	"flag"

	"github.com/sloe/gosloe/internal/app"
	"github.com/sloe/gosloe/internal/config"

	log "github.com/sirupsen/logrus"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "configFile", "config.yml", "Config file")
	flag.Parse()

	config, err := config.NewConfigFromYaml(configFile)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{"configFile": configFile}).Fatal("Error loading config")
	}
	var app = app.NewApp().WithConfig(*config)
	app.Run()
}
