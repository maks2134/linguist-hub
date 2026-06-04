package language

import (
	"errors"
	"strings"
)

var ErrInvalidLanguage = errors.New("invalid or unsupported language")

func ParseLangCode(s string) (LangCode, error) {
	code, ok := NormalizeLangCode(s)
	if !ok {
		return "", ErrInvalidLanguage
	}
	return code, nil
}

func NormalizeLangCode(s string) (LangCode, bool) {
	key := normalizeKey(s)
	if key == "" {
		return "", false
	}
	if code, ok := aliasToCode[key]; ok {
		return code, true
	}
	if _, ok := byCode[LangCode(key)]; ok {
		return LangCode(key), true
	}
	return "", false
}

func ParseUILocale(s string) (UILocale, bool) {
	loc := UILocale(strings.ToLower(strings.TrimSpace(s)))
	_, ok := uiLocales[loc]
	return loc, ok
}

func ParseLearningTarget(s string) (LangCode, error) {
	code, err := ParseLangCode(s)
	if err != nil {
		return "", err
	}
	if !IsLearningTarget(code) {
		return "", ErrInvalidLanguage
	}
	return code, nil
}

func ParseNative(s string) (LangCode, error) {
	code, err := ParseLangCode(s)
	if err != nil {
		return "", err
	}
	if !IsNative(code) {
		return "", ErrInvalidLanguage
	}
	return code, nil
}
