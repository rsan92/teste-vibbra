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

type (
	IUserUpsertService interface {
		ValidateUpsertUserInput(*http.Request, models.UpsertUserInputRequest, uint64) (models.UpsertUserInputRequest, error)
		UpsertUser(models.UpsertUserInputRequest, user.IUserRepository, string) (entitys.User, error)
		BuildUserUpsertResponse(entitys.User) models.UpsertUserOutputRequest
	}
)

type UpsertUserService struct {
	Encrypt encryptdata.IEncrypt
	Token   webtoken.ISecurityToken
}

func NewUpsertUserService() *UpsertUserService {
	return &UpsertUserService{
		Encrypt: encryptdata.NewDataEncrypt(),
		Token:   webtoken.NewJWTSecurityToken(),
	}
}

func (service UpsertUserService) ValidateUpsertUserInput(r *http.Request, expectedInput models.UpsertUserInputRequest, userID uint64) (models.UpsertUserInputRequest, error) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		return expectedInput, err
	}

	if err := json.Unmarshal(requestBody, &expectedInput); err != nil {
		return expectedInput, err
	}

	expectedInput.ID = userID

	tokenPermissionsRaw := r.Context().Value("tokenPermissions")
	tokenPermissions, ok := tokenPermissionsRaw.(webtoken.UserPermissions)
	if !ok {
		return expectedInput, errors.New("invalid tokenPermissions")
	}
	if expectedInput.ID != tokenPermissions.UserID {
		return expectedInput, errors.New("invalid userID for update")
	}

	return expectedInput, nil
}

func (service UpsertUserService) UpsertUser(input models.UpsertUserInputRequest, repository user.IUserRepository, method string) (entitys.User, error) {
	var newUser entitys.User
	var err error
	oldpass, err := repository.GetUserPass(input.Login)
	if err != nil {
		return entitys.User{}, errors.New("can not get oldpass from database")
	}

	if input.Password != oldpass {
		var isDiff error
		if oldpass != "" {
			isDiff = service.Encrypt.VerifyHash(input.Password, oldpass)
		}

		if isDiff != nil || oldpass == "" {
			hashedPass, err := service.Encrypt.GenerateHash(input.Password)
			if err != nil {
				return entitys.User{}, errors.New("can not encrypt password")
			}
			input.Password = string(hashedPass)
		} else if isDiff == nil && oldpass != "" {
			input.Password = oldpass
		}
	}

	if method == http.MethodPut {
		newUser, err = repository.UpdateUser(input)
	} else if method == http.MethodPost {
		newUser, err = repository.InsertUser(input)
	} else {
		return entitys.User{}, errors.New("invalid http method")
	}

	if err != nil {
		return entitys.User{}, err
	}

	return newUser, nil
}

func (service UpsertUserService) BuildUserUpsertResponse(newUser entitys.User) models.UpsertUserOutputRequest {
	return models.UpsertUserOutputRequest{
		User: newUser,
	}
}
