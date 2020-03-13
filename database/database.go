package database

import (
	// database driver

	"database/sql"

	"github.com/leeif/go-web-template/config"
)

func NewDatabase(config *config.Config) (*sql.DB, error) {
	dbcfg := config.Database
	db, err := sql.Open("mysql", generateConnectionSchema(dbcfg))
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(dbcfg.MaxOpenConns)
	db.SetMaxIdleConns(dbcfg.MaxIdleConns)
	return db, nil
}

func generateConnectionSchema(dbcfg *config.DatabaseConfig) string {
	switch dbcfg.Type.String() {
	case "mysql":
		return dbcfg.User + ":" + dbcfg.Password + "@tcp(" + dbcfg.Host + ":" + dbcfg.Port.String() + ")/" + dbcfg.DB +
			"?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return ""
}
