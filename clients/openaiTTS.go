package openai

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func GenerateSpeech(text string) ([]byte, error) {
	payload := map[string]string{
		"model": "tts-1",
		"input": text,
		"voice": "fable", // or fable/echo
	}
	body, _ := json.Marshal(payload)

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/audio/speech", bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}
