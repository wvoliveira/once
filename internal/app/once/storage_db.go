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
	Get(fileID string) (Content, error)
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
		index   = util.GenerateID()
		timeNow = time.Now()
	)

	_, err := s.db.Exec(`
		INSERT INTO content (id, file_id, file_name, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		index, fileID, filename, fmt.Sprintf("%s", ttl.Format(time.RFC3339)), timeNow.Format(time.RFC3339),
	)
	return err
}

func (s storageDB) Get(fileID string) (Content, error) {
	row := s.db.QueryRow(`
		DELETE FROM content 
        WHERE file_id = ?
        RETURNING id, file_name, expires_at, created_at`,
		fileID,
	)

	var (
		id        string
		fileName  string
		expiresAt string
		createdAt string
	)

	err := row.Scan(&id, &fileName, &expiresAt, &createdAt)
	if err != nil {
		return Content{}, err
	}

	ExpiresAt, err := time.Parse(time.RFC3339, expiresAt)
	if err != nil {
		return Content{}, err
	}

	CreatedAt, err := time.Parse(time.RFC3339, createdAt)
	if err != nil {
		return Content{}, err
	}

	if time.Now().After(ExpiresAt) {
		return Content{}, ErrContentExpired
	}

	return Content{
		ID:        id,
		FileName:  fileName,
		ExpiresAt: ExpiresAt,
		CreatedAt: CreatedAt,
	}, err
}
