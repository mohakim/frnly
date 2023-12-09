package main

import (
	"fmt"
	"log"
  "github.com/gdamore/tcell/v2"
)

var (
  config Settings
  session Session
  sf *StatefulFormatter
  history string
  ui *uiElements
)

func init() {
	var err error

	if err = initializeConfig(); err != nil {
		log.Fatal(err)
	}

	if config, err = readSettings(settingsPath); err != nil {
		log.Fatal("Error reading configuration file: ", err)
	}

	sf = NewStatefulFormatter()

	if err = initializeColors(sf, config); err != nil {
		log.Fatal(err)
	}

  ui = newUI()

  if err := ui.app.Run(); err != nil {
    log.Fatal("Couldn't load UI")
  }
  
	if config.APIKey == "" {
		log.Fatal("You need to add your API key in ~/.config/frnly/settings.conf")
	}

	if config.Session {
		if session, err = readSession(); err != nil {
			fmt.Printf("Error: %v. Will proceed without persistent sessions\n", err)
		}

		if history, err = readHistory(); err != nil {
      fmt.Printf("Error: %v", err)
    }
	}
}

func main() {
   // Initialize TUI
    ui := newUI()  // Assuming newUI() sets up the TUI as discussed earlier

    // Set up TUI input handling
    ui.inputField.SetDoneFunc(func(key tcell.Key) {
        if key == tcell.KeyEnter {
            input := ui.inputField.GetText()
            if input == "!fin" {
                ui.inputField.SetText("")
            } else {
                readInput(input, ui)
            }
        }
    })

    // Run the TUI application
    if err := ui.app.Run(); err != nil {
        log.Fatalf("Failed to run application: %v", err)
    }
}
