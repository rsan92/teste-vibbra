package models

import "github.com/rsan92/teste-vibbra/internal/domain/entitys"

type (
	GetUserInputRequest struct {
		ID         uint64 `json:"id"`
		IsSameUser bool
	}

	GetUserOutputRequest struct {
		User entitys.User `json:"user"`
	}

	UpsertUserInputRequest struct {
		ID       uint64   `json:"id"`
		Name     string   `json:"name"`
		Email    string   `json:"email"`
		Login    string   `json:"login"`
		Password string   `json:"password"`
		Location Location `json:"location"`
	}

	Location struct {
		Lat     float64 `json:"lat"`
		Ing     float64 `json:"ing"`
		Address string  `json:"address"`
		City    string  `json:"city"`
		State   string  `json:"state"`
		ZipCode int64   `json:"zip_code"`
	}

	UpsertUserOutputRequest struct {
		User entitys.User `json:"user"`
	}
)
