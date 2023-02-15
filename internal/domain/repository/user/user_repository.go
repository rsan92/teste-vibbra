package user

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/rsan92/teste-vibbra/internal/domain/entitys"
	"github.com/rsan92/teste-vibbra/internal/domain/models"
	"github.com/rsan92/teste-vibbra/internal/infra/configuration"
	"github.com/rsan92/teste-vibbra/internal/infra/database"
)

const (
	TABLE_NAME = "usuarios"
)

var (
	INSERT_USER_QUERY              = fmt.Sprintf("INSERT INTO %v (name, email, login, password, lat, ing, address, city, state, zip_code) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", TABLE_NAME)
	UPDATE_USER_QUERY              = fmt.Sprintf("UPDATE %v SET name=?, email=?, password=?, lat=?, ing=?, address=?, city=?, state=?, zip_code=? WHERE id=?;", TABLE_NAME)
	SELECT_USER_BY_LOGIN_QUERY     = fmt.Sprintf("SELECT id, name, email, login, password, lat, ing, address, city, state, zip_code FROM %v WHERE login=?", TABLE_NAME)
	SELECT_USERID_BY_LOGIN_QUERY   = fmt.Sprintf("SELECT id FROM %v WHERE login=?", TABLE_NAME)
	SELECT_USERPASS_BY_LOGIN_QUERY = fmt.Sprintf("SELECT password FROM %v WHERE login=?", TABLE_NAME)
	SELECT_USER_BY_ID_QUERY        = fmt.Sprintf("SELECT id, name, email, login, password, lat, ing, address, city, state, zip_code FROM %v WHERE id=?", TABLE_NAME)
)

type IUserRepository interface {
	GetUserByLogin(string) (entitys.User, error)
	GetUserByID(uint64) (entitys.User, error)
	InsertUser(models.UpsertUserInputRequest) (entitys.User, error)
	UpdateUser(models.UpsertUserInputRequest) (entitys.User, error)
	GetUserPass(string) (string, error)
	GetOtherUserByID(uint64) (entitys.User, error)
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

func (u UserRepository) GetUserByLogin(login string) (entitys.User, error) {
	rows, err := u.database.Query(SELECT_USER_BY_LOGIN_QUERY, login)

	if err != nil {
		return entitys.User{}, err
	}

	if rows.Next() {
		var user entitys.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Login,
			&user.Password,
			&user.Location.Lat,
			&user.Location.Ing,
			&user.Location.Address,
			&user.Location.City,
			&user.Location.State,
			&user.Location.ZipCode,
		); err != nil {
			return entitys.User{}, err
		}
		return user, nil
	}
	return entitys.User{}, nil
}

func (u UserRepository) GetOtherUserByID(id uint64) (entitys.User, error) {
	otherUserInfo, err := u.GetUserByID(id)
	if err != nil {
		return entitys.User{}, err
	}
	otherUserInfo.ClearSensitiveInformation()
	return otherUserInfo, nil

}

func (u UserRepository) GetUserByID(id uint64) (entitys.User, error) {
	rows, err := u.database.Query(SELECT_USER_BY_ID_QUERY, id)

	if err != nil {
		return entitys.User{}, err
	}

	if rows.Next() {
		var user entitys.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Login,
			&user.Password,
			&user.Location.Lat,
			&user.Location.Ing,
			&user.Location.Address,
			&user.Location.City,
			&user.Location.State,
			&user.Location.ZipCode,
		); err != nil {
			return entitys.User{}, err
		}
		return user, nil
	}
	return entitys.User{}, nil
}

func (u UserRepository) getUserID(login string) (uint64, error) {
	rows, err := u.database.Query(SELECT_USERID_BY_LOGIN_QUERY, login)

	if err != nil {
		return 0, err
	}

	if rows.Next() {
		var userID uint64
		if err := rows.Scan(
			&userID,
		); err != nil {
			return 0, err
		}
		return userID, nil
	}
	return 0, nil
}

func (u UserRepository) GetUserPass(login string) (string, error) {
	rows, err := u.database.Query(SELECT_USERPASS_BY_LOGIN_QUERY, login)

	if err != nil {
		return "", err
	}

	if rows.Next() {
		var userPass string
		if err := rows.Scan(
			&userPass,
		); err != nil {
			return "", err
		}
		return userPass, nil
	}
	return "", nil
}

func (u UserRepository) UpdateUser(newUser models.UpsertUserInputRequest) (entitys.User, error) {
	var result entitys.User
	statement, err := u.database.Prepare(UPDATE_USER_QUERY)

	if err != nil {
		return result, fmt.Errorf("failed to prepare query_statement. err: %v", err)
	}

	defer statement.Close()

	db_result, err := statement.Exec(newUser.Name, newUser.Email, newUser.Password, newUser.Location.Lat, newUser.Location.Ing, newUser.Location.Address, newUser.Location.City, newUser.Location.State, newUser.Location.ZipCode, newUser.ID)

	if err != nil {
		return result, fmt.Errorf("failed to update user. err: %v", err)
	}
	if rows, _ := db_result.RowsAffected(); rows == 0 {
		return result, errors.New("no rows affected")
	}
	userID, err := u.getUserID(newUser.Login)
	if err != nil {
		return result, fmt.Errorf("failed to get id from user. err: %v", err)

	}

	result = entitys.User{
		ID:       uint64(userID),
		Name:     newUser.Name,
		Email:    newUser.Email,
		Login:    newUser.Login,
		Password: newUser.Password,
		Location: entitys.Location{
			Lat:     newUser.Location.Lat,
			Ing:     newUser.Location.Ing,
			Address: newUser.Location.Address,
			City:    newUser.Location.City,
			State:   newUser.Location.State,
			ZipCode: newUser.Location.ZipCode,
		},
	}

	return result, nil
}

func (u UserRepository) InsertUser(newUser models.UpsertUserInputRequest) (entitys.User, error) {
	var result entitys.User
	statement, err := u.database.Prepare(INSERT_USER_QUERY)

	if err != nil {
		return result, fmt.Errorf("failed to prepare query_statement. err: %v", err)
	}

	defer statement.Close()

	db_result, err := statement.Exec(newUser.Name, newUser.Email, newUser.Login, newUser.Password, newUser.Location.Lat, newUser.Location.Ing,
		newUser.Location.Address, newUser.Location.City, newUser.Location.State, newUser.Location.ZipCode)

	if err != nil {
		return result, fmt.Errorf("failed to insert new user. err: %v", err)
	}
	userID, err := db_result.LastInsertId()
	if err != nil {
		return result, fmt.Errorf("failed to get id from new user. err: %v", err)

	}

	result = entitys.User{
		ID:       uint64(userID),
		Name:     newUser.Name,
		Email:    newUser.Email,
		Login:    newUser.Login,
		Password: newUser.Password,
		Location: entitys.Location{
			Lat:     newUser.Location.Lat,
			Ing:     newUser.Location.Ing,
			Address: newUser.Location.Address,
			City:    newUser.Location.City,
			State:   newUser.Location.State,
			ZipCode: newUser.Location.ZipCode,
		},
	}

	return result, nil
}

func (u UserRepository) End() error {
	return u.database.Close()
}
