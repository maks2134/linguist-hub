const ERROR_KEYS: Record<string, string> = {
  INVALID_LANGUAGE: "Errors.invalidLanguage",
  BAD_REQUEST: "Errors.badRequest",
  UNAUTHORIZED: "Errors.unauthorized",
  FORBIDDEN: "Errors.forbidden",
  NOT_FOUND: "Errors.notFound",
  INTERNAL_SERVER_ERROR: "Errors.internal",
  INVALID_CREDENTIALS: "Errors.invalidCredentials",
  EMAIL_EXISTS: "Errors.emailExists",
  LEVEL_REQUIRED: "Errors.levelRequired",
  ONBOARDING_COMPLETE: "Errors.onboardingComplete",
};

export function translationKeyForApiError(code?: string): string | undefined {
  if (!code) return undefined;
  const full = ERROR_KEYS[code];
  if (!full) return undefined;
  return full.replace("Errors.", "");
}

export function getLocaleFromPathname(): string {
  if (typeof window === "undefined") return "en";
  const match = window.location.pathname.match(/^\/([a-z]{2})(\/|$)/);
  return match?.[1] ?? "en";
}
