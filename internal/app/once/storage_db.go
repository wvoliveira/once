package once

import "database/sql"

type StorageDB interface{}

type storageDB struct {
	db *sql.DB
}

func NewStorageDB(uri string) (StorageDB, error) {
	db, err := sql.Open("sqlite", uri)
	if err != nil {
		return nil, err
	}
	return storageDB{db: db}, nil
}
