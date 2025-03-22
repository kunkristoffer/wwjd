package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
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
	database.InitDB("data/prompts.db")

	// Init session
	var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_VOTE_KEY")))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}

	// Start server
	addr := ":8080"
	logger.Info("Starting server", slog.String("address", addr))
	if err := http.ListenAndServe(addr, New(store)); err != nil {
		logger.Error("Server crashed", slog.String("error", err.Error()))
	}
}
