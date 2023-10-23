package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const configDirName = ".config/frnly"

var (
	configDir    string
	historyPath  string
	settingsPath string
)

func InitializeConfigFiles() error {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return fmt.Errorf("Couldn't determine the user's home directory: %w", err)
	}

	configDir = filepath.Join(homeDir, configDirName)

	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal("Couldn't create the configuration directory: ", err)
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
Context=3000

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

	err = createFile(historyPath, "")
	err = createFile(settingsPath, defaultSettings)
	return err
}

func createFile(filePath string, defaultContent string) error {
	_, err := os.Stat(filePath)

	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return fmt.Errorf("Failed to stat the file %s: %w", filePath, err)
	}

	if err := os.WriteFile(filePath, []byte(defaultContent), 0644); err != nil {
		return fmt.Errorf("Failed to write to the file %s: %w", filePath, err)
	}

	return nil
}

func readHistoryFromFile(fileName string) (string, string, error) {
	fileData, err := os.ReadFile(fileName)

	if err != nil {
		return "", "", fmt.Errorf("Failed to read the history at %s: %w", fileName, err)
	}

	if len(fileData) == 0 {
		return "", "", nil
	}

	sections := strings.Split(string(fileData), "Dynamic:")
	return strings.TrimSpace(sections[0]), strings.TrimSpace(sections[1]), nil
}

func writeHistoryToFile(fileName string, permanentHistory string, dynamicHistory string) error {
	content := permanentHistory + "\n\nDynamic:\n" + dynamicHistory
	err := os.WriteFile(fileName, []byte(content), 0644)

	if err != nil {
		return fmt.Errorf("Failed to write to file: %w", err)
	}

	return nil
}

func truncateDynamicHistory(dynamicHistory string) string {
	lines := strings.Split(dynamicHistory, "\n")
	return strings.Join(lines[2:], "\n")
}
