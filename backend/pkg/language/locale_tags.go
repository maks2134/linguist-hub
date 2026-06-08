package language

import "strings"

// LocaleTagSet holds BCP-47 tags for Web Speech API / TTS.
type LocaleTagSet struct {
	Primary  string   `json:"primary"`
	Fallback []string `json:"fallback,omitempty"`
}

// LocaleTags returns STT/TTS locale tags for a learning language code.
func (c LangCode) LocaleTags() LocaleTagSet {
	if e, ok := entryFor(c); ok && len(e.LocaleTags) > 0 {
		primary := e.LocaleTags[0]
		var fallback []string
		if len(e.LocaleTags) > 1 {
			fallback = append([]string(nil), e.LocaleTags[1:]...)
		}
		return LocaleTagSet{Primary: primary, Fallback: fallback}
	}
	return LocaleTagSet{Primary: defaultLocaleFor(c)}
}

// ParseAcceptLocale maps a BCP-47 tag to a LangCode if known.
func ParseAcceptLocale(tag string) (LangCode, bool) {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return "", false
	}
	lower := strings.ToLower(tag)
	for code, e := range byCode {
		for _, lt := range e.LocaleTags {
			if strings.EqualFold(lt, tag) || strings.HasPrefix(lower, strings.ToLower(lt)) {
				return code, true
			}
		}
		if strings.HasPrefix(lower, string(code)) {
			return code, true
		}
	}
	return "", false
}

// ValidateLocale checks that locale is supported for the given profile target language.
func ValidateLocale(locale string, target LangCode) bool {
	if locale == "" {
		return false
	}
	if parsed, ok := ParseAcceptLocale(locale); ok {
		return parsed == target
	}
	set := target.LocaleTags()
	if strings.EqualFold(locale, set.Primary) {
		return true
	}
	for _, f := range set.Fallback {
		if strings.EqualFold(locale, f) {
			return true
		}
	}
	return false
}

func defaultLocaleFor(c LangCode) string {
	defaults := map[LangCode]string{
		"english":    "en-US",
		"german":     "de-DE",
		"french":     "fr-FR",
		"spanish":    "es-ES",
		"italian":    "it-IT",
		"portuguese": "pt-PT",
		"russian":    "ru-RU",
		"chinese":    "zh-CN",
		"japanese":   "ja-JP",
		"korean":     "ko-KR",
		"arabic":     "ar-SA",
	}
	if loc, ok := defaults[c]; ok {
		return loc
	}
	return "en-US"
}
