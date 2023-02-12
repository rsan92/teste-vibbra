package webtoken

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTSecurityToken struct{}

func NewJWTSecurityToken() JWTSecurityToken {
	return JWTSecurityToken{}
}

func (jst JWTSecurityToken) CreateToken(userID uint64) (string, error) {
	permissions := jwt.MapClaims{}
	permissions["auth"] = true
	permissions["user_id"] = userID
	permissions["expires"] = time.Now().Add(time.Hour * 2).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, permissions)
	tokenAsString, err := token.SignedString("it_will_be_a_secret_key")
	if err != nil {
		return "", err
	}
	return tokenAsString, nil
}

func (jst JWTSecurityToken) ValidateToken(request *http.Request) error {
	tokenAsString := jst.extractTokenFromRequest(request)

	token, err := jwt.Parse(tokenAsString, jst.getVerificationKey)

	if err != nil {
		return err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return nil
	}

	return errors.New("invalid securitytoken")
}

func (jst JWTSecurityToken) GetUserID(request *http.Request) (uint64, error) {
	tokenAsString := jst.extractTokenFromRequest(request)

	token, err := jwt.Parse(tokenAsString, jst.getVerificationKey)

	if err != nil {
		return 0, err
	}

	if permissions, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, err := strconv.ParseUint(fmt.Sprintf("%.0f", permissions["user_id"]), 10, 64)
		if err != nil {
			return 0, err
		}
		return userID, nil
	}

	return 0, errors.New("invalid securitytoken or invalid permissions")
}

func (jst JWTSecurityToken) extractTokenFromRequest(req *http.Request) string {
	token := req.Header.Get("Authorization")

	if len(strings.Split(token, " ")) == 2 {
		return strings.Split(token, " ")[1]
	}
	return ""
}

func (jst JWTSecurityToken) getVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("metodo de assinatura inesperado %v", token.Header["alg"])
	}

	return "config.SecretKey", nil
}
