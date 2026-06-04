package language

import "testing"

func TestNormalizeLangCode_aliases(t *testing.T) {
	tests := []struct {
		in   string
		want LangCode
	}{
		{"English", "english"},
		{"english", "english"},
		{" ENGLISH ", "english"},
		{"Russian", "russian"},
		{"Русский", "russian"},
		{"Español", "spanish"},
	}
	for _, tt := range tests {
		got, ok := NormalizeLangCode(tt.in)
		if !ok || got != tt.want {
			t.Fatalf("NormalizeLangCode(%q) = %q, %v; want %q", tt.in, got, ok, tt.want)
		}
	}
}

func TestNormalizeLangCode_unknown(t *testing.T) {
	if _, ok := NormalizeLangCode("klingon"); ok {
		t.Fatal("expected false for unknown language")
	}
}

func TestSuggestNativeFromUILocale(t *testing.T) {
	tests := map[UILocale]LangCode{
		"en": "english",
		"ru": "russian",
		"es": "spanish",
		"ar": "arabic",
		"ja": "japanese",
	}
	for loc, want := range tests {
		got, ok := SuggestNativeFromUILocale(loc)
		if !ok || got != want {
			t.Fatalf("SuggestNativeFromUILocale(%q) = %q, %v; want %q", loc, got, ok, want)
		}
	}
}

func TestIsLearningTarget_koreanOnlyLearning(t *testing.T) {
	if !IsLearningTarget("korean") {
		t.Fatal("korean should be learning target")
	}
	if IsNative("korean") {
		t.Fatal("korean should not be native")
	}
}

func TestParseLearningTarget_rejectsInvalidNativeOnly(t *testing.T) {
	// korean is learning-only, should work as target
	if _, err := ParseLearningTarget("korean"); err != nil {
		t.Fatalf("korean as target: %v", err)
	}
}

func TestPromptName(t *testing.T) {
	if got := LangCode("english").PromptName(); got != "English" {
		t.Fatalf("PromptName() = %q", got)
	}
}
