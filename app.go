package main

import (
	"log"
	"net/http"

	"github.com/kunkristoffer/wwjd/pages/index"

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

	log.Println("listening on", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
