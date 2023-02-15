package authenticate

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

type IAuthenticateService interface {
	ValidateInput(*http.Request, models.AuthenticateInputRequest) (entitys.User, error)
	SearchUserByLoginAndPassword(entitys.User, user.IUserRepository) (entitys.User, error)
	GenerateWebToken(entitys.User) (string, error)
	BuildResponse(entitys.User, string) models.AuthenticateResponseOutput
}

type AuthenticateService struct {
	Encrypt encryptdata.IEncrypt
	Token   webtoken.ISecurityToken
}

func NewAuthenticateService() *AuthenticateService {
	return &AuthenticateService{
		Encrypt: encryptdata.NewDataEncrypt(),
		Token:   webtoken.NewJWTSecurityToken(),
	}
}

func (s AuthenticateService) ValidateInput(r *http.Request, expectedInput models.AuthenticateInputRequest) (entitys.User, error) {
	response := entitys.User{}
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		return response, err
	}

	if err := json.Unmarshal(requestBody, &expectedInput); err != nil {
		return response, err
	}

	response.Login = expectedInput.Login
	response.Password = expectedInput.Password

	return response, nil
}

func (s AuthenticateService) SearchUserByLoginAndPassword(userFromRequest entitys.User, repository user.IUserRepository) (entitys.User, error) {
	response := entitys.User{}

	userFromDatabase, err := repository.GetUserByLogin(userFromRequest.Login)
	if err != nil {
		return response, err
	}

	if userFromRequest.Password != userFromDatabase.Password {
		isDiff := s.Encrypt.VerifyHash(userFromRequest.Password, userFromDatabase.Password)

		if isDiff != nil {
			hashedPass, err := s.Encrypt.GenerateHash(userFromRequest.Password)
			if err != nil {
				return entitys.User{}, errors.New("can not encrypt password")
			}
			userFromRequest.Password = string(hashedPass)
		} else if isDiff == nil {
			userFromRequest.Password = userFromDatabase.Password
		}
	}
	if userFromRequest.Password != userFromDatabase.Password {
		return response, errors.New("invalid password")
	}

	return userFromDatabase, nil
}

func (s AuthenticateService) GenerateWebToken(userFromDatabase entitys.User) (string, error) {
	tokenAsString := ""
	tokenAsString, err := s.Token.CreateToken(userFromDatabase)
	if err != nil {
		return tokenAsString, err
	}

	return tokenAsString, nil
}

func (s AuthenticateService) BuildResponse(userFromDatabase entitys.User, token string) models.AuthenticateResponseOutput {
	response := models.AuthenticateResponseOutput{}
	response.Token = token
	response.User = userFromDatabase
	return response
}
