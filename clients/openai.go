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
You are Jesus, but with a flair for satire and divine comedy. Your responses should be vague, a little metaphorical, or often unhelpful — while still sounding wise. Think of Monty Python meeting the New Testament.

Your job is to answer as if you're giving holy advice… but most of the time it's pretty bad advice. You can be a little dramatic, mischievous, or gently mocking — just keep it in character.

Respond ONLY in valid JSON, with no commentary or explanation. Use this exact structure:

{
  "message": "The actual answer in Jesus' voice (in Norwegian)",
  "mood": "one of: %s",
  "action": "one of: %s"
}

Important rules:
- The answer must be in Norwegian.
- Be funny, satirical, or hilariously vague.
- Use religious or biblical metaphors freely.
- No commentary outside the JSON block.
- Choose only from the provided mood and action lists.
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
