package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Settings struct {
	APIKey        string
	Temperature   float32
	Model         string
	Context       int
	UserColor     string
	BotColor      string
  BoldColor     string 
	CodeBlock     string
	TextBlock     string
	Comments      string
	References    string
	ClearCmd      string
	SubmitCmd     string
  ResetCmd      string
  PermCmd       string
	ExitCmd       string
}

func readSettings(filePath string) (Settings, error) {
	file, err := os.Open(filePath)

	if err != nil {
		return Settings{}, err
	}

	defer file.Close()

	settings := Settings{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		value = strings.Trim(value, "\"")

		switch key {
		case "API_KEY":
			settings.APIKey = value
		case "Temperature":
			settings.Temperature = parseFloat32(value)
		case "Model":
			settings.Model = value
		case "Context":
			settings.Context = parseInt(value)
		case "UserColor":
			settings.UserColor = value
		case "BotColor":
			settings.BotColor = value
    case "BoldColor":
      settings.BoldColor = value
		case "CodeBlock":
			settings.CodeBlock = value
		case "TextBlock":
			settings.TextBlock = value
		case "Comments":
			settings.Comments = value
		case "References":
			settings.References = value
		case "Clear":
			settings.ClearCmd = value
		case "Submit":
			settings.SubmitCmd = value
    case "Reset":
      settings.ResetCmd = value
    case "Permanent":
			settings.PermCmd = value
		case "Exit":
			settings.ExitCmd = value
		}
	}

	if err := scanner.Err(); err != nil {
		return Settings{}, err
	}

	return settings, nil
}

func parseFloat32(str string) float32 {
	var value float32
	fmt.Sscanf(str, "%f", &value)
	return value
}

func parseBool(str string) bool {
	return strings.ToLower(str) == "true"
}

func parseInt(str string) int {
	var value int
	fmt.Sscanf(str, "%d", &value)
	return value
}
