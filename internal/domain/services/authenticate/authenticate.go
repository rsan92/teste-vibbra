package authenticate

import (
	"encoding/json"
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
	SearchUser(entitys.User, user.IUserRepository) (entitys.User, error)
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

func (s AuthenticateService) SearchUser(userFromRequest entitys.User, repository user.IUserRepository) (entitys.User, error) {
	response := entitys.User{}
	hashPassword, err := s.Encrypt.GenerateHash(userFromRequest.Password)
	if err != nil {
		return response, err
	}
	userFromDatabase, err := repository.GetUserByLoginAndPassword(userFromRequest.Login, string(hashPassword))
	if err != nil {
		return response, err
	}

	return userFromDatabase, nil
}

func (s AuthenticateService) GenerateWebToken(userFromDatabase entitys.User) (string, error) {
	tokenAsString := ""
	tokenAsString, err := s.Token.CreateToken(userFromDatabase.ID)
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
