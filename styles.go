package main

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/gookit/color"
)

var DynamicColorMap map[string]func(string, ...interface{}) string

func InitializeColors(settings Settings) {
	DynamicColorMap = make(map[string]func(string, ...interface{}) string)
	colorNames := []string{"BotColor", "CodeBlock", "TextBlock", "Comments", "References", "Prompt"}

	for _, colorName := range colorNames {
		colorHex := reflect.ValueOf(settings).FieldByName(colorName).String()
		DynamicColorMap[colorName] = func(text string, args ...interface{}) string {
			return color.HEX(colorHex, false).Basic().Render(text)
		}
	}
}

func applyFormatting(text string) string {
	lines := strings.Split(text, "\n")
	stateStack := make([]string, 0)
	formattedLines := make([]string, 0)
	boldPattern := regexp.MustCompile(`\*{1,3}[^\*]+\*{1,3}`)
	referencePattern := regexp.MustCompile("`([^`]+)`")

	for _, line := range lines {
		var currentState string

		if len(stateStack) > 0 {
			currentState = stateStack[len(stateStack)-1]
		}

		switch {
		case strings.HasPrefix(line, "```"):
			if currentState == "isCode" {
				stateStack = stateStack[:len(stateStack)-1]
			} else {
				stateStack = append(stateStack, "isCode")
			}
			continue
		case line == "`" || line == "``":
			if currentState == "isTextBlock" {
				stateStack = stateStack[:len(stateStack)-1]
			} else {
				stateStack = append(stateStack, "isTextBlock")
			}
			continue
		case strings.Contains(line, "// "):
			formattedLines = append(formattedLines, DynamicColorMap["Comments"](line))
			continue
		}

		if currentState == "" {
			line = boldPattern.ReplaceAllStringFunc(line, func(in string) string {
				trimmed := strings.Trim(in, "*")

				if strings.HasPrefix(trimmed, " ") || strings.HasSuffix(trimmed, " ") {
					return in
				}

				return color.Style{color.FgGreen, color.OpBold}.Sprint(trimmed)
			})
			line = referencePattern.ReplaceAllStringFunc(line, func(in string) string {
				return DynamicColorMap["References"](strings.Trim(in, "`"))
			})
		}

		switch currentState {
		case "isCode":
			formattedLines = append(formattedLines, DynamicColorMap["CodeBlock"](line))
		case "isTextBlock":
			formattedLines = append(formattedLines, DynamicColorMap["TextBlock"](line))
		default:
			formattedLines = append(formattedLines, DynamicColorMap["BotColor"](line))
		}
	}

	return strings.Join(formattedLines, "\n")
}
