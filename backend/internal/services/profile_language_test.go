package services

import (
	"testing"

	langpkg "github.com/makskozlov/linguist-hub/pkg/language"
)

func TestProfileLanguageNormalization_inputs(t *testing.T) {
	target, err := langpkg.ParseLearningTarget("English")
	if err != nil || target != "english" {
		t.Fatalf("target: %v %v", target, err)
	}
	native, err := langpkg.ParseNative("Russian")
	if err != nil || native != "russian" {
		t.Fatalf("native: %v %v", native, err)
	}
}
