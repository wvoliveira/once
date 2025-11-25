package once

import (
	"time"
)

type Service interface {
	Upload(id, filename string, data []byte, ttl time.Duration) error
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

func (s service) Upload(id, filename string, data []byte, ttl time.Duration) error {
	err := s.storageDB.Save(id, filename, ttl)
	if err != nil {
		return err
	}

	err = s.storageFiles.Write(id, filename, data)
	if err != nil {
		return err
	}
	return nil
}

func (s service) Download(id string) ([]byte, Content, error) {
	return nil, Content{}, nil
}
