package main

import (
	"fmt"
	"os"
	"strings"
  "path/filepath"
)

const (
	configDirName = ".config/frnly"
	maxContext  = 8192
)

var (
  configDir     string
  historyPath   string
  settingsPath  string
)

func InitializeConfigFiles() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("Couldn't determine the user's home directory.")
	}

	configDir = filepath.Join(homeDir, configDirName)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		panic("Couldn't create the configuration directory.")
	}

	historyPath = filepath.Join(configDir, "history.log")
	settingsPath = filepath.Join(configDir, "settings.conf")

	defaultSettings := `# settings.conf
# OpenAI API Configuration
API_KEY=""

# GPT Model and Tuning
Temperature=0.2
Model="gpt-4-0314"
History=True
Context=8192

# Styling
UserColor="#4499FF"
BotColor="#00FF00"
CodeBlock="#FF00FF"
TextBlock="#FFFF00"
Comments="#666666"
References="#FFAA00"

# Interaction
Prompt=">>>"
clear="!clear"
submit="!fin"
exit="!exit"`

	createFileIfNotExist(historyPath, "")
	createFileIfNotExist(settingsPath, defaultSettings)
}

func createFileIfNotExist(filePath string, defaultContent string) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		if err := os.WriteFile(filePath, []byte(defaultContent), 0644); err != nil {
			panic(fmt.Sprintf("Failed to create the file %s.", filePath))
		}
	}
}

func readHistoryFromFile(fileName string) (string, string) {
	fileData, err := os.ReadFile(fileName)
	if err != nil {
		return "", ""
	}
  if len(fileData) == 0 {
    fmt.Println("The file is empty")
    return "", ""
  }
	sections := strings.Split(string(fileData), "Dynamic:")
	return strings.TrimSpace(sections[0]), strings.TrimSpace(sections[1])
}

func writeHistoryToFile(fileName string, permanentHistory string, dynamicHistory string) {
	content := permanentHistory + "\n\nDynamic:\n" + dynamicHistory
	err := os.WriteFile(fileName, []byte(content), 0644)
	if err != nil {
		fmt.Println("Failed to write to file:", err)
	}
}

func truncateDynamicHistory(dynamicHistory string) string {
	if len(dynamicHistory) > maxContext {
		lines := strings.Split(dynamicHistory, "\n")
		return strings.Join(lines[2:], "\n")
	}
	return dynamicHistory
}
