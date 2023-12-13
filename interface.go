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
)

func setupUI() {
	tview.Styles.PrimitiveBackgroundColor = tcell.ColorDefault
	app = tview.NewApplication()
	userInput = tview.NewTextArea()
	chatbotOutput = tview.NewTextView()

	userInput.SetBorder(true).SetBorderColor(tcell.ColorGreen)
	userInput.SetBackgroundColor(tcell.ColorDefault)
	userInput.SetTextStyle(tcell.StyleDefault)

	chatbotOutput.SetBackgroundColor(tcell.ColorDefault)
	chatbotOutput.SetDynamicColors(true)
  chatbotOutput.SetBorderPadding(1, 0, 1, 1)
  chatbotOutput.SetRegions(true)

	userInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
		  currentText := userInput.GetText()

      ctx, cancel := context.WithCancel(context.Background())
      if handleCmd(currentText, ctx, cancel) {
        return nil
      }
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

func handleCmd(input string, ctx context.Context, cancel context.CancelFunc) bool {
  if strings.Contains(input, config.SubmitCmd) {
    processInput(userInput.GetText(), ctx, cancel)
    userInput.SetText("", true)
  } else if strings.Contains(input, config.ClearCmd) {
    cancel()
    chatbotOutput.SetText("")
    userInput.SetText("", true)
  } else if strings.Contains(input, config.ExitCmd) {
    app.Stop()
  } else {
    return false
  }
  return true
}

func writeChatbotMessage(ctx context.Context, char rune, color string) {
	coloredChar := fmt.Sprintf("[%s]%s[white]", color, string(char))
	select {
	case <-ctx.Done():
		return
	default:
		chatbotOutput.Write([]byte(string(coloredChar)))
		app.Draw()
		time.Sleep(time.Millisecond * 10)
	}
}

func processInput(userInput string, ctx context.Context, cancel context.CancelFunc) {
  var apiOutput = make(chan string, 10000)

	session.Dynamic = append(session.Dynamic, ChatMessage{
		Role:    "user",
		Content: userInput,
	})

	go streamCompletion(config, session, apiOutput)
	go typeResponse(apiOutput, ctx, cancel)
}

func typeResponse(apiOutput chan string, ctx context.Context, cancel context.CancelFunc) {
  var response strings.Builder

	for token := range apiOutput {
		response.WriteString(token)
    for _, char := range token {
      sf.Print(char, ctx)
    }
  }

  sf.Print('\n', ctx)
  sf.Print('\n', ctx)
  cancel()

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
