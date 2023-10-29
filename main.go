package main

import (
	"fmt"
	"log"
	"sync"
)

var (
  config Settings
  session Session
  formatter *StatefulFormatter
  history string
)

func init() {
	var err error

	if err = initializeConfig(); err != nil {
		log.Fatal(err)
	}

	if config, err = readSettings(settingsPath); err != nil {
		log.Fatal("Error reading configuration file: ", err)
	}

	formatter = NewStatefulFormatter()

	if err = initializeColors(formatter, config); err != nil {
		log.Fatal(err)
	}

	if config.APIKey == "" {
		log.Fatal("You need to add your API key in ~/.config/frnly/settings.conf")
	}

	if config.Session {
		if session, err = readSession(); err != nil {
			fmt.Printf("Error: %v. Will proceed without persistent sessions\n", err)
		}

		history, _ = readHistory()
	}
}

func main() {
	var wg sync.WaitGroup
	fmt.Println("\033[H\033[2J")

	for {
		apiOutput := make(chan string, 10000)
		historyChannel := make(chan ChatMessage, 2)

		userInput := readInput()

    if isCommand(&userInput) {
			handleCommand(&userInput)
			continue
		}
    
		wg.Add(1)
		processInput(&userInput, apiOutput, &wg, historyChannel)
		wg.Wait()

		if config.Session {
			updateSession()
			for msg := range historyChannel {
				updateHistory(msg.Role, msg.Content)
			}
		}
  }
}

func updateSession() {
	var sessionSize, messagesToRemove int
	for _, msg := range session.Dynamic {
		sessionSize += len(msg.Content)
	}
	excess := sessionSize - config.Context
	if excess > 0 {
		for i, msg := range session.Dynamic {
			excess -= len(msg.Content)
			if excess <= 0 {
				messagesToRemove = i + 1
				break
			}
		}
		if messagesToRemove < len(session.Dynamic) {
			session.Dynamic = session.Dynamic[messagesToRemove:]
		} else {
			session.Dynamic = []ChatMessage{}
		}
	}
	writeSession()
}

func updateHistory(role, content string) {
	history += fmt.Sprintf("%s: %s\n\n", role, content)
	writeHistory()
}
