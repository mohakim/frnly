package main

import (
	"github.com/rivo/tview"
	"github.com/gdamore/tcell/v2"
)

// uiElements holds the UI components
type uiElements struct {
	app       *tview.Application
	inputField *tview.InputField
	chatArea  *tview.TextView
}

// newUI initializes the UI elements and returns the uiElements struct
func newUI() *uiElements {
	ui := &uiElements{
		app: tview.NewApplication(),
		inputField: tview.NewInputField().
			SetLabel("Input: ").
			SetLabelColor(tcell.ColorGreen).
			SetFieldBackgroundColor(tcell.ColorDefault).
			SetFieldTextColor(tcell.ColorGreen).
			SetBorder(true).
			SetBorderColor(tcell.ColorGreen).
			SetBorderAttributes(tcell.AttrBold).
			SetRounded(true),
		chatArea: tview.NewTextView().
			SetDynamicColors(true).
			SetRegions(true).
			SetScrollable(true).
			SetWordWrap(true).
			SetBorder(true).
			SetBorderColor(tcell.ColorGreen).
			SetBorderAttributes(tcell.AttrBold).
			SetRounded(true),
	}


  // Set input field behavior
  ui.inputField.SetDoneFunc(func(key tcell.Key) {
      if key == tcell.KeyEnter {
          input := ui.inputField.GetText()
          if input == "!fin" {
              ui.inputField.SetText("")
          } else {
              readInput(input, ui) // Replace with your actual input handling function
              ui.inputField.SetText("")
          }
      }
  })

	// Layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(ui.inputField, 1, 1, true).
		AddItem(ui.chatArea, 0, 4, false)

	ui.app.SetRoot(flex, true)

	return ui
}

// updateChatArea is an existing function adjusted for the new chatArea component
func (ui *uiElements) updateChatArea(message string) {
	ui.chatArea.SetText(message)
}

