package controllers

import (
	"net/http"

	"github.com/rsan92/teste-vibbra/app/application/http_responses"
	"github.com/rsan92/teste-vibbra/internal/domain/models"
	"github.com/rsan92/teste-vibbra/internal/domain/repository/user"
	"github.com/rsan92/teste-vibbra/internal/domain/services/authenticate"
	"github.com/rsan92/teste-vibbra/internal/infra/configuration"
)

func NewAuthenticateController_V1() (*AuthenticateController_V1, error) {
	conf := configuration.GetConfigurations()
	repository, err := user.NewUserRepository(*conf)
	if err != nil {
		return nil, err
	}
	AuthenticateController_V1 := AuthenticateController_V1{
		ExpectedInput:  models.AuthenticateInputRequest{},
		Service:        authenticate.NewAuthenticateService(),
		UserRepository: repository,
		Response:       models.AuthenticateResponseOutput{},
	}
	return &AuthenticateController_V1, nil
}

type AuthenticateController_V1 struct {
	ExpectedInput  models.AuthenticateInputRequest
	Service        authenticate.IAuthenticateService
	UserRepository user.IUserRepository
	Response       models.AuthenticateResponseOutput
}

func AuthenticateHandler_V1(w http.ResponseWriter, r *http.Request) {
	controller, err := NewAuthenticateController_V1()
	if err != nil {
		http_responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	//defer controller.UserRepository.End()

	userFromRequest, err := controller.Service.ValidateInput(r, controller.ExpectedInput)
	if err != nil {
		http_responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userFromDatabase, err := controller.Service.SearchUser(userFromRequest, controller.UserRepository)
	if err != nil {
		http_responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	if userFromDatabase.Login == "" {
		http_responses.Error(w, http.StatusNotFound, err)
		return
	}

	token, err := controller.Service.GenerateWebToken(userFromDatabase)
	if err != nil {
		http_responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	response := controller.Service.BuildResponse(userFromDatabase, token)

	http_responses.JSON(w, http.StatusOK, response)

}

func (c AuthenticateController_V1) AuthenticateSSOHandler_V1(w http.ResponseWriter, r *http.Request) {

}
