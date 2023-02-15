package controllers

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/rsan92/teste-vibbra/app/application/http_responses"
	"github.com/rsan92/teste-vibbra/internal/domain/models"
	"github.com/rsan92/teste-vibbra/internal/domain/repository/user"
	service "github.com/rsan92/teste-vibbra/internal/domain/services/user"
	"github.com/rsan92/teste-vibbra/internal/infra/configuration"
)

type UserGetController_V1 struct {
	ExpectedInput  models.GetUserInputRequest
	Service        service.IGetUserService
	UserRepository user.IUserRepository
	Response       models.GetUserOutputRequest
}

func NewUserGetController_V1() (*UserGetController_V1, error) {
	conf := configuration.GetConfigurations()
	repository, err := user.NewUserRepository(*conf)
	if err != nil {
		return nil, err
	}
	UserGetController_V1 := UserGetController_V1{
		ExpectedInput:  models.GetUserInputRequest{},
		Service:        service.NewGetUserService(),
		UserRepository: repository,
		Response:       models.GetUserOutputRequest{},
	}
	return &UserGetController_V1, nil
}

func UserGetHandler_V1(w http.ResponseWriter, r *http.Request) {
	controller, err := NewUserGetController_V1()
	if err != nil {
		http_responses.Error(w, http.StatusInternalServerError, err)
		return
	}
	params := mux.Vars(r)
	userID, err := strconv.ParseUint(params["ID"], 10, 64)
	if err != nil {
		http_responses.Error(w, http.StatusBadRequest, err)
		return
	}
	validInput, err := controller.Service.ValidateGetInput(r, controller.ExpectedInput, userID)

	if err != nil {
		http_responses.Error(w, http.StatusBadRequest, err)
		return
	}

	userFromDatabase, err := controller.Service.GetUser(validInput, controller.UserRepository)
	if err != nil {
		http_responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	response := controller.Service.BuildResponseGetUser(userFromDatabase)

	http_responses.JSON(w, http.StatusOK, response)

}
