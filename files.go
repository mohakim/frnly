package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	historyPath  string
	settingsPath string
	sessionPath  string
)

func initializeConfig() error {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return fmt.Errorf("Couldn't determine the user's home directory: %w", err)
	}

  configDir := filepath.Join(homeDir, ".config/frnly")

	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatal("Couldn't create the configuration directory: ", err)
	}

	historyPath = filepath.Join(configDir, "history.log")
	settingsPath = filepath.Join(configDir, "settings.conf")
	sessionPath = filepath.Join(configDir, "session.log")

	defaultSettings := `# settings.conf
# OpenAI API Configuration
API_KEY=""

# GPT Model and Tuning
Temperature=0.5
Model="gpt-4-0314"
Session=False
Context=8192

# Styling
UserColor="#55DD55"
BotColor="#A6E3A1"
CodeBlock="#7700AA"
TextBlock="#94E2D5"
Comments="#BAC2DE"
References="#FFAA00"

# Interaction
Prompt=">>>"
Clear="!clear"
Submit="!fin"
History="!hist"
Exit="!exit"`

	initialSession := Session{
		Permanent: "",
		Dynamic:   []ChatMessage{},
	}

	initialSessionJSON, err := json.Marshal(initialSession)

	if err != nil {
		return fmt.Errorf("Couldn't marshal initial session to JSON: %w", err)
	}

	err = createFile(sessionPath, string(initialSessionJSON))
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

func readHistory() (string, error) {
	fileData, err := os.ReadFile(historyPath)

	if err != nil {
		return "", fmt.Errorf("Failed to read the history at %s: %w", historyPath, err)
	}

	if len(fileData) == 0 {
		return "", fmt.Errorf("History is empty")
	} else {
		return string(fileData), nil
	}
}

func writeHistory() error {
	err := os.WriteFile(historyPath, []byte(history), 0644)

	if err != nil {
		return fmt.Errorf("Failed to write to file: %w", err)
	}

	return nil
}

func readSession() (Session, error) {
	fileData, err := os.ReadFile(sessionPath)

	if err != nil {
		return Session{}, fmt.Errorf("Failed to open file: %w", err)
	}

	if err := json.Unmarshal(fileData, &session); err != nil {
		return Session{}, fmt.Errorf("Failed to unmarshal session: %w", err)
	}

	return session, nil
}

func writeSession() error {
	content, err := json.MarshalIndent(session, "", " ")

	if err != nil {
		return fmt.Errorf("Failed to marshal session: %w", err)
	}

	err = os.WriteFile(sessionPath, content, 0644)

	if err != nil {
		return fmt.Errorf("Failed to write session to file: %w", err)
	}

	return nil
}
