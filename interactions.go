package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Session struct {
	Permanent string        `json:"permanent"`
	Dynamic   []ChatMessage `json:"dynamic"`
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatPayload struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Temperature float32       `json:"temperature"`
}

type APIResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func getAssistantReply(config Settings, session Session) (string, error) {
	messages := make([]ChatMessage, 0)
	messages = append(messages, ChatMessage{Role: "system", Content: session.Permanent})
	messages = append(messages, session.Dynamic...)

	payload := ChatPayload{
		Model:       config.Model,
		Messages:    messages,
		Temperature: config.Temperature,
	}

	payloadJSON, err := json.Marshal(payload)

	if err != nil {
		return "", fmt.Errorf("Failed to marshal payload: %v", err)
	}

	req, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(payloadJSON))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.APIKey))
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var apiResponse APIResponse

	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", fmt.Errorf("Failed to unmarshal response: %v", err)
	}

	if resp.StatusCode != 200 || len(apiResponse.Choices) == 0 {
		return "", fmt.Errorf("received an invalid response from API: HTTP %d, Body: %s", resp.StatusCode, string(body))
	}

	return apiResponse.Choices[0].Message.Content, nil
}
