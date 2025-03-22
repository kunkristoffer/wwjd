package sessions

import (
	"os"

	"github.com/gorilla/sessions"
)

// Init sessions
var VoteStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_VOTE_KEY")))
var ChatStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_VOTE_KEY")))

func init() {
	VoteStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   10,
		HttpOnly: true,
	}

	ChatStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	}
}
