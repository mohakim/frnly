package main

import (
	"context"
	"strings"
	"time"
  "fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const chatbotTypingSpeed = time.Millisecond * 10

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
  chatbotOutput.SetBorder(true).SetBorderColor(tcell.ColorGreen)

  userInput.SetBackgroundColor(tcell.ColorDefault)
  userInput.SetTextStyle(tcell.StyleDefault)
  chatbotOutput.SetBackgroundColor(tcell.ColorDefault)
  chatbotOutput.SetTextColor(tcell.ColorGreen)
  chatbotOutput.SetDynamicColors(true)

  ctx, cancel := context.WithCancel(context.Background())

  userInput.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
    if event.Key() == tcell.KeyEnter {
      currentText := userInput.GetText()

      if strings.Contains(currentText, "!fin") {
        cancel()
        ctx, cancel = context.WithCancel(context.Background())
        processInput(userInput.GetText(), ctx)
        userInput.SetText("", true)
        return nil
      }
    }
    return event
  })

  flex := tview.NewFlex().SetDirection(tview.FlexRow).
    AddItem(userInput, 0, 1, true).
    AddItem(chatbotOutput, 0, 4, false)

  app.SetRoot(flex, true).SetFocus(userInput)

  if err := app.Run(); err != nil {
    panic(err)
  }
}

func writeChatbotMessage(ctx context.Context, char rune, color tcell.Color) {
  // Logic to change color of the current rune (TODO)
  coloredChar := fmt.Sprintf("[#%06x]%s[white]", color.Hex(), string(char))
  select {
  case <-ctx.Done():
    return
  default:
    chatbotOutput.ScrollToEnd()
    _, _, _, height := chatbotOutput.GetInnerRect()
    row, _ := chatbotOutput.GetScrollOffset()
    if row >= height {
      chatbotOutput.ScrollToBeginning()
    }
    _, _ = chatbotOutput.Write([]byte(string(coloredChar)))
    app.Draw()
    time.Sleep(chatbotTypingSpeed)
  }
}
