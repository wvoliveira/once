package once

import (
	"os"
	"path/filepath"

	"github.com/wvoliveira/once/configs"
)

type StorageFiles interface {
	Write(id, filename string, content []byte) error
}

type storageFiles struct {
	cfg configs.Config
}

func NewStorageFiles(cfg configs.Config) (StorageFiles, error) {
	err := os.MkdirAll(cfg.StorageFolder, 0755)
	if err != nil {
		return nil, err
	}
	return storageFiles{cfg: cfg}, nil
}

func (s storageFiles) GetFile(id string) ([]byte, error) {
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

func (s storageFiles) Write(id, filename string, content []byte) error {
	filePath := filepath.Join(s.cfg.StorageFolder, id)

	// Salva arquivo físico
	if err := os.WriteFile(filePath, content, 0600); err != nil {
		return err
	}
	return nil
}
