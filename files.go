package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const configDirName = ".config/frnly"

var (
	configDir    string
	historyPath  string
	settingsPath string
	sessionPath  string
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
	sessionPath = filepath.Join(configDir, "session.log")

	defaultSettings := `# settings.conf
# OpenAI API Configuration
API_KEY=""

# GPT Model and Tuning
Temperature=0.2
Model="gpt-4-0314"
Session=True
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

func readHistory(fileName string) (string, error) {
	fileData, err := os.ReadFile(fileName)

	if err != nil {
		return "", fmt.Errorf("Failed to read the history at %s: %w", fileName, err)
	}

	if len(fileData) == 0 {
		return "", fmt.Errorf("History is empty")
	} else {
		return string(fileData), nil
	}
}

func writeHistory(fileName string, history string) error {
	err := os.WriteFile(fileName, []byte(history), 0644)

	if err != nil {
		return fmt.Errorf("Failed to write to file: %w", err)
	}

	return nil
}

func readSession(fileName string) (Session, error) {
	fileData, err := os.ReadFile(fileName)

	if err != nil {
		return Session{}, fmt.Errorf("Failed to open file: %w", err)
	}

	var session Session

	if err := json.Unmarshal(fileData, &session); err != nil {
		return Session{}, fmt.Errorf("Failed to unmarshal session: %w", err)
	}

	return session, nil
}

func writeSession(filename string, session Session) error {
	content, err := json.MarshalIndent(session, "", " ")

	if err != nil {
		return fmt.Errorf("Failed to marshal session: %w", err)
	}

	err = os.WriteFile(filename, content, 0644)

	if err != nil {
		return fmt.Errorf("Failed to write session to file: %w", err)
	}

	return nil
}
