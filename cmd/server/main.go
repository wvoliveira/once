package main

import (
	"log/slog"
	"os"

	"github.com/wvoliveira/once/configs"
	"github.com/wvoliveira/once/internal/app/once"
	"github.com/wvoliveira/once/internal/pkg/util"
)

func init() {
	var (
		level       = slog.LevelInfo
		textHandler = slog.NewTextHandler(os.Stdout, util.LogHandler(level))
		logger      = slog.New(textHandler)
	)

	slog.SetDefault(logger)
}

func main() {
	cfg := configs.New()

	storageDB, err := once.NewStorageDB(cfg)
	if err != nil {
		slog.Error("error to create new storage DB", "error", err.Error())
		os.Exit(1)
	}

	storageFiles, err := once.NewStorageFiles(cfg)
	if err != nil {
		slog.Error("error to create new storage files", "error", err.Error())
		os.Exit(1)
	}

	service := once.NewService(storageDB, storageFiles)
	httpServer := once.NewTransportHTTP(cfg, service)
	httpServer.Start()
}
