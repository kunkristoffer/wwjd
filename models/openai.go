package models

type ChatResponse struct {
	Message  string `json:"message"`
	Mood     string `json:"mood"`
	Action   string `json:"action"`
	AudioURL string `json:"audio_url"`
}
