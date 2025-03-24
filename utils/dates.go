package utils

import (
	"fmt"
	"time"
)

func TimeAgo(t time.Time) string {
	duration := time.Since(t)

	switch {
	case duration < time.Minute:
		return "akkurat nå"
	case duration < time.Hour:
		return fmt.Sprintf("for %d minutter siden", int(duration.Minutes()))
	case duration < 24*time.Hour:
		return fmt.Sprintf("for %d timer siden", int(duration.Hours()))
	case duration < 48*time.Hour:
		return "i går"
	case duration < 7*24*time.Hour:
		return fmt.Sprintf("for %d dager siden", int(duration.Hours()/24))
	case duration < 30*24*time.Hour:
		return fmt.Sprintf("for %d uker siden", int(duration.Hours()/(24*7)))
	case duration < 365*24*time.Hour:
		return fmt.Sprintf("for %d måneder siden", int(duration.Hours()/(24*30)))
	default:
		return fmt.Sprintf("for %d år siden", int(duration.Hours()/(24*365)))
	}
}
