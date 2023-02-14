package application

import (
	"github.com/gorilla/mux"

	"github.com/rsan92/teste-vibbra/app/application/router"
	"github.com/rsan92/teste-vibbra/internal/infra/configuration"
)

type Application struct {
	Configuration *configuration.AppConfigurations
	HttpRouter    *mux.Router
}

func NewApplication() (*Application, error) {
	app := &Application{}

	app.Configuration = configuration.GetConfigurations()

	router := router.NewApplicationRouter()
	app.HttpRouter = router

	return app, nil
}
