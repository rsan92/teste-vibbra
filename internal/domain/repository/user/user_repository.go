package user

import (
	"database/sql"
	"fmt"

	"github.com/rsan92/teste-vibbra/internal/domain/entitys"
	"github.com/rsan92/teste-vibbra/internal/infra/configuration"
	"github.com/rsan92/teste-vibbra/internal/infra/database"
)

const (
	TABLE_NAME = "usuarios"
)

type IUserRepository interface {
	GetUserByLoginAndPassword(login, password string) (entitys.User, error)
	GetUserByID(id uint64) (entitys.User, error)
	End() error
}

type UserRepository struct {
	database *sql.DB
}

func NewUserRepository(conf configuration.AppConfigurations) (*UserRepository, error) {
	var err error
	userRepository := &UserRepository{}
	database := database.SQLDatabase{}
	userRepository.database, err = database.GetInstance(conf)
	if err != nil {
		return nil, err
	}
	return userRepository, nil
}

func (u UserRepository) GetUserByLoginAndPassword(login, password string) (entitys.User, error) {
	query := fmt.Sprintf("SELECT id, login, password FROM %v WHERE login=? AND password=?", TABLE_NAME)
	rows, err := u.database.Query(query, login, password)

	if err != nil {
		return entitys.User{}, err
	}

	if rows.Next() {
		var user entitys.User
		if err := rows.Scan(
			&user.ID,
			&user.Login,
			&user.Password,
		); err != nil {
			return entitys.User{}, err
		}
		return user, nil
	}
	return entitys.User{}, nil
}

func (u UserRepository) GetUserByID(id uint64) (entitys.User, error) {
	return entitys.User{}, nil
}

func (u UserRepository) End() error {
	return u.database.Close()
}
