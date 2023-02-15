package authenticate

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/rsan92/teste-vibbra/internal/domain/entitys"
	"github.com/rsan92/teste-vibbra/internal/domain/models"
	"github.com/rsan92/teste-vibbra/internal/domain/repository/user"
	"github.com/rsan92/teste-vibbra/internal/infra/security/encryptdata"
	"github.com/rsan92/teste-vibbra/internal/infra/security/webtoken"
)

type IAuthenticateSSOService interface {
	ValidateInputSSO(*http.Request, models.AuthenticateSSOInputRequest) (models.AuthenticateSSOInputRequest, error)
	SearchUserByLoginAndToken(context.Context, models.AuthenticateSSOInputRequest, user.IUserRepository) (entitys.User, error)
	BuildResponseSSO(entitys.User, string) models.AuthenticateSSOOutputRequest
}

type AuthenticateSSOService struct {
	Encrypt encryptdata.IEncrypt
	Token   webtoken.ISecurityToken
}

func NewAuthenticateSSOService() *AuthenticateSSOService {
	return &AuthenticateSSOService{
		Encrypt: encryptdata.NewDataEncrypt(),
		Token:   webtoken.NewJWTSecurityToken(),
	}
}

func (s AuthenticateSSOService) ValidateInputSSO(r *http.Request, expectedInput models.AuthenticateSSOInputRequest) (models.AuthenticateSSOInputRequest, error) {
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		return expectedInput, err
	}

	if err := json.Unmarshal(requestBody, &expectedInput); err != nil {
		return expectedInput, err
	}

	return expectedInput, nil
}

func (s AuthenticateSSOService) SearchUserByLoginAndToken(ctx context.Context, input models.AuthenticateSSOInputRequest, repository user.IUserRepository) (entitys.User, error) {
	response := entitys.User{}

	userFromDatabase, err := repository.GetUserByLogin(input.Login)
	if err != nil {
		return response, err
	}

	return userFromDatabase, nil
}

func (s AuthenticateSSOService) BuildResponseSSO(userFromDatabase entitys.User, token string) models.AuthenticateSSOOutputRequest {
	response := models.AuthenticateSSOOutputRequest{}
	response.Token = token
	response.User = userFromDatabase
	return response
}
