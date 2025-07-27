package secure

import (
	"html"
	"html/template"
	"strings"
)

func escapeHTML(input string) template.HTML {
	escaped := html.EscapeString(input)
	return template.HTML(escaped)
}

func safeJSONString(input string) string {
	input = strings.ReplaceAll(input, "\\", "\\\\")
	input = strings.ReplaceAll(input, "\"", "\\\"")
	input = strings.ReplaceAll(input, "\n", "\\n")
	input = strings.ReplaceAll(input, "\r", "\\r")
	input = strings.ReplaceAll(input, "\t", "\\t")

	input = strings.ReplaceAll(input, "<", "\\u003c")
	input = strings.ReplaceAll(input, ">", "\\u003e")
	input = strings.ReplaceAll(input, "&", "\\u0026")

	return input
}

func ValidateTextInput(input string, maxLength int) (string, bool) {
	// Check length
	if len(input) > maxLength {
		return "", false
	}

	// Remove null bytes and other control characters except newlines and tabs
	cleanInput := ""
	for _, r := range input {
		if r == '\n' || r == '\r' || r == '\t' || (r >= 32 && r < 127) || r >= 160 {
			cleanInput += string(r)
		}
	}

	return cleanInput, true
}
