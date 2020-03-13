package manage

import (
	"database/sql"

	"github.com/leeif/go-web-template/config"
	"github.com/leeif/go-web-template/log"
)

type Manager struct {
	logger *log.Log
	config *config.Config
	db     *sql.DB
}

func NewManager(db *sql.DB, config *config.Config, logger *log.Log) *Manager {
	return &Manager{
		logger: logger.With("compoment", "manager"),
		config: config,
		db:     db,
	}
}
