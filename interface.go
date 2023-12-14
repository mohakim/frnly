package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	chatbotOutput *tview.TextView
	userInput     *tview.TextArea
	app           *tview.Application
  typingCtx     context.Context
  cancelTyping  context.CancelFunc
)

func setupUI() {

  typingCtx, cancelTyping = context.WithCancel(context.Background())
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	app = tview.NewApplication()
	userInput = tview.NewTextArea()
	chatbotOutput = tview.NewTextView()

	userInput.SetBorder(true).SetBorderColor(tcell.ColorGreen)
	userInput.SetBackgroundColor(tcell.ColorDefault)
	userInput.SetTextStyle(tcell.StyleDefault)

	
	chatbotOutput.SetBorder(true).SetBorderColor(tcell.ColorGreen)
  chatbotOutput.SetBackgroundColor(tcell.ColorDefault)
	chatbotOutput.SetDynamicColors(true)
  chatbotOutput.SetDisabled(true)

	userInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
		  currentText := userInput.GetText()
      sf.Reset()
      if handleCmd(currentText) {
        return nil
      }
    }
		return event
	})

  app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    if event.Key() == tcell.KeyBacktab {
      if userInput.HasFocus() {
        app.SetFocus(chatbotOutput)
      } else {
        app.SetFocus(userInput)
      }
      return nil
    }
    return event
  })  

	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(userInput, 0, 1, true).
		AddItem(chatbotOutput, 0, 5, false)

	app.SetRoot(flex, true).EnableMouse(true)

	if err := app.Run(); err != nil {
		panic(err)
	}
}

func handleCmd(input string) bool {
  if strings.Contains(input, config.SubmitCmd) {
    input = strings.Replace(input, "!fin", "", -1)
    processInput(input)
    userInput.SetText("", true)
  } else if strings.Contains(input, config.ClearCmd) {
    cancelTyping()
    time.Sleep(time.Second * 1)
    typingCtx, cancelTyping = context.WithCancel(context.Background())
    chatbotOutput.SetText("")
    userInput.SetText("", true)
  } else if strings.Contains(input, config.PermCmd) {
    input = strings.Replace(input, "!perm", "", -1)
    session.Permanent = input
    userInput.SetText("", true)
  } else if strings.Contains(input, config.ResetCmd) {
    session.Dynamic = []ChatMessage{}
    userInput.SetText("", true)
  } else if strings.Contains(input, config.ExitCmd) {
    app.Stop()
  } else {
    return false
  }
  return true
}

func writeChatbotMessage(char rune, color string) {
	coloredChar := fmt.Sprintf("[%s]%s[white]", color, string(char))
	chatbotOutput.Write([]byte(string(coloredChar)))
	app.Draw()
	time.Sleep(time.Millisecond * 10)
}

func processInput(userInput string) {
  var apiOutput = make(chan string, 10000)

	session.Dynamic = append(session.Dynamic, ChatMessage{
		Role:    "user",
		Content: userInput,
	})

	go streamCompletion(config, session, apiOutput)
	go typeResponse(apiOutput)
}

func typeResponse(apiOutput chan string) {
  var response strings.Builder

  for token := range apiOutput {
    select {
      case <-typingCtx.Done():
        return
      default:
        response.WriteString(token)
        for _, char := range token {
          sf.Print(char)
      }
    }
  }
  
  sf.Print('\n')
  sf.Print('\n')
  sf.Print('\n')
  sf.Print('\n')

  if config.Context > 0 {
    session.Dynamic = append(session.Dynamic, ChatMessage{
      Role:    "assistant",
      Content: response.String(),
    })
  } else {
    session.Dynamic = []ChatMessage{}
  }
  updateSession()
}
