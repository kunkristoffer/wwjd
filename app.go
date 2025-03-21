package main

import (
	"net/http"

	"github.com/kunkristoffer/wwjd/pages/best"
	"github.com/kunkristoffer/wwjd/pages/disclaimer"
	"github.com/kunkristoffer/wwjd/pages/index"
	"github.com/kunkristoffer/wwjd/pages/vote"

	"github.com/a-h/templ"
)

func main() {
	// Root page with optional form handling
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			r.ParseForm()
			question := r.FormValue("question")
			// Replace this with real AI call later
			response := "You asked: " + question
			templ.Handler(index.Index(question, response)).ServeHTTP(w, r)
			return
		}

		// GET: show empty form
		templ.Handler(index.Index("", "")).ServeHTTP(w, r)
	})

	// Static pages
	http.Handle("/best", templ.Handler(best.BestQuestions()))
	http.Handle("/vote", templ.Handler(vote.VotePage()))
	http.Handle("/disclaimer", templ.Handler(disclaimer.DisclaimerPage()))

	// Start the server
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	component := templates.Index("hello", "world")
	templ.Handler(component).ServeHTTP(w, r)
}
