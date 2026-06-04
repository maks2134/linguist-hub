import catalog from "./catalog.json";

export type UILocale = (typeof catalog.uiLocales)[number];
export type LangCode = (typeof catalog.languages)[number]["code"];

export type LanguageEntry = {
  code: LangCode;
  promptName: string;
  flag: string;
  roles: ("learning" | "native")[];
  endonym: string;
};

export const UI_LOCALES = catalog.uiLocales as readonly UILocale[];

export const ALL_LANGUAGES = catalog.languages as LanguageEntry[];

export const LEARNING_LANGUAGES = ALL_LANGUAGES.filter((l) =>
  l.roles.includes("learning")
);

export const NATIVE_LANGUAGES = ALL_LANGUAGES.filter((l) =>
  l.roles.includes("native")
);

const aliasToCode = new Map<string, LangCode>();

function registerAlias(alias: string, code: LangCode) {
  const key = alias.trim().toLowerCase();
  if (key) aliasToCode.set(key, code);
}

for (const lang of ALL_LANGUAGES) {
  registerAlias(lang.code, lang.code);
  registerAlias(lang.promptName, lang.code);
  registerAlias(lang.endonym, lang.code);
}

export function normalizeLangCode(input: string): LangCode | null {
  const key = input.trim().toLowerCase();
  if (!key) return null;
  return aliasToCode.get(key) ?? (ALL_LANGUAGES.some((l) => l.code === key) ? (key as LangCode) : null);
}

export function findLanguage(code: string): LanguageEntry | undefined {
  const normalized = normalizeLangCode(code);
  if (!normalized) return undefined;
  return ALL_LANGUAGES.find((l) => l.code === normalized);
}

export function suggestNativeFromUILocale(locale: string): LangCode | undefined {
  const map = catalog.uiLocaleToNative as Record<string, string>;
  const mapped = map[locale];
  if (!mapped) return undefined;
  return mapped as LangCode;
}

export function displayLanguageName(code: string, uiLocale?: string): string {
  const entry = findLanguage(code);
  if (!entry) return code;
  return entry.endonym || entry.promptName;
}
