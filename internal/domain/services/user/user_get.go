package user

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/rsan92/teste-vibbra/internal/domain/entitys"
	"github.com/rsan92/teste-vibbra/internal/domain/models"
	"github.com/rsan92/teste-vibbra/internal/domain/repository/user"
	"github.com/rsan92/teste-vibbra/internal/infra/security/encryptdata"
	"github.com/rsan92/teste-vibbra/internal/infra/security/webtoken"
)

type IGetUserService interface {
	ValidateGetInput(*http.Request, models.GetUserInputRequest, uint64) (models.GetUserInputRequest, error)
	GetUser(models.GetUserInputRequest, user.IUserRepository) (entitys.User, error)
	BuildResponseGetUser(entitys.User) models.GetUserOutputRequest
}

type GetUserService struct {
	Encrypt encryptdata.IEncrypt
	Token   webtoken.ISecurityToken
}

func NewGetUserService() *GetUserService {
	return &GetUserService{
		Encrypt: encryptdata.NewDataEncrypt(),
		Token:   webtoken.NewJWTSecurityToken(),
	}
}

func (s GetUserService) ValidateGetInput(req *http.Request, input models.GetUserInputRequest, userID uint64) (models.GetUserInputRequest, error) {
	if userID == 0 {
		requestBody, err := io.ReadAll(req.Body)
		if err != nil {
			return input, err
		}

		if err := json.Unmarshal(requestBody, &input); err != nil {
			return input, err
		}
	} else {
		input.ID = userID
	}

	tokenPermissionsRaw := req.Context().Value("tokenPermissions")
	tokenPermissions, ok := tokenPermissionsRaw.(webtoken.UserPermissions)
	if !ok {
		return input, errors.New("invalid tokenPermissions")
	}

	if input.ID != tokenPermissions.UserID {
		input.IsSameUser = false
	} else {
		input.IsSameUser = true
	}

	return input, nil
}

func (s GetUserService) GetUser(input models.GetUserInputRequest, repository user.IUserRepository) (entitys.User, error) {
	if input.IsSameUser {
		userFromDatabase, err := repository.GetUserByID(input.ID)
		if err != nil {
			return entitys.User{}, errors.New("cant get user from database")
		}
		return userFromDatabase, nil
	} else {
		userFromDatabase, err := repository.GetOtherUserByID(input.ID)
		if err != nil {
			return entitys.User{}, errors.New("cant get user from database")
		}
		return userFromDatabase, nil

	}

}

func (s GetUserService) BuildResponseGetUser(input entitys.User) models.GetUserOutputRequest {
	return models.GetUserOutputRequest{User: input}
}
