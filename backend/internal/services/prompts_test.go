package services

import (
	"strings"
	"testing"
)

func TestGetPlanGenerationPrompt_usesNativeContentLanguage(t *testing.T) {
	prompt := GetPlanGenerationPrompt("spanish", "russian", "B1")
	if !strings.Contains(prompt, "Russian") {
		t.Fatal("expected native prompt name Russian")
	}
	if !strings.Contains(prompt, "Spanish") {
		t.Fatal("expected target prompt name Spanish")
	}
	if !strings.Contains(prompt, "Write ALL module titles") {
		t.Fatal("expected instruction to write modules in native language")
	}
}

func TestGetOnboardingSystemPrompt_usesPromptNames(t *testing.T) {
	prompt := GetOnboardingSystemPrompt("english", "russian")
	if !strings.Contains(prompt, "Russian") || !strings.Contains(prompt, "English") {
		t.Fatalf("expected display names in prompt, got excerpt: %s", prompt[:120])
	}
}
