package main

import (
	"encoding/json"
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	openai "github.com/kunkristoffer/wwjd/clients"
	"github.com/kunkristoffer/wwjd/models"
	"github.com/kunkristoffer/wwjd/pages/best"
	"github.com/kunkristoffer/wwjd/pages/disclaimer"
	"github.com/kunkristoffer/wwjd/pages/index"
	"github.com/kunkristoffer/wwjd/pages/vote"
)

func New() http.Handler {
	r := chi.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	// GET homepage
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(index.IndexPage("", models.ChatResponse{})).ServeHTTP(w, r)
	})

	// POST form
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		question := r.FormValue("question")

		raw, err := openai.AskChatGPT(question)
		if err != nil {
			http.Error(w, "AI error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var chatResp models.ChatResponse
		if err := json.Unmarshal([]byte(raw), &chatResp); err != nil {
			http.Error(w, "Failed to parse AI response", http.StatusInternalServerError)
			return
		}

		templ.Handler(index.IndexPage(question, chatResp)).ServeHTTP(w, r)
	})

	// Other pages
	r.Get("/best", templ.Handler(best.BestQuestions()).ServeHTTP)
	r.Get("/vote", templ.Handler(vote.VotePage()).ServeHTTP)
	r.Get("/disclaimer", templ.Handler(disclaimer.DisclaimerPage()).ServeHTTP)

	return r
}
