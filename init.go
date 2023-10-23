package main

import (
	"bufio"
  "fmt"
	"os"
	"strings"
)

type Settings struct {
	APIKey          string
	Temperature     float32
	Model           string
	History         bool
	Context         int
	UserColor       string
	BotColor        string
	CodeBlock       string
  TextBlock       string
	Comments        string
	References      string
	Prompt          string
	ClearCommand    string
	SubmitCommand   string
	ExitCommand     string
}

func readSettingsFromFile(filePath string) Settings {
	file, err := os.Open(filePath)
	if err != nil {
		panic("Failed to open settings.conf file")
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
		case "History":
			settings.History = parseBool(value)
		case "Context":
			settings.Context = parseInt(value)
		case "UserColor":
			settings.UserColor = value
		case "BotColor":
			settings.BotColor = value
		case "CodeBlock":
			settings.CodeBlock = value
    case "TextBlock":
			settings.TextBlock = value
		case "Comments":
			settings.Comments = value
		case "References":
			settings.References = value
		case "Prompt":
			settings.Prompt = value
		case "clear":
			settings.ClearCommand = value
		case "submit":
			settings.SubmitCommand = value
		case "exit":
			settings.ExitCommand = value
		}
	}

	if err := scanner.Err(); err != nil {
		panic("Failed to read settings.conf file")
	}

	return settings
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
