DROP INDEX IF EXISTS idx_accounts_user_id;
CREATE INDEX IF NOT EXISTS idx_accounts_owner ON "accounts" ("owner");

ALTER TABLE "accounts"
    DROP CONSTRAINT IF EXISTS accounts_user_id_fkey;

ALTER TABLE "accounts"
    DROP CONSTRAINT IF EXISTS accounts_user_currency_unique;

COMMENT ON COLUMN "transfers"."amount" IS NULL;
COMMENT ON COLUMN "accounts"."balance" IS NULL;

ALTER TABLE "accounts"
    DROP COLUMN IF EXISTS "user_id",
    ADD COLUMN IF NOT EXISTS "owner" varchar NOT NULL;

DROP TABLE IF EXISTS "users";
