package router

import (
	"github.com/gorilla/mux"

	"github.com/rsan92/teste-vibbra/app/application/router/endpoints"
)

func NewApplicationRouter() *mux.Router {
	router := mux.NewRouter()
	endpoints.ConfigureEndpoints(router)
	return router
}
