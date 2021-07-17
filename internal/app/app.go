package app

import (
	"github.com/sloe/gosloe/internal/config"
	"github.com/sloe/gosloe/internal/domain"

	log "github.com/sirupsen/logrus"
)

type App struct {
	config config.Config
	tree   domain.SloeTree
}

func NewApp() *App {
	return &App{}
}

func (app *App) WithConfig(config config.Config) *App {
	app.config = config
	return app
}

func (app *App) Run() error {
	log.Infof("Running app with config %+v", app.config)
	err := app.tree.LoadFromSource(app.config.TreeRoot)
	if err != nil {
		return err
	}
	log.Infof("Tree loaded: %+v", app.tree)
	return nil
}
