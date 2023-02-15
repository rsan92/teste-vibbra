package database

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"

	"github.com/rsan92/teste-vibbra/internal/communs/helpers"
	"github.com/rsan92/teste-vibbra/internal/infra/configuration"
)

var DabaseInstance *sql.DB

type ISQLDatabases interface {
	GetInstance() *sql.DB
}

type SQLDatabase struct{}

func (db SQLDatabase) GetInstance(conf configuration.AppConfigurations) (*sql.DB, error) {
	if DabaseInstance == nil || DabaseInstance.Stats().OpenConnections == 0 {
		dbase, err := db.configuraConnection(conf)
		if err != nil {
			return nil, err
		}
		DabaseInstance = dbase
	}

	return DabaseInstance, nil
}

func (db SQLDatabase) configuraConnection(conf configuration.AppConfigurations) (*sql.DB, error) {
	connString := helpers.GetConnString(conf)
	if connString == "" {
		return nil, fmt.Errorf("invalid connection in database: %v", connString)
	}
	database, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	if err = database.Ping(); err != nil {
		database.Close()
		return nil, err
	}

	return database, nil
}
