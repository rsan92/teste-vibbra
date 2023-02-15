package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rsan92/teste-vibbra/internal/infra/middleware"
)

type Route struct {
	URI         string
	Method      string
	HandlerFunc func(http.ResponseWriter, *http.Request)
	IsPublic    bool
}

func ConfigureEndpoints(r *mux.Router) *mux.Router {
	endpoints := []Route{}

	endpoints = append(endpoints, exportAuthenthicateEndpoints()...)
	endpoints = append(endpoints, exportUserEndpoints()...)

	for _, endpoint := range endpoints {
		if !endpoint.IsPublic {
			r.HandleFunc(endpoint.URI, middleware.Logger(middleware.Auth(endpoint.HandlerFunc))).Methods(endpoint.Method)
		} else {
			r.HandleFunc(endpoint.URI, middleware.Logger(endpoint.HandlerFunc)).Methods(endpoint.Method)
		}
	}
	return r
}
