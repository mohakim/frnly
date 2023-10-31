package main

import (
	"errors"
	"fmt"
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

	var result strings.Builder
	char := string(ch)
	state := getState(sf.stateStack)
  stateChange := false
	switch ch {
	case '`', '*', '#', '/':
		sf.symbolCount[char]++
    sf.charBuffer.WriteString(char)
	case '\n', ' ':
		for symbol, count := range sf.symbolCount {
			if ch == '\n' {
				if symbol == "`" && (count == 1 || count == 2) {
					sf.stateStack = updateStateStack(sf.stateStack, "isTextBlock")
          stateChange = true
				} else if symbol == "/" && state == "isComment" {
					sf.stateStack = updateStateStack(sf.stateStack, "isComment")
          stateChange = true
				} else if symbol == "`" && count == 3 {
					sf.stateStack = updateStateStack(sf.stateStack, "isCode")
          stateChange = true
				} else if symbol == "*" && count > 0 {
					sf.stateStack = updateStateStack(sf.stateStack, "isBold")
          stateChange = true
        }
			} else if ch == ' ' {
				if symbol == "/" && count == 2 {
					sf.stateStack = updateStateStack(sf.stateStack, "isComment")
					//formatted, _ := sf.ColorMap["isComment"]("//")
					//result.WriteString(formatted)
				} else if symbol == "#" && count > 0 {
					sf.stateStack = updateStateStack(sf.stateStack, "isBold")
          stateChange = true
				} else if symbol == "`" && count == 1 {
					sf.stateStack = updateStateStack(sf.stateStack, "isReference")
          stateChange = true
				} else if symbol == "*" && count > 0 && getState(sf.stateStack) == "isBold" {
					sf.stateStack = updateStateStack(sf.stateStack, "isBold")
          stateChange = true
        }
			}
			sf.symbolCount[symbol] = 0
		}
    if !stateChange {
      result.WriteString(sf.flushCharBuffer())
    }
	  sf.charBuffer.Reset()
		result.WriteString(char)
	default:
		// Process accumulated symbols
		for symbol, count := range sf.symbolCount {
			switch {
			case symbol == "`" && count == 3:
				sf.stateStack = updateStateStack(sf.stateStack, "isCode")
        stateChange = true
			case symbol == "*" && count > 0 && state != "isCode":
				sf.stateStack = updateStateStack(sf.stateStack, "isBold")
        stateChange = true
			case symbol == "`" && count == 1 && state != "isCode":
				sf.stateStack = updateStateStack(sf.stateStack, "isReference")
        stateChange = true
			}
			sf.symbolCount[symbol] = 0
		}
    if !stateChange {
      result.WriteString(sf.flushCharBuffer())
    }
	  sf.charBuffer.Reset()
		formatted, _ := sf.ColorMap[getState(sf.stateStack)](string(ch))
		result.WriteString(formatted)
	}

	return result.String()
}

func (sf *StatefulFormatter) flushCharBuffer() string {
	state := getState(sf.stateStack)
	colorFunc := sf.ColorMap[state]
	coloredStr, _ := colorFunc(sf.charBuffer.String())
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
	red, err := strconv.ParseInt(hex[0:2], 16, 64)
	if err != nil {
		return "", err
	}
	green, err := strconv.ParseInt(hex[2:4], 16, 64)
	if err != nil {
		return "", err
	}
	blue, err := strconv.ParseInt(hex[4:6], 16, 64)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", red, green, blue), nil
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
		{"isComment", settings.Comments},
		{"isTextBlock", settings.TextBlock},
		{"isReference", settings.References},
	}
	for _, color := range colorFields {
		sf.ColorMap[color.Name] = colorText(color.Hex, settings.BotColor)
	}
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
