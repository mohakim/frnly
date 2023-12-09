package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/gdamore/tcell/v2"
	//"time"
)


func readInput() {
  var sb strings.Builder
	userColor, _ := hexToANSI(config.UserColor)
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(userColor + config.Prompt + " ")

	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if isCommand(&line) { 
			sb.WriteString(line)
			fmt.Print("\n")
      fullStr := sb.String()
      handleCommand(&fullStr)
			break
		}

		sb.WriteString(line + "\n")
	}	
}

func isCommand(userInput *string) bool {
  commands := []string{config.ClearCmd, config.ExitCmd, config.HistoryCmd, config.SubmitCmd}

  for _, command := range commands {
    if strings.Contains(*userInput, command) {
      return true
    }
  }

  return false
}

func handleCommand(userInput *string) {
	switch {
  case strings.Contains(*userInput, config.SubmitCmd):    
    *userInput = strings.ReplaceAll(*userInput, config.SubmitCmd, "")
    var wg sync.WaitGroup
    wg.Add(1)
		//processInput(userInput, &wg)
		wg.Wait()
	case strings.Contains(*userInput, config.ClearCmd):
		fmt.Print("\033[H\033[2J")
	case strings.Contains(*userInput, config.HistoryCmd):
		fmt.Print("\033[H\033[2J")
		if history, err := readHistory(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(history)
		}
	case strings.Contains(*userInput, config.ExitCmd):
		os.Exit(0)
	}
	*userInput = ""
}

func processInput(userInput string, ctx context.Context) {
  var (
    apiOutput       = make(chan string, 10000)
    historyChannel  = make(chan ChatMessage, 2)
  )

	userColor, _ := hexToANSI(config.UserColor)

	session.Dynamic = append(session.Dynamic, ChatMessage{
		Role:    "user",
		Content: userInput,
	})

	historyChannel <- ChatMessage{
		Role:    "\033[0muser",
		Content: fmt.Sprintf("%s%s", userColor, userInput),
	}

	go streamCompletion(config, session, apiOutput)
	go typeResponse(apiOutput, ctx)

  if config.Session {
    updateSession()
    for msg := range historyChannel {
      updateHistory(msg.Role, msg.Content)
    }
  }
}

func typeResponse(apiOutput chan string, ctx context.Context) {
	var (
    response strings.Builder
  )

	for token := range apiOutput {
		response.WriteString(token)
    for _, char := range token {
      writeChatbotMessage(ctx, char, tcell.ColorRed)
    }
  }

  writeChatbotMessage(ctx, '\n', tcell.Color100)

	session.Dynamic = append(session.Dynamic, ChatMessage{
		Role:    "assistant",
		Content: response.String(),
	})
}
