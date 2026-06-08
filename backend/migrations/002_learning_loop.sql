ALTER TABLE language_profiles ADD COLUMN IF NOT EXISTS weaknesses jsonb NOT NULL DEFAULT '[]';
