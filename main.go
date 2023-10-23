package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func HexToANSI(hex string) string {

	hex = strings.TrimPrefix(hex, "#")

	red, _ := strconv.ParseInt(hex[0:2], 16, 64)
	green, _ := strconv.ParseInt(hex[2:4], 16, 64)
	blue, _ := strconv.ParseInt(hex[4:6], 16, 64)

	return fmt.Sprintf("\033[38;2;%d;%d;%dm", red, green, blue)
}

func main() {

	var (
		permanentHistory string
		dynamicHistory   string
		userInput        string
		formattedReply   string
	)

	if err := InitializeConfigFiles(); err != nil {
		log.Fatal(err)
	}

	config, err := readSettingsFromFile(settingsPath)

	if err != nil {
		log.Fatal("Error reading configuration file: ", err)
	}

	InitializeColors(config)

	if config.APIKey == "" {
		log.Fatal("You need to add your API key in ~/.config/frnly/settings.conf")
	}

	if config.History {
		permanentHistory, dynamicHistory, err = readHistoryFromFile(historyPath)

		if err != nil {
			fmt.Printf("Failed to read from %s! Will proceed without persistent history\nError: %v", historyPath, err)
		}
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(HexToANSI(config.UserColor))
		fmt.Print(config.Prompt)

		for {
			line, err := reader.ReadString('\n')

			if err != nil {
				log.Fatal("Failed to read user input!", err)
			}

			line = strings.TrimSpace(line)

			if strings.Contains(line, config.SubmitCommand) || strings.Contains(line, config.ClearCommand) || strings.Contains(line, config.ExitCommand) {
				userInput += line
				fmt.Print("\n")
				break
			}

			userInput += line + "\n"
		}

		if strings.Contains(userInput, config.ClearCommand) {
			fmt.Print("\033[H\033[2J")
			userInput = ""
			continue
		}

		if strings.Contains(userInput, config.ExitCommand) {
			break
		}

		userInput = strings.ReplaceAll(userInput, config.SubmitCommand, "")

		if config.History {
			dynamicHistory += "user: " + userInput + "\n"
			reply, err := getAssistantReply(config, permanentHistory+dynamicHistory)

			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			formattedReply = applyFormatting(reply)
			dynamicHistory += "assistant: " + formattedReply + "\n"
		} else {
			reply, err := getAssistantReply(config, userInput)

			if err != nil {
				fmt.Println("Error: ", err)
				return
			}

			formattedReply = applyFormatting(reply)
		}

		fmt.Println(formattedReply + "\n")

		if config.History {
			if len(dynamicHistory) > config.Context {
				dynamicHistory = truncateDynamicHistory(dynamicHistory)
			}

			err := writeHistoryToFile(historyPath, permanentHistory, dynamicHistory)

			if err != nil {
				fmt.Println(err)
			}
		}

		userInput = ""
	}
}
