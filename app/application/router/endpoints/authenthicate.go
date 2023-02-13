package endpoints

import (
	"net/http"

	"github.com/rsan92/teste-vibbra/app/application/controllers"
)

func exportAuthenthicateEndpoints() []Route {
	authEndpoints := make([]Route, 0)
	authEndpoints = append(authEndpoints, Route{
		URI:         "/v1/authenticate",
		Method:      http.MethodPost,
		HandlerFunc: controllers.AuthenticateHandler_V1,
		IsPublic:    true,
	})

	authEndpoints = append(authEndpoints, Route{
		URI:         "/v1/authenticate/sso",
		Method:      http.MethodPost,
		HandlerFunc: controllers.AuthenticateSSOHandler_V1,
		IsPublic:    true,
	})

	return authEndpoints
}
