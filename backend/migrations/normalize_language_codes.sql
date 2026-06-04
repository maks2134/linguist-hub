-- Optional: normalize legacy profile language strings to canonical slugs.
-- Run manually against PostgreSQL after deploying pkg/language validation.

UPDATE language_profiles SET target_language = LOWER(target_language);
UPDATE language_profiles SET native_language = LOWER(native_language);

UPDATE language_profiles SET target_language = 'english' WHERE target_language IN ('english', 'en') OR target_language = 'English';
UPDATE language_profiles SET native_language = 'english' WHERE native_language IN ('english', 'en') OR native_language = 'English';
UPDATE language_profiles SET native_language = 'russian' WHERE native_language IN ('russian', 'ru') OR native_language = 'Russian';
-- Extend with other aliases as needed for your production data.
