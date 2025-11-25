package once

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/wvoliveira/once/configs"
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

func (s storageDB) GetFile(id string) ([]byte, error) {
	filePath := filepath.Join(s.cfg.StorageFolder, id)
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Limpeza final: Apaga o arquivo físico
	// (Fazemos isso em background ou defer para não atrasar a resposta,
	//  mas aqui faremos síncrono por simplicidade)
	defer os.Remove(filePath)
	return data, nil
}
