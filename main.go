package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	openai "github.com/kunkristoffer/wwjd/clients"
	"github.com/kunkristoffer/wwjd/database"
)

func main() {
	// Init logger
	var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Init OpenAI client
	openai.Init()

	// Init db
	path := "data"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			slog.Info("Error creating data folder")
		}
	}
	database.InitDB()

	// Start server
	addr := ":8080"
	logger.Info("Starting server", slog.String("address", addr))
	if err := http.ListenAndServe(addr, New()); err != nil {
		logger.Error("Server crashed", slog.String("error", err.Error()))
	}
}
