package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kunkristoffer/wwjd/models"
	"github.com/sashabaranov/go-openai"
)

var Client *openai.Client

func Init() {
	// Load .env file
	if os.Getenv("FLY_APP_NAME") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("❌ Error loading .env file")
		}
	}
	apiKey := os.Getenv("OPENAI_API_KEY")

	// Client init
	Client = openai.NewClient(apiKey)
}

func AskChatGPT(prompt string) (string, error) {
	// Available strings for promp
	moods := []string{"calm", "angry", "joyful", "compassionate", "serious", "disappointed"}
	actions := []string{"glow", "shake", "fade", "pulse", "shine"}
	systemPrompt := fmt.Sprintf(`
You are a helpful assistant that responds in the voice of Jesus. Your answers must be satirical, vague, or funny — but must follow a strict format and be plausible.

Respond ONLY in valid JSON, with no commentary or explanation. The structure must be:

{
  "message": "The actual answer in Jesus' voice",
  "mood": "one of: %s",
  "action": "one of: %s"
}

Important rules:
- Always respond in in Norwegian
- Do not add any text outside the JSON block.
- The mood and action must be chosen from the provided options.
`, strings.Join(moods, ", "), strings.Join(actions, ", "))

	resp, err := Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: systemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	var chatResp models.ChatResponse
	err = json.Unmarshal([]byte(resp.Choices[0].Message.Content), &chatResp)

	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
