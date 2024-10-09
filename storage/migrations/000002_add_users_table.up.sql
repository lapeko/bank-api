CREATE TABLE IF NOT EXISTS "users" (
    "id" bigserial PRIMARY KEY,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "hashed_password" varchar NOT NULL,
    "password_changes_at" timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "accounts"
    ADD COLUMN IF NOT EXISTS "user_id" bigint NOT NULL,
    DROP COLUMN IF EXISTS "owner";

COMMENT ON COLUMN "accounts"."balance" IS 'should be positive';
COMMENT ON COLUMN "transfers"."amount" IS 'should be positive';

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'accounts_user_currency_unique') THEN
        ALTER TABLE "accounts"
            ADD CONSTRAINT accounts_user_currency_unique UNIQUE ("user_id", "currency");
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'accounts_user_id_fkey') THEN
        ALTER TABLE "accounts"
            ADD CONSTRAINT accounts_user_id_fkey FOREIGN KEY ("user_id") REFERENCES "users" ("id");
    END IF;
END $$;

DROP INDEX IF EXISTS idx_accounts_owner;
CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON "accounts" ("user_id");
