package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	openai "github.com/kunkristoffer/wwjd/clients"
	"github.com/kunkristoffer/wwjd/database"
	"github.com/kunkristoffer/wwjd/models"
	"github.com/kunkristoffer/wwjd/pages/best"
	"github.com/kunkristoffer/wwjd/pages/disclaimer"
	"github.com/kunkristoffer/wwjd/pages/index"
	"github.com/kunkristoffer/wwjd/pages/newest"
	"github.com/kunkristoffer/wwjd/pages/tired"
	"github.com/kunkristoffer/wwjd/pages/vote"
	"github.com/kunkristoffer/wwjd/sessions"
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
		// Get counter from session
		session, _ := sessions.ChatStore.Get(r, "chat-session")
		count := 0
		if v, ok := session.Values["ask_count"].(int); ok {
			count = v
		}

		// Increment counter and check if above 10
		count++
		session.Values["ask_count"] = count
		session.Save(r, w)
		if count >= 10 {
			http.Redirect(w, r, "/tired", http.StatusSeeOther)
			return
		}

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

		// Check if reply already exists, prevents dup on reload, fix tba
		var exists bool
		err = database.DB.QueryRow(
			"SELECT EXISTS(SELECT 1 FROM prompts WHERE question = ? LIMIT 1)",
			question,
		).Scan(&exists)
		if err != nil {
			slog.Error("DB check failed", slog.String("error", err.Error()))
		} else if !exists {
			_, errDB := database.DB.Exec(
				"INSERT INTO prompts (question, reply) VALUES (?, ?)",
				question,
				chatResp.Message,
			)
			if errDB != nil {
				slog.Error("DB insert failed", slog.String("error", errDB.Error()))
			} else {
				slog.Info("Prompt saved", slog.String("question", question))
			}
		} else {
			slog.Info("Reply already exists, skipping insert", slog.String("reply", chatResp.Message))
		}

		// TTS
		responseText := chatResp.Message
		audioBytes, err := openai.GenerateSpeech(responseText) // 🆕 write this helper
		if err == nil {
			filename := fmt.Sprintf("assets/audio/%d.mp3", time.Now().UnixNano())
			os.WriteFile(filename, audioBytes, 0644)
			chatResp.AudioURL = "/" + filename // serve via static assets
		}

		templ.Handler(index.IndexPage(question, chatResp)).ServeHTTP(w, r)
	})

	// Other pages
	r.Get("/best", func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.DB.Query("SELECT id, date_asked, question, reply, votes FROM prompts ORDER BY votes DESC LIMIT 10")
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var prompts []models.Prompt
		for rows.Next() {
			var p models.Prompt
			err := rows.Scan(&p.ID, &p.DateAsked, &p.Question, &p.Reply, &p.Votes)
			if err != nil {
				http.Error(w, "Failed to parse row: "+err.Error(), http.StatusInternalServerError)
				return
			}
			prompts = append(prompts, p)
		}
		templ.Handler(best.BestQuestions(prompts)).ServeHTTP(w, r)
	})

	r.Get("/newest", func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.DB.Query("SELECT id, date_asked, question, reply, votes FROM prompts ORDER BY date_asked DESC LIMIT 10")
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var prompts []models.Prompt
		for rows.Next() {
			var p models.Prompt
			err := rows.Scan(&p.ID, &p.DateAsked, &p.Question, &p.Reply, &p.Votes)
			if err != nil {
				http.Error(w, "Failed to parse row: "+err.Error(), http.StatusInternalServerError)
				return
			}
			prompts = append(prompts, p)
		}
		templ.Handler(newest.NewestQuestions(prompts)).ServeHTTP(w, r)
	})

	r.Get("/vote", func(w http.ResponseWriter, r *http.Request) {
		rows, err := database.DB.Query("SELECT id, question, reply, votes FROM prompts ORDER BY RANDOM() LIMIT 3")
		if err != nil {
			http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var prompts []models.Prompt
		for rows.Next() {
			var p models.Prompt
			err := rows.Scan(&p.ID, &p.Question, &p.Reply, &p.Votes)
			if err != nil {
				http.Error(w, "Scan error: "+err.Error(), http.StatusInternalServerError)
				return
			}
			prompts = append(prompts, p)
		}

		templ.Handler(vote.VotePage(prompts)).ServeHTTP(w, r)
	})

	r.Post("/vote", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		id := r.FormValue("id")

		// Get session
		session, _ := sessions.VoteStore.Get(r, "vote-session")
		voted := session.Values["voted"]

		// Check if user has already voted
		if voted == true {
			slog.Info("User has already voted", slog.String("id", id))
			http.Redirect(w, r, "/vote", http.StatusSeeOther)
			return
		} else {
			session.Values["voted"] = true
			err := session.Save(r, w)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

		_, err := database.DB.Exec("UPDATE prompts SET votes = votes + 1 WHERE id = ?", id)
		if err != nil {
			http.Error(w, "Failed to register vote", http.StatusInternalServerError)
			return
		}

		slog.Info("Voted!", slog.String("question", id))

		// Redirect to refresh the list
		http.Redirect(w, r, "/vote", http.StatusSeeOther)
	})

	r.Get("/disclaimer", templ.Handler(disclaimer.DisclaimerPage()).ServeHTTP)
	r.Get("/tired", templ.Handler(tired.TiredPage()).ServeHTTP)
	return r
}
