package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

func readInput() string {
  var sb strings.Builder
	userColor, _ := hexToANSI(config.UserColor)
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(userColor + config.Prompt + " ")

	for {
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if strings.Contains(line, config.SubmitCommand) || isCommand(&line) {
			sb.WriteString(line)
			fmt.Print("\n")
			break
		}

		sb.WriteString(line + "\n")
	}

	return sb.String()
}

func isCommand(userInput *string) bool {
  commands := []string{config.ClearCommand, config.ExitCommand, config.ExitCommand}

  for _, command := range commands {
    if strings.Contains(*userInput, command) {
      return true
    }
  }

  return false
}

func handleCommand(userInput *string) {
	switch {
  case strings.Contains(*userInput, config.SubmitCommand):    
    *userInput = strings.ReplaceAll(*userInput, config.SubmitCommand, "")
	case strings.Contains(*userInput, config.ClearCommand):
		fmt.Print("\033[H\033[2J")
	case strings.Contains(*userInput, config.HistoryCommand):
		fmt.Print("\033[H\033[2J")
		if history, err := readHistory(); err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(history)
		}
	case strings.Contains(*userInput, config.ExitCommand):
		os.Exit(0)
	}
	*userInput = ""
}

func processInput(userInput *string, apiOutput chan string, wg *sync.WaitGroup, historyChannel chan ChatMessage) {
	userColor, _ := hexToANSI(config.UserColor)

	session.Dynamic = append(session.Dynamic, ChatMessage{
		Role:    "user",
		Content: *userInput,
	})

	historyChannel <- ChatMessage{
		Role:    "\033[0muser",
		Content: fmt.Sprintf("%s%s", userColor, *userInput),
	}

	go streamCompletion(config, session, apiOutput)
	go typeResponse(apiOutput, wg, historyChannel)
}

func typeResponse(apiOutput chan string, wg *sync.WaitGroup, historyChannel chan ChatMessage) {
	var (
    response, formattedResponse strings.Builder
    maxWidth, lineLength int
	  skipSpace bool
  )

  maxWidth = getTerminalWidth()

	for token := range apiOutput {
		response.WriteString(token)

    shouldWrap := getCurrentState(formatter.stateStack) != "isCode" && getCurrentState(formatter.stateStack) != "isComment"
    if shouldWrap {
      wordLength := len(token)
      lineLength += wordLength
      if lineLength+wordLength > maxWidth {
        fmt.Print("\n")
        skipSpace = true
        lineLength = wordLength
      }
    }
    for _, char := range token {
      if skipSpace {
        skipSpace = false
        continue
      }
      if char == '\n' {
        lineLength = 0
      }
      if char == '\t' {
        char = ' '
        fmt.Print(string(char))
      }
      formattedChar := formatter.applyFormatting(char)
      time.Sleep(24 * time.Millisecond)
      fmt.Print(formattedChar)
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
