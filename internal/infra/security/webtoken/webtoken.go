package webtoken

import "net/http"

type ISecurityToken interface {
	CreateToken(userID uint64) (string, error)
	ValidateToken(request *http.Request) error
	GetUserID(request *http.Request) (uint64, error)
}
