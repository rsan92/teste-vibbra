package models

import "github.com/rsan92/teste-vibbra/internal/domain/entitys"

type (
	AuthenticateInputRequest struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}

	AuthenticateResponseOutput struct {
		Token string       `json:"token"`
		User  entitys.User `json:"user"`
	}
)
