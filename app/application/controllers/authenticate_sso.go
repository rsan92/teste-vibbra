package controllers

import (
	"net/http"

	"github.com/rsan92/teste-vibbra/app/application/http_responses"
	"github.com/rsan92/teste-vibbra/internal/domain/models"
	"github.com/rsan92/teste-vibbra/internal/domain/repository/user"
	"github.com/rsan92/teste-vibbra/internal/domain/services/authenticate"
	"github.com/rsan92/teste-vibbra/internal/infra/configuration"
)

func NewAuthenticateSSOController_V1() (*AuthenticateSSOController_V1, error) {
	conf := configuration.GetConfigurations()
	repository, err := user.NewUserRepository(*conf)
	if err != nil {
		return nil, err
	}
	AuthenticateSSOController_V1 := AuthenticateSSOController_V1{
		ExpectedInput:  models.AuthenticateSSOInputRequest{},
		Service:        authenticate.NewAuthenticateSSOService(),
		UserRepository: repository,
		Response:       models.AuthenticateSSOOutputRequest{},
	}
	return &AuthenticateSSOController_V1, nil
}

type AuthenticateSSOController_V1 struct {
	ExpectedInput  models.AuthenticateSSOInputRequest
	Service        authenticate.IAuthenticateSSOService
	UserRepository user.IUserRepository
	Response       models.AuthenticateSSOOutputRequest
}

func AuthenticateSSOHandler_V1(w http.ResponseWriter, r *http.Request) {
	controller, err := NewAuthenticateSSOController_V1()
	if err != nil {
		http_responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	inputFromRequest, err := controller.Service.ValidateInputSSO(r, controller.ExpectedInput)
	if err != nil {
		http_responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userFromDatabase, err := controller.Service.SearchUserByLoginAndToken(r.Context(), inputFromRequest, controller.UserRepository)
	if err != nil {
		http_responses.Error(w, http.StatusUnauthorized, err)
		return
	}

	response := controller.Service.BuildResponseSSO(userFromDatabase, inputFromRequest.Token)

	http_responses.JSON(w, http.StatusOK, response)

}
