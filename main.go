package main

import (
	"fmt"
	"log"
)

var (
  config Settings
  session Session
  sf *StatefulFormatter
  history string
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
	fmt.Println("\033[H\033[2J")

	for {
    readInput()
	}
}
