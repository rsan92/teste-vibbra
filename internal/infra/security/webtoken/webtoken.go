package webtoken

import (
	"net/http"

	"github.com/rsan92/teste-vibbra/internal/domain/entitys"
)

type ISecurityToken interface {
	CreateToken(user entitys.User) (string, error)
	GetUserFromToken(tokenAsString string) (entitys.User, error)
	ValidateToken(*http.Request) (UserPermissions, error)
}
