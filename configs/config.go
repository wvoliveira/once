package configs

import (
	"log"
	"os"
	"strconv"
	"time"
)

var (
	DatabaseURI = getEnv("ONCE_DATABASE_URI", "file:app.db?_journal_mode=WAL&_cache_size=2000&_foreign_keys=on&_busy_timeout=5000&_synchronous=NORMAL")

	httpServerPortString = getEnv("ONCE_HTTP_SERVER_PORT", "8080")
	HTTPServerPort       int

	httpServerReadTimeoutString = getEnv("ONCE_HTTP_SERVER_READ_TIMEOUT", "10")
	HTTPServerReadTimeout       time.Duration

	httpServerWriteTimeoutString = getEnv("ONCE_HTTP_SERVER_WRITE_TIMEOUT", "10")
	HTTPServerWriteTimeout       time.Duration
)

func Init() {
	HTTPServerPort = toIntOrError(httpServerPortString)
	HTTPServerReadTimeout = time.Duration(toIntOrError(httpServerReadTimeoutString)) * time.Second
	HTTPServerWriteTimeout = time.Duration(toIntOrError(httpServerWriteTimeoutString))
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
