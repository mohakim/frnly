package main

import "fmt"

func updateSession() {
	var sessionSize, messagesToRemove int
	for _, msg := range session.Dynamic {
		sessionSize += len(msg.Content)
	}
	excess := sessionSize - config.Context
	if excess > 0 {
		for i, msg := range session.Dynamic {
			excess -= len(msg.Content)
			if excess <= 0 {
				messagesToRemove = i + 1
				break
			}
		}
		if messagesToRemove < len(session.Dynamic) {
			session.Dynamic = session.Dynamic[messagesToRemove:]
		} else {
			session.Dynamic = []ChatMessage{}
		}
	}
	writeSession()
}

func updateHistory(role, content string) {
	history += fmt.Sprintf("%s: %s\n\n", role, content)
	writeHistory()
}
