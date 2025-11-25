package once

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/wvoliveira/once/configs"
	"github.com/wvoliveira/once/internal/pkg/util"

	_ "modernc.org/sqlite"
)

type StorageDB interface {
	Save(fileID, filename string, ttl time.Time) error
}

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

func (s storageDB) Save(fileID, filename string, ttl time.Time) error {
	var (
		id      = util.GenerateID()
		timeNow = time.Now()
	)

	_, err := s.db.Exec(`
		INSERT INTO content (id, file_id, file_name, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		id, fileID, filename, fmt.Sprintf("%s", ttl.Format(time.RFC3339)), timeNow.Format(time.RFC3339),
	)
	return err
}

func (s storageDB) Get(fileID string) error {
	row := s.db.QueryRow(`
		SELECT file_id, expires_at FROM content WHERE file_id = ?`,
		fileID,
	)

	var ExpiresAt string
	err := row.Scan(&fileID, &ExpiresAt)
	if err != nil {
		return err
	}

	expiresAt, err := time.Parse(time.RFC3339, ExpiresAt)
	if err != nil {
		return err
	}

	if time.Now().After(expiresAt) {
		return ErrContentExpired
	}

	return err
}
