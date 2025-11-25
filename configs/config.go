package configs

import (
	"log"
	"os"
	"strconv"
	"time"
)

var (
	databaseURI                  = getEnv("ONCE_DATABASE_URI", "file:app.db?_journal_mode=WAL&_cache_size=2000&_foreign_keys=on&_busy_timeout=5000&_synchronous=NORMAL")
	storageFolder                = getEnv("ONCE_STORAGE_FOLDER", "./uploads")
	storageCleanerInterval       = getEnv("ONCE_STORAGE_CLEANER_INTERVAL", "./uploads")
	httpServerPortString         = getEnv("ONCE_HTTP_SERVER_PORT", "8080")
	httpServerReadTimeoutString  = getEnv("ONCE_HTTP_SERVER_READ_TIMEOUT", "5")
	httpServerWriteTimeoutString = getEnv("ONCE_HTTP_SERVER_WRITE_TIMEOUT", "5")
)

type Config struct {
	DatabaseURI            string
	StorageFolder          string
	StorageCleanerInterval time.Duration
	HTTPServerPort         int
	HTTPServerReadTimeout  time.Duration
	HTTPServerWriteTimeout time.Duration
}

func New() Config {
	return Config{
		DatabaseURI:            databaseURI,
		StorageFolder:          storageFolder,
		StorageCleanerInterval: time.Duration(toIntOrError(storageCleanerInterval)) * time.Second,
		HTTPServerPort:         toIntOrError(httpServerPortString),
		HTTPServerReadTimeout:  time.Duration(toIntOrError(httpServerReadTimeoutString)) * time.Second,
		HTTPServerWriteTimeout: time.Duration(toIntOrError(httpServerWriteTimeoutString)),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	return defaultValue
}

func toIntOrError(content string) int {
	value, err := strconv.Atoi(content)
	if err != nil {
		log.Fatal("error to get http server port", err)
	}
	return value
}
