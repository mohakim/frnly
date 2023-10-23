package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
  "strconv"
)

var config Settings

func HexToANSI(hex string) string {

	hex = strings.TrimPrefix(hex, "#")

	red, _ := strconv.ParseInt(hex[0:2], 16, 64)
	green, _ := strconv.ParseInt(hex[2:4], 16, 64)
	blue, _ := strconv.ParseInt(hex[4:6], 16, 64)

	return fmt.Sprintf("\033[38;2;%d;%d;%dm", red, green, blue)
}

func main() {

  InitializeConfigFiles()
  config = readSettingsFromFile(settingsPath)
  InitializeColors(config)

	if config.APIKey == "" {
		fmt.Println("API_KEY environment variable not set")
		os.Exit(1)
	}
	
  var permanentHistory, dynamicHistory string
	if config.History {
		permanentHistory, dynamicHistory = readHistoryFromFile(historyPath)
	}

	reader := bufio.NewReader(os.Stdin)
	var userInput string
  
	for {
    fmt.Print(HexToANSI(config.UserColor))
    fmt.Print(config.Prompt)
		for {
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)
			if strings.Contains(line, config.SubmitCommand) || strings.Contains(line, config.ClearCommand) || strings.Contains(line, config.ExitCommand) {
				userInput += line
        fmt.Print("\n")
				break
			}
			userInput += line + "\n"
		}
    fmt.Print("\033[0m")

		if strings.Contains(userInput, config.ClearCommand) {
			fmt.Print("\033[H\033[2J")
			userInput = ""
			continue
		}

		if strings.Contains(userInput, config.ExitCommand) {
			break
		}

		
    userInput = strings.ReplaceAll(userInput, config.SubmitCommand, "")
    if config.History{
      dynamicHistory += "user: " + userInput + "\n"
    }

		
		reply, err := getAssistantReply(config.APIKey, permanentHistory+dynamicHistory)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
    
    formattedReply := applyFormatting(reply)

    if config.History {
      dynamicHistory += "assistant: " + formattedReply + "\n"
	  }

		
		fmt.Println(formattedReply + "\n")

		if config.History {
			if len(dynamicHistory) > maxContext {
				dynamicHistory = truncateDynamicHistory(dynamicHistory)
			}
			writeHistoryToFile(historyPath, permanentHistory, dynamicHistory)
		}
		
		userInput = ""
	}
}
