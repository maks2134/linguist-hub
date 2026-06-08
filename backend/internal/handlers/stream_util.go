package handlers

import (
	"regexp"
	"strings"
)

var sentenceEnd = regexp.MustCompile(`([.!?…]+[\s\n]*)`)

// SplitBySentences splits text into chunks ending at sentence boundaries for TTS streaming.
func SplitBySentences(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	parts := sentenceEnd.Split(text, -1)
	ends := sentenceEnd.FindAllString(text, -1)
	if len(ends) == 0 {
		return []string{text}
	}
	chunks := make([]string, 0, len(parts))
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if i < len(ends) {
			part += strings.TrimSpace(ends[i])
		}
		chunks = append(chunks, part)
	}
	if len(chunks) == 0 {
		return []string{text}
	}
	return chunks
}
