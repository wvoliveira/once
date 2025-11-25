package once

import (
	"database/sql"

	"github.com/wvoliveira/once/configs"

	_ "modernc.org/sqlite"
)

type StorageDB interface{}

type storageDB struct {
	db  *sql.DB
	cfg configs.Config
}

func NewStorageDB(cfg configs.Config) (StorageDB, error) {
	db, err := sql.Open("sqlite", cfg.DatabaseURI)
	if err != nil {
		return nil, err
	}
	return storageDB{db: db, cfg: cfg}, nil
}
