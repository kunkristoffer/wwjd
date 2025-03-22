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

func main() {
	// Setup router
	router := chi.NewRouter()

	// Init openai client
	openai.Init()

	// Serve static assets file
	fs := http.FileServer(http.Dir("assets"))
	router.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	// Home page + POST form
	router.Get("/", func(w http.ResponseWriter, router *http.Request) {
		templ.Handler(index.IndexPage("", models.ChatResponse{})).ServeHTTP(w, router)
	})

	router.Post("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		question := r.FormValue("question")

		rawResponse, err := openai.AskChatGPT(question)
		if err != nil {
			http.Error(w, "AI error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		var chatResp models.ChatResponse
		err = json.Unmarshal([]byte(rawResponse), &chatResp)
		if err != nil {
			http.Error(w, "Failed to parse AI response: "+err.Error(), http.StatusInternalServerError)
			return
		}

		templ.Handler(index.IndexPage(question, chatResp)).ServeHTTP(w, r)
	})

	// Other pages
	router.Get("/best", templ.Handler(best.BestQuestions()).ServeHTTP)
	router.Get("/vote", templ.Handler(vote.VotePage()).ServeHTTP)
	router.Get("/disclaimer", templ.Handler(disclaimer.DisclaimerPage()).ServeHTTP)

	http.ListenAndServe(":8080", router)
}
