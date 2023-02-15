package middleware

import (
	"context"
	"net/http"

	"github.com/rsan92/teste-vibbra/app/application/http_responses"
	"github.com/rsan92/teste-vibbra/internal/infra/security/webtoken"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := webtoken.NewJWTSecurityToken().ValidateToken(r)
		if err != nil {
			http_responses.Error(w, http.StatusUnauthorized, err)
			return
		}
		ctx := r.Context()
		ctx = context.WithValue(ctx, "tokenPermissions", token)
		req := r.Clone(ctx)
		next(w, req)
	}
}
