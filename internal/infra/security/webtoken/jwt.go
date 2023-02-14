package webtoken

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/rsan92/teste-vibbra/internal/domain/entitys"
	"github.com/rsan92/teste-vibbra/internal/infra/configuration"
)

type JWTSecurityToken struct{}

func NewJWTSecurityToken() JWTSecurityToken {
	return JWTSecurityToken{}
}

func (jst JWTSecurityToken) CreateToken(user entitys.User) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["authorized"] = true
	permissions["user_id"] = user.ID
	permissions["user_login"] = user.Login
	permissions["user_pass"] = user.Password
	permissions["exp"] = time.Now().Add(time.Hour * 6).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)

	tokenAsString, err := token.SignedString([]byte(configuration.GetConfigurations().SecretKey))
	if err != nil {
		return "", err
	}
	return tokenAsString, nil
}

func (jst JWTSecurityToken) GetUserFromToken(tokenAsString string) (entitys.User, error) {
	token, err := jwt.Parse(tokenAsString, jst.getVerificationKey)

	if err != nil {
		return entitys.User{}, err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		response := entitys.User{}
		var err error
		response.ID, err = strconv.ParseUint(fmt.Sprintf("%.0f", permissions["user_id"]), 10, 64)
		if err != nil {
			return response, err
		}
		response.Login = fmt.Sprintf("%v", permissions["user_login"])
		response.Password = fmt.Sprintf("%v", permissions["user_pass"])

		return response, nil
	}

	return entitys.User{}, errors.New("invalid securitytoken or invalid permissions")
}

func (jst JWTSecurityToken) getVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("metodo de assinatura inesperado %v", token.Header["alg"])
	}

	return []byte(configuration.GetConfigurations().SecretKey), nil
}
