package main

import (
	"regexp"
	"strings"
	"sync"
)

type StatefulFormatter struct {
	stateStack  []string
	mu          sync.Mutex
  charBuffer  strings.Builder
  lang        strings.Builder
  readLang    bool
	ColorMap    map[string]string
}

func NewStatefulFormatter() *StatefulFormatter {
	return &StatefulFormatter{
 		stateStack:  make([]string, 0),
	}
}

func (sf *StatefulFormatter) Reset() {
  sf.stateStack = []string{}
  sf.charBuffer.Reset()
  sf.lang.Reset()
}

func (sf *StatefulFormatter) Print(ch rune) {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	var (
    char        = string(ch)
    //stateChange bool
    buffer      = sf.charBuffer.String()
    lang        = sf.lang.String()
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
  case regexp.MustCompile(`[\p{P}\p{S}]`).MatchString(char) && char != "." && char != ",":
    sf.charBuffer.WriteString(char)
	case ch == '\n', ch == ' ':
    
		if ch == '\n' {
      if sf.getState() == "isCode" && sf.readLang {
        sf.readLang = false
      } else if sf.getState() == "isCommentSingle" {
        sf.stateStack = sf.updateStateStack("isCommentSingle")
        //stateChange = true
      } else if len(commentMap[lang]) > 2 && delimeter(buffer, commentMap[lang][2:]) {
        if sf.getState() == "isCommentMulti" {
          sf.flushBuffer(nil)
        }
        sf.stateStack = sf.updateStateStack("isCommentMulti")
      }

      if (buffer == "`" || buffer == "``") && !strings.Contains(sf.getState(), "isC") {
        sf.stateStack = sf.updateStateStack("isTextBlock")
        //stateChange = true
      } else if buffer == "```" {
        sf.charBuffer.Reset()
        sf.stateStack = sf.updateStateStack("isCode")
        //stateChange = true
      } else if strings.Contains(buffer, "*") && (sf.getState() == "Default" || sf.getState() == "isBold")  {
        sf.flushBuffer([]rune{'*'})
        sf.stateStack = sf.updateStateStack("isBold")
        //stateChange = true
      }
    } else {
      if len(commentMap[lang]) > 2 && delimeter(buffer, commentMap[lang][:2]) {
        sf.stateStack = sf.updateStateStack("isCommentSingle")
      } else if len(commentMap[lang]) > 2 && delimeter(buffer, commentMap[lang][2:]) {
        if sf.getState() == "isCommentMulti" {
          writeChatbotMessage(ch, sf.ColorMap[sf.getState()])
        }
        sf.stateStack = sf.updateStateStack("isCommentMulti")
      }

      if strings.Contains(buffer, "`") && (sf.getState() == "Default" || sf.getState() == "isReference") {
        sf.flushBuffer([]rune{'`'})
        sf.stateStack = sf.updateStateStack("isReference")
        //stateChange = true
      } else if strings.Contains(buffer, "*") && sf.getState() == "isBold"{
        sf.flushBuffer([]rune{'*'})
        sf.stateStack = sf.updateStateStack("isBold")
        //stateChange = true
      }
    }
    sf.flushBuffer(nil)

    writeChatbotMessage(ch, sf.ColorMap[sf.getState()])
	default: 
    switch {
    case buffer == "```":
      sf.stateStack = sf.updateStateStack("isCode")
      sf.lang.Reset()
      sf.readLang = true
      //stateChange = true
    case (buffer == "*" || buffer == "**") &&  (sf.getState() == "Default" || sf.getState() == "isBold") :
      sf.stateStack = sf.updateStateStack("isBold")
      //stateChange = true
    case strings.Contains(buffer, "`") &&  (sf.getState() == "Default" || sf.getState() == "isReference") :
      sf.flushBuffer([]rune{'`'})
      sf.stateStack = sf.updateStateStack("isReference")
      //stateChange = true
    }
    
    sf.flushBuffer([]rune{'`', '*'})

    if sf.readLang {
      sf.lang.WriteRune(ch)
    } else {
      writeChatbotMessage(ch, sf.ColorMap[sf.getState()])
    }
	}

}

func (sf *StatefulFormatter) flushBuffer(exclude []rune) {

  excludeMap := make(map[rune]bool)
  for _, r := range exclude {
    excludeMap[r] = true
  }

  for _, char := range sf.charBuffer.String() {
    if !excludeMap[char] {
      writeChatbotMessage(char, sf.ColorMap[sf.getState()])
    }
  }
  sf.charBuffer.Reset()
}

func (sf *StatefulFormatter) getState() string {
	if len(sf.stateStack) == 0 {
		return "Default"
	}
	return sf.stateStack[len(sf.stateStack)-1]
}

func (sf *StatefulFormatter) updateStateStack(state string) []string {
	if len(sf.stateStack) > 0 && sf.stateStack[len(sf.stateStack)-1] == state {
		return sf.stateStack[:len(sf.stateStack)-1]
	}
	return append(sf.stateStack, state)
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
