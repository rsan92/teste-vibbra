package webtoken

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/rsan92/teste-vibbra/internal/domain/entitys"
	"github.com/rsan92/teste-vibbra/internal/infra/configuration"
)

type tokenObject struct {
	Token    string `json:"token"`
	TokenApp string `json:"app_token"`
}

type UserPermissions struct {
	Authorized bool   `json:"authorized"`
	UserID     uint64 `json:"user_id"`
	UserLogin  string `json:"user_login"`
	UserPass   string `json:"user_pass"`
	Exp        int64  `json:"exp"`
}

func (t *tokenObject) String() string {
	if t.Token != "" {
		return t.Token
	}
	if t.TokenApp != "" {
		return t.TokenApp
	}
	return ""
}

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

func (jst JWTSecurityToken) ValidateToken(req *http.Request) (UserPermissions, error) {
	tokenAsString, err := jst.extractTokenFromRequest(req)
	if err != nil {
		return UserPermissions{}, err
	}

	realToken, err := jwt.Parse(tokenAsString, jst.getVerificationKey)
	if err != nil {
		return UserPermissions{}, err
	}

	if mapToken, ok := realToken.Claims.(jwt.MapClaims); ok && realToken.Valid {
		authorized, ok := mapToken["authorized"].(bool)
		if !ok {
			return UserPermissions{}, errors.New("invalid authorized value in token")
		}

		userID, ok := mapToken["user_id"].(float64)
		if !ok {
			return UserPermissions{}, errors.New("invalid user_id value in token")
		}

		userLogin, ok := mapToken["user_login"].(string)
		if !ok {
			return UserPermissions{}, errors.New("invalid user_login value in token")
		}

		userPass, ok := mapToken["user_pass"].(string)
		if !ok {
			return UserPermissions{}, errors.New("invalid user_pass value in token")
		}

		Exp, ok := mapToken["exp"].(float64)
		if !ok {
			return UserPermissions{}, errors.New("invalid expiration value in token")
		}

		return UserPermissions{
			Authorized: authorized,
			UserID:     uint64(userID),
			UserLogin:  userLogin,
			UserPass:   userPass,
			Exp:        int64(Exp),
		}, nil
	}

	return UserPermissions{}, errors.New("invalid token")
}

func (jst JWTSecurityToken) extractTokenFromRequest(req *http.Request) (string, error) {
	var token string
	token = req.Header.Get("Authorization")
	if token == "" {
		bodyReq, err := io.ReadAll(req.Body)
		if err != nil {
			return "", err
		}
		tk := tokenObject{}
		err = json.Unmarshal(bodyReq, &tk)
		if err != nil {
			return "", err
		}
		token = tk.String()

		if token == "" {
			return "", fmt.Errorf("empty token")
		}
	} else {
		if len(strings.Split(token, " ")) == 2 {
			return strings.Split(token, " ")[1], nil
		} else {
			return "", fmt.Errorf("invalid token format in header")
		}
	}
	return token, nil
}

func (jst JWTSecurityToken) getVerificationKey(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		return nil, fmt.Errorf("metodo de assinatura inesperado %v", token.Header["alg"])
	}

	return []byte(configuration.GetConfigurations().SecretKey), nil
}
