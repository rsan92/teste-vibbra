package endpoints

import (
	"net/http"

	"github.com/rsan92/teste-vibbra/app/application/controllers"
)

func exportUserEndpoints() []Route {
	userEndpoints := make([]Route, 0)
	userEndpoints = append(userEndpoints, Route{
		URI:         "/v1/users/{ID}",
		Method:      http.MethodGet,
		HandlerFunc: controllers.UserGetHandler_V1,
		IsPublic:    false,
	})

	userEndpoints = append(userEndpoints, Route{
		URI:         "/v1/users/{ID}",
		Method:      http.MethodPut,
		HandlerFunc: controllers.UserUpsertHandler_V1,
		IsPublic:    false,
	})

	userEndpoints = append(userEndpoints, Route{
		URI:         "/v1/users",
		Method:      http.MethodPost,
		HandlerFunc: controllers.UserUpsertHandler_V1,
		IsPublic:    false,
	})

	return userEndpoints
}
