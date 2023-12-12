package main

import (
	"context"
	"regexp"
	"strings"
	"sync"
)

type StatefulFormatter struct {
	stateStack  []string
	symbolCount map[string]int
	mu          sync.Mutex
  charBuffer  strings.Builder
  lang        strings.Builder
  readLang    bool
	ColorMap    map[string]string
}

func NewStatefulFormatter() *StatefulFormatter {
	return &StatefulFormatter{
 		stateStack:  make([]string, 0),
		symbolCount: make(map[string]int),
	}
}

func (sf *StatefulFormatter) Print(ch rune, ctx context.Context) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	var (
    result      strings.Builder
    char        = string(ch)
    state       = getState(sf.stateStack)
    buffer      = sf.charBuffer.String()
    lang        = sf.lang.String()
    stateChange bool
    delimeter   = func (s string, slice []string) bool {
      for _, item := range slice {
        if item == s && s != ""{
          return true
        }
      }
      return false
    }
  )

	switch {
  case regexp.MustCompile(`[\p{P}\p{S}]`).MatchString(char) && char != ".":
    if char == "*" || char == "`"{
      sf.symbolCount[char]++
    }
    sf.charBuffer.WriteString(char)
	case ch == '\n', ch == ' ':
		if ch == '\n' {
      if state == "isCode" && sf.readLang {
        sf.readLang = false
      } else if state == "isCommentSingle" {
        sf.stateStack = updateStateStack(sf.stateStack, "isCommentSingle")
        stateChange = true
      } else if delimeter(buffer, commentMap[lang][2:]) {
        if state == "isCommentMulti" {
          writeChatbotMessage(ctx, ch, sf.ColorMap[getState(sf.stateStack)])
          sf.charBuffer.Reset()
        }
        sf.stateStack = updateStateStack(sf.stateStack, "isCommentMulti")
      }
      for symbol, count := range sf.symbolCount {
        if symbol == "`" && (count == 1 || count == 2) && !strings.Contains(state, "isC") {
					sf.stateStack = updateStateStack(sf.stateStack, "isTextBlock")
          stateChange = true
				} else if symbol == "`" && count == 3 {
					sf.stateStack = updateStateStack(sf.stateStack, "isCode")
          stateChange = true
				} else if symbol == "*" && count > 0 && !strings.Contains(state, "isC") {
					sf.stateStack = updateStateStack(sf.stateStack, "isBold")
          stateChange = true
        }
      }
    } else if ch == ' ' {
      if delimeter(buffer, commentMap[lang][:2]) {
        sf.stateStack = updateStateStack(sf.stateStack, "isCommentSingle")
      } else if delimeter(buffer, commentMap[lang][2:]) {
        if state == "isCommentMulti" {
          writeChatbotMessage(ctx, ch, sf.ColorMap[getState(sf.stateStack)])
        }
        sf.stateStack = updateStateStack(sf.stateStack, "isCommentMulti")
      }
      for symbol, count := range sf.symbolCount {
        if symbol == "`" && count == 1 && !strings.Contains(state, "isC") {
          sf.stateStack = updateStateStack(sf.stateStack, "isReference")
          stateChange = true
        } else if symbol == "*" && count > 0 && state == "isBold"{
          sf.stateStack = updateStateStack(sf.stateStack, "isBold")
          stateChange = true
        }
      }
    }
    for k := range sf.symbolCount {
      delete(sf.symbolCount, k)
    }
       
    if !stateChange {
      for _, char := range sf.charBuffer.String() {
        writeChatbotMessage(ctx, char, sf.ColorMap[getState(sf.stateStack)])
      }
    }
	  sf.charBuffer.Reset()
    writeChatbotMessage(ctx, ch, sf.ColorMap[getState(sf.stateStack)])
		result.WriteString(char)
	default:
		for symbol, count := range sf.symbolCount {
			switch {
			case symbol == "`" && count == 3:
				sf.stateStack = updateStateStack(sf.stateStack, "isCode")
        sf.lang.Reset()
        sf.readLang = true
        stateChange = true
			case symbol == "*" && count > 0 && !strings.Contains(state, "isC"):
				sf.stateStack = updateStateStack(sf.stateStack, "isBold")
        stateChange = true
			case symbol == "`" && count == 1 && !strings.Contains(state, "isC"):
				sf.stateStack = updateStateStack(sf.stateStack, "isReference")
        stateChange = true
			}
			sf.symbolCount[symbol] = 0
		}
    if !stateChange {
      for _, char := range sf.charBuffer.String() {
        writeChatbotMessage(ctx, char, sf.ColorMap[getState(sf.stateStack)])
      }
      sf.charBuffer.Reset()
    }
    if sf.readLang {
      sf.lang.WriteRune(ch)
    } else {
      writeChatbotMessage(ctx, ch, sf.ColorMap[getState(sf.stateStack)])
    }
    sf.charBuffer.Reset()
	}

}

func getState(stateStack []string) string {
	if len(stateStack) == 0 {
		return "Default"
	}
	return stateStack[len(stateStack)-1]
}

func updateStateStack(stack []string, state string) []string {
	if len(stack) > 0 && stack[len(stack)-1] == state {
		return stack[:len(stack)-1]
	}
	return append(stack, state)
}

func initializeColors(sf *StatefulFormatter, settings Settings) {
  sf.ColorMap = map[string]string{
    "Default":           settings.BotColor,
    "isCode":            settings.CodeBlock,
    "isBold":            settings.BoldColor,
    "isCommentSingle":   settings.Comments,
    "isCommentMulti":    settings.Comments,
    "isTextBlock":       settings.TextBlock,
    "isReference":       settings.References,
  }
	sf.lang.WriteString("python")
}
