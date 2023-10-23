package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func HexToANSI(hex string) (string, error) {
	hex = strings.TrimPrefix(hex, "#")
	red, err := strconv.ParseInt(hex[0:2], 16, 64)

	if err != nil {
		return "", err
	}

	green, err := strconv.ParseInt(hex[2:4], 16, 64)

	if err != nil {
		return "", err
	}

	blue, err := strconv.ParseInt(hex[4:6], 16, 64)

	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\033[38;2;%d;%d;%dm", red, green, blue), nil
}

func readInput(config Settings) string {
	var userInput string
	userColor, err := HexToANSI(config.UserColor)

	if err != nil {
		fmt.Printf("Failed to process user color. Make sure it's in hex format! %v", err)
	}

	fmt.Print(userColor + config.Prompt)
	reader := bufio.NewReader(os.Stdin)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatal("Failed to read user input!", err)
		}

		line = strings.TrimSpace(line)

		if strings.Contains(line, config.SubmitCommand) || strings.Contains(line, config.ClearCommand) || strings.Contains(line, config.ExitCommand) || strings.Contains(line, config.HistoryCommand) {
			userInput += line
			fmt.Print("\n")
			break
		}

		userInput += line + "\n"
	}
	return userInput
}

func main() {
	var (
		session          Session
		sessionSize      int
		messagesToRemove int
		history          string
		formattedReply   string
	)

	if err := InitializeConfigFiles(); err != nil {
		log.Fatal(err)
	}

	config, err := readSettings(settingsPath)

	if err != nil {
		log.Fatal("Error reading configuration file: ", err)
	}

	InitializeColors(config)
  history, err = readHistory(historyPath)

  if err != nil {
    fmt.Printf("%v", err)
  }

	if config.APIKey == "" {
		log.Fatal("You need to add your API key in ~/.config/frnly/settings.conf")
	}

	if config.Session {
		session, err = readSession(sessionPath)
		if err != nil {
			fmt.Printf("Failed to read from %s! Will proceed without persistent sessions\nError: %v", sessionPath, err)
		}
	}

	for {
		userInput := readInput(config)

		if strings.Contains(userInput, config.ClearCommand) {
			fmt.Print("\033[H\033[2J")
			userInput = ""
			continue
		}

    if strings.Contains(userInput, config.HistoryCommand) {
      fmt.Print("\033[H\033[2J")
			userInput = ""
  	  botColor, err := HexToANSI(config.BotColor)

	    if err != nil {
		    fmt.Printf("Failed to process bot color. Make sure it's in hex format! %v", err)
	    }
 
      fmt.Println(botColor + history)
      continue
    }

		if strings.Contains(userInput, config.ExitCommand) {
			break
		}

		userInput = strings.ReplaceAll(userInput, config.SubmitCommand, "")

		if config.Session {

			session.Dynamic = append(session.Dynamic, ChatMessage{
				Role:    "user",
				Content: userInput,
			})

      history += "user: " + userInput + "\n"

			reply, err := getAssistantReply(config, session)

			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

      session.Dynamic = append(session.Dynamic, ChatMessage{
				Role:    "assistant",
				Content: userInput,
			})

			formattedReply = applyFormatting(reply)
			history += "assistant: " + formattedReply + "\n"
		} else {
			session.Dynamic = append(session.Dynamic, ChatMessage{
				Role:    "user",
				Content: userInput,
			})

			reply, err := getAssistantReply(config, session)

			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			formattedReply = applyFormatting(reply)
      session = Session{}
		}

		fmt.Println(formattedReply + "\n")

		if config.Session {
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
			}

			if messagesToRemove > 0 && messagesToRemove <= len(session.Dynamic) {
				session.Dynamic = session.Dynamic[messagesToRemove:]
			}

			err = writeSession(sessionPath, session)

			if err != nil {
				fmt.Println(err)
			}

			err = writeHistory(historyPath, history)

			if err != nil {
				fmt.Println(err)
			}
		}

		userInput = ""
	}
}
