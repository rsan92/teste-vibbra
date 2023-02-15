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

type UserUpsertController_V1 struct {
	ExpectedInput  models.UpsertUserInputRequest
	Service        service.IUserUpsertService
	UserRepository user.IUserRepository
	Response       models.UpsertUserOutputRequest
}

func NewUserUpsertController_V1() (*UserUpsertController_V1, error) {
	conf := configuration.GetConfigurations()
	repository, err := user.NewUserRepository(*conf)
	if err != nil {
		return nil, err
	}
	UserUpsertController_V1 := UserUpsertController_V1{
		ExpectedInput:  models.UpsertUserInputRequest{},
		Service:        service.NewUpsertUserService(),
		UserRepository: repository,
		Response:       models.UpsertUserOutputRequest{},
	}
	return &UserUpsertController_V1, nil
}

func UserUpsertHandler_V1(w http.ResponseWriter, r *http.Request) {
	controller, err := NewUserUpsertController_V1()
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
	validInput, err := controller.Service.ValidateUpsertUserInput(r, controller.ExpectedInput, userID)

	if err != nil {
		http_responses.Error(w, http.StatusBadRequest, err)
		return
	}

	newUser, err := controller.Service.UpsertUser(validInput, controller.UserRepository, r.Method)

	if err != nil {
		http_responses.Error(w, http.StatusInternalServerError, err)
		return
	}

	response := controller.Service.BuildUserUpsertResponse(newUser)

	var status int
	if r.Method == "POST" {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}
	http_responses.JSON(w, status, response)

}
