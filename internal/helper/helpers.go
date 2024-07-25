package helper

import (
	"html"
	"regexp"
	"strings"
)

// cleanDescription cleans the description text by unescaping HTML entities, removing HTML tags, and replacing newlines with spaces
func CleanDescription(description string) string {
	return strings.ReplaceAll(removeHTMLTags(html.UnescapeString(description)), "\n", " ")
}

// removeHTMLTags removes HTML tags from a string
func removeHTMLTags(input string) string {
	re := regexp.MustCompile(`<.*?>`)
	return re.ReplaceAllString(input, "")
}

// getStringValue safely dereferences a string pointer, returning an empty string if the pointer is nil.
func GetStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
