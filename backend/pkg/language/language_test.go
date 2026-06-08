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

func TestLocaleTags_english(t *testing.T) {
	tags := LangCode("english").LocaleTags()
	if tags.Primary != "en-US" {
		t.Fatalf("Primary = %q, want en-US", tags.Primary)
	}
	if len(tags.Fallback) == 0 || tags.Fallback[0] != "en-GB" {
		t.Fatalf("Fallback = %v", tags.Fallback)
	}
}

func TestParseAcceptLocale(t *testing.T) {
	code, ok := ParseAcceptLocale("de-DE")
	if !ok || code != "german" {
		t.Fatalf("ParseAcceptLocale(de-DE) = %q, %v", code, ok)
	}
}

func TestValidateLocale(t *testing.T) {
	if !ValidateLocale("fr-FR", "french") {
		t.Fatal("fr-FR should be valid for french")
	}
	if ValidateLocale("de-DE", "french") {
		t.Fatal("de-DE should not be valid for french")
	}
}
