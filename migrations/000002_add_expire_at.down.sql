ALTER TABLE url DROP COLUMN IF EXISTS expire_at
;

DROP INDEX IF EXISTS idx_url_expire_at
;