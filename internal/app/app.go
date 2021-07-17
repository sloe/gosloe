package app

import (
	"github.com/sloe/gosloe/internal/config"

	log "github.com/sirupsen/logrus"
)

type App struct {
	config config.Config
}

func NewApp() App {
	return App{}
}

func (app App) WithConfig(config config.Config) App {
	app.config = config
	return app
}

func (app App) Run() error {
	log.Infof("Running app with config %+v", app.config)
	return nil
}
