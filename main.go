package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

var config Settings
var session Session
var formatter *StatefulFormatter
var history string

func init() {
	var err error
	if err = InitializeConfigFiles(); err != nil {
		log.Fatal(err)
	}
	if config, err = readSettings(settingsPath); err != nil {
		log.Fatal("Error reading configuration file: ", err)
	}

	formatter = NewStatefulFormatter()

	if err = InitializeColors(formatter, config); err != nil {
		log.Fatal(err)
	}
	if config.APIKey == "" {
		log.Fatal("You need to add your API key in ~/.config/frnly/settings.conf")
	}
	if config.Session {
		if session, err = readSession(sessionPath); err != nil {
			fmt.Printf("Error: %v. Will proceed without persistent sessions\n", err)
		}
		history, _ = readHistory(historyPath)
	}
}

func readInput() string {
	userColor, _ := HexToANSI(config.UserColor)
	reader := bufio.NewReader(os.Stdin)
	var sb strings.Builder

	fmt.Print(userColor + config.Prompt + " ")
	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if strings.Contains(line, config.SubmitCommand) || isSpecialCommand(line) {
			sb.WriteString(line)
			fmt.Print("\n")
			break
		}
		sb.WriteString(line + "\n")
	}
	return sb.String()
}

func isSpecialCommand(cmd string) bool {
	return strings.Contains(cmd, config.ClearCommand) || strings.Contains(cmd, config.ExitCommand) || strings.Contains(cmd, config.HistoryCommand)
}

func main() {
	var wg sync.WaitGroup
	fmt.Println("\033[H\033[2J")

	for {
		apiOutput := make(chan string)
		typingQueue := make(chan string)
		historyChannel := make(chan ChatMessage, 2)

		wg.Add(1)

		userInput := readInput()
		if isSpecialCommand(userInput) {
			handleSpecialCommands(&userInput, &wg)
			continue
		}

		processUserInput(&userInput, apiOutput, typingQueue, &wg, historyChannel)
		wg.Wait()

		if config.Session {
			updateSession()
			for msg := range historyChannel {
				updateHistory(msg.Role, msg.Content)
			}
		} else {
			session.Dynamic = session.Dynamic[:0]
		}
	}
}

func handleSpecialCommands(userInput *string, wg *sync.WaitGroup) {
	switch {
	case strings.Contains(*userInput, config.ClearCommand):
		fmt.Print("\033[H\033[2J")
		wg.Done()
	case strings.Contains(*userInput, config.HistoryCommand):
		fmt.Print("\033[H\033[2J")
		if history, err := readHistory(historyPath); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(history)
		}
		wg.Done()
	case strings.Contains(*userInput, config.ExitCommand):
		os.Exit(0)
	}
	*userInput = ""
}

func processUserInput(userInput *string, apiOutput chan string, typingQueue chan string, wg *sync.WaitGroup, historyChannel chan ChatMessage) {
	*userInput = strings.ReplaceAll(*userInput, config.SubmitCommand, "")
	userColor, _ := HexToANSI(config.UserColor)

	session.Dynamic = append(session.Dynamic, ChatMessage{
		Role:    "user",
		Content: *userInput,
	})

	historyChannel <- ChatMessage{
		Role:    "\033[0muser",
		Content: fmt.Sprintf("%s%s", userColor, *userInput),
	}

	go streamCompletion(config, session, apiOutput)
	go processAPIOutput(apiOutput, typingQueue)
	go typeResponse(typingQueue, wg, historyChannel)
}

func processAPIOutput(apiOutput chan string, typingQueue chan string) {
	for token := range apiOutput {
		typingQueue <- token
	}
	close(typingQueue)
}

func typeResponse(typingQueue chan string, wg *sync.WaitGroup, historyChannel chan ChatMessage) {
	var response, formattedResponse strings.Builder
	for text := range typingQueue {
		response.WriteString(text)
		for _, char := range text {
			formattedChar := formatter.ApplyFormatting(char)
			formattedResponse.WriteString(formattedChar)
			fmt.Print(formattedChar)
			time.Sleep(24 * time.Millisecond)
		}
	}
	fmt.Print("\n\n")
	session.Dynamic = append(session.Dynamic, ChatMessage{
		Role:    "assistant",
		Content: response.String(),
	})

	historyChannel <- ChatMessage{
		Role:    "\033[0massistant",
		Content: formattedResponse.String(),
	}
	close(historyChannel)
	response.Reset()
	formattedResponse.Reset()
	wg.Done()
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
	writeSession(sessionPath, session)
}

func updateHistory(role, content string) {
	history += fmt.Sprintf("%s: %s\n\n", role, content)
	writeHistory(historyPath, history)
}
