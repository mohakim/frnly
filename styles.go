package main

import (
	"errors"
	"fmt"
  "regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

type StatefulFormatter struct {
	stateStack  []string
	symbolCount map[string]int
	mu          sync.Mutex
  charBuffer  strings.Builder
  lang        strings.Builder
  readLang    bool
	ColorMap    map[string]func(string) (string, error)
}

func NewStatefulFormatter() *StatefulFormatter {
	return &StatefulFormatter{
		stateStack:  make([]string, 0),
		symbolCount: make(map[string]int),
	}
}

func (sf *StatefulFormatter) applyFormatting(ch rune) string {
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
        result.WriteString(sf.flushCharBuffer())
        sf.stateStack = updateStateStack(sf.stateStack, "isCommentSingle")
        stateChange = true
      } else if delimeter(buffer, commentMap[lang][2:]) {
        if state == "isCommentMulti" {
          result.WriteString(sf.flushCharBuffer())
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
          result.WriteString(sf.flushCharBuffer())
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
      result.WriteString(sf.flushCharBuffer())
    }
	  sf.charBuffer.Reset()
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
      result.WriteString(sf.flushCharBuffer())
    }
    if sf.readLang {
      sf.lang.WriteRune(ch)
    } else {
      formatted, _ := sf.ColorMap[getState(sf.stateStack)](string(ch))
      result.WriteString(formatted)
    }
    sf.charBuffer.Reset()
	}

	return result.String()
}

func (sf *StatefulFormatter) flushCharBuffer() string {
	state := getState(sf.stateStack)
	colorFunc := sf.ColorMap[state]
	coloredStr, _ := colorFunc(sf.charBuffer.String())
  sf.charBuffer.Reset()
	return coloredStr
}

func getState(stateStack []string) string {
	if len(stateStack) == 0 {
		return "Default"
	}
	return stateStack[len(stateStack)-1]
}

func getTerminalWidth() int {
	var ws struct {
		rows    uint16
		cols    uint16
		xpixels uint16
		ypixels uint16
	}

	wsPtr := uintptr(unsafe.Pointer(&ws))

	ret, _, _ := syscall.Syscall(syscall.SYS_IOCTL, uintptr(syscall.Stdout), uintptr(syscall.TIOCGWINSZ), wsPtr)
	if int(ret) == -1 {
		return 100
	}

	return int(ws.cols)
}

func updateStateStack(stack []string, state string) []string {
	if len(stack) > 0 && stack[len(stack)-1] == state {
		return stack[:len(stack)-1]
	}
	return append(stack, state)
}

func hexToANSI(hex string) (string, error) {
  if hex == "#GGGGGG" {
    botColor, _ := hexToANSI(config.BotColor)
    return fmt.Sprintf("%s\033[1m", botColor), nil
  }
	hex = strings.TrimPrefix(hex, "#")
	r, err := strconv.ParseInt(hex[0:2], 16, 64)
	g, err2 := strconv.ParseInt(hex[2:4], 16, 64)
	b, err3 := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil || err2 != nil || err3 != nil {
		return "", errors.New("Invalid HEX")
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", r, g, b), nil
}

func initializeColors(sf *StatefulFormatter, settings Settings) error {
	sf.ColorMap = make(map[string]func(string) (string, error))
	colorFields := []struct {
		Name string
		Hex  string
	}{
		{"Default", settings.BotColor},
		{"isCode", settings.CodeBlock},
		{"isBold", "#GGGGGG"},
		{"isCommentSingle", settings.Comments},
		{"isCommentMulti", settings.Comments},
		{"isTextBlock", settings.TextBlock},
		{"isReference", settings.References},
	}
	for _, color := range colorFields {
		sf.ColorMap[color.Name] = colorText(color.Hex, settings.BotColor)
	}
  sf.lang.WriteString("python")
	return nil
}

func colorText(colorHex, resetColor string) func(string) (string, error) {
	return func(text string) (string, error) {
		ansiColor, err := hexToANSI(colorHex)
		if err != nil {
			return "", errors.New("Hex to ANSI conversion failed")
		}
		return fmt.Sprintf("%s%s\033[0m", ansiColor, text), nil
	}
}
