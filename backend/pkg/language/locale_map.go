package language

func SuggestNativeFromUILocale(loc UILocale) (LangCode, bool) {
	code, ok := uiToNative[loc]
	return code, ok
}
