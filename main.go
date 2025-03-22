package main

import (
	"log/slog"
	"net/http"
	"os"

	openai "github.com/kunkristoffer/wwjd/clients"
)

func main() {
	// Init logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Init OpenAI client
	openai.Init()

	// Start server
	addr := ":8080"
	logger.Info("Starting server", slog.String("address", addr))
	if err := http.ListenAndServe(addr, New()); err != nil {
		logger.Error("Server crashed", slog.String("error", err.Error()))
	}
}
