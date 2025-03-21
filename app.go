package main

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"

	"github.com/kunkristoffer/wwjd/pages/best"
	"github.com/kunkristoffer/wwjd/pages/disclaimer"
	"github.com/kunkristoffer/wwjd/pages/index"
	"github.com/kunkristoffer/wwjd/pages/vote"
)

func main() {
	r := chi.NewRouter()

	// Home page + POST form
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		templ.Handler(index.Index("", "")).ServeHTTP(w, r)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		question := r.FormValue("question")
		response := "You asked: " + question // Replace with AI logic
		templ.Handler(index.Index(question, response)).ServeHTTP(w, r)
	})

	// Other pages
	r.Get("/best", templ.Handler(best.BestQuestions()).ServeHTTP)
	r.Get("/vote", templ.Handler(vote.VotePage()).ServeHTTP)
	r.Get("/disclaimer", templ.Handler(disclaimer.DisclaimerPage()).ServeHTTP)

	http.ListenAndServe(":8080", r)
}
