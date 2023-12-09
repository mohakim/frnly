package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
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

	  ui.updateChatArea(line)
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
		processInput(userInput, &wg)
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

func processInput(userInput *string, wg *sync.WaitGroup) {
  var (
    apiOutput       = make(chan string, 10000)
    historyChannel  = make(chan ChatMessage, 2)
  )

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

  if config.Session {
    updateSession()
    for msg := range historyChannel {
      updateHistory(msg.Role, msg.Content)
    }
  }
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

    if getState(sf.stateStack) != "isCode" && getState(sf.stateStack) != "isComment" {
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
      formattedChar := sf.applyFormatting(char)
      time.Sleep(24 * time.Millisecond)
      fmt.Print(formattedChar)
    }
  }
	fmt.Print("\n")
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
