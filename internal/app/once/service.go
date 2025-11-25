package once

import (
	"time"

	"github.com/wvoliveira/once/internal/pkg/util"
)

type Service interface {
	Upload(filename string, data []byte, ttl time.Time) (string, error)
	Download(id string) ([]byte, Content, error)
}

type service struct {
	storageDB    StorageDB
	storageFiles StorageFiles
}

func NewService(storageDB StorageDB, storageFiles StorageFiles) Service {
	return service{
		storageDB:    storageDB,
		storageFiles: storageFiles,
	}
}

func (s service) Upload(filename string, data []byte, ttl time.Time) (string, error) {
	fileID := util.GenerateID()

	err := s.storageDB.Save(fileID, filename, ttl)
	if err != nil {
		return fileID, err
	}

	err = s.storageFiles.Write(fileID, filename, data)
	if err != nil {
		return fileID, err
	}
	return fileID, nil
}

func (s service) Download(id string) ([]byte, Content, error) {
	content, err := s.storageDB.Get(id)
	if err != nil {
		return nil, Content{}, err
	}

	data, err := s.storageFiles.GetFile(id)
	if err != nil {
		return nil, Content{}, err
	}

	return data, content, nil
}
