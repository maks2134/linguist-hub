package language

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed languages.json
var catalogJSON []byte

type catalogFile struct {
	UILocales        []string          `json:"uiLocales"`
	UILocaleToNative map[string]string `json:"uiLocaleToNative"`
	Languages        []languageEntry   `json:"languages"`
}

type languageEntry struct {
	Code       string   `json:"code"`
	PromptName string   `json:"promptName"`
	Flag       string   `json:"flag"`
	Roles      []string `json:"roles"`
	Endonym    string   `json:"endonym"`
	LocaleTags []string `json:"localeTags"`
}

var (
	byCode      map[LangCode]languageEntry
	aliasToCode map[string]LangCode
	uiLocales   map[UILocale]struct{}
	uiToNative  map[UILocale]LangCode
)

func init() {
	var raw catalogFile
	if err := json.Unmarshal(catalogJSON, &raw); err != nil {
		panic(fmt.Sprintf("language: invalid catalog: %v", err))
	}

	byCode = make(map[LangCode]languageEntry, len(raw.Languages))
	aliasToCode = make(map[string]LangCode)

	for _, lang := range raw.Languages {
		code := LangCode(lang.Code)
		byCode[code] = lang
		registerAlias(string(code), code)
		registerAlias(lang.PromptName, code)
		registerAlias(lang.Endonym, code)
	}

	uiLocales = make(map[UILocale]struct{}, len(raw.UILocales))
	for _, loc := range raw.UILocales {
		uiLocales[UILocale(loc)] = struct{}{}
	}

	uiToNative = make(map[UILocale]LangCode, len(raw.UILocaleToNative))
	for loc, native := range raw.UILocaleToNative {
		uiToNative[UILocale(loc)] = LangCode(native)
	}
}

func registerAlias(alias string, code LangCode) {
	key := normalizeKey(alias)
	if key == "" {
		return
	}
	aliasToCode[key] = code
}

func normalizeKey(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func entryFor(code LangCode) (languageEntry, bool) {
	e, ok := byCode[code]
	return e, ok
}

func (c LangCode) PromptName() string {
	if e, ok := entryFor(c); ok {
		return e.PromptName
	}
	return string(c)
}

func IsLearningTarget(c LangCode) bool {
	e, ok := entryFor(c)
	if !ok {
		return false
	}
	for _, r := range e.Roles {
		if Role(r) == RoleLearning {
			return true
		}
	}
	return false
}

func IsNative(c LangCode) bool {
	e, ok := entryFor(c)
	if !ok {
		return false
	}
	for _, r := range e.Roles {
		if Role(r) == RoleNative {
			return true
		}
	}
	return false
}

func UILocales() []UILocale {
	out := make([]UILocale, 0, len(uiLocales))
	for loc := range uiLocales {
		out = append(out, loc)
	}
	return out
}
