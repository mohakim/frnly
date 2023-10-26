package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type ChatCompletionChunk struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Delta        map[string]string `json:"delta"`
		Index        int               `json:"index"`
		FinishReason string            `json:"finish_reason"`
	} `json:"choices"`
}

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
	Stream      bool          `json:"stream"`
}

type APIResponse struct {
	Choices []struct {
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

func streamCompletion(config Settings, session Session, apiOutput chan<- string) {
	url := "https://api.openai.com/v1/chat/completions"
	messages := make([]ChatMessage, 0)
	messages = append(messages, ChatMessage{Role: "system", Content: session.Permanent})
	messages = append(messages, session.Dynamic...)

	payload := ChatPayload{
		Model:       config.Model,
		Messages:    messages,
		Temperature: config.Temperature,
		Stream:      true,
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	req.Header.Add("Authorization", "Bearer "+config.APIKey)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	var buffer bytes.Buffer

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil && err == io.EOF {
			break
		}

		if len(line) > 6 {
			var chunk ChatCompletionChunk
			if string(line) == "data: [DONE]\n" {
				close(apiOutput)
				break
			} else if err := json.Unmarshal(line[6:], &chunk); err == nil {
				buffer.WriteString(chunk.Choices[0].Delta["content"])
				apiOutput <- buffer.String()
				buffer.Reset()
			}
		}
	}
}
