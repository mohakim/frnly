package main

import (
	"log"
)

var (
  config Settings
  session Session
  sf *StatefulFormatter
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

	initializeColors(sf, config)

	if config.APIKey == "" {
		log.Fatal("You need to add your API key in ~/.config/frnly/settings.conf")
	}

	if config.Session {
		if session, err = readSession(); err != nil {
			log.Panic(err)
		}
	}
}

func main() {
  setupUI()
}
