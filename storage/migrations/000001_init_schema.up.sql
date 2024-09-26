CREATE TABLE IF NOT EXISTS "accounts" (
    "id" bigserial PRIMARY KEY,
    "owner" varchar NOT NULL,
    "currency" char(3) NOT NULL,
    "balance" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "entries" (
    "id" bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "transfers" (
    "id" bigserial PRIMARY KEY,
    "account_from" bigint NOT NULL,
    "account_to" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX IF NOT EXISTS "idx_accounts_owner" ON "accounts" ("owner");
CREATE INDEX IF NOT EXISTS idx_entries_account_id ON "entries" ("account_id");
CREATE INDEX IF NOT EXISTS idx_transfers_account_from ON "transfers" ("account_from");
CREATE INDEX IF NOT EXISTS idx_transfers_account_to ON "transfers" ("account_to");
CREATE INDEX IF NOT EXISTS idx_transfers_account_from_to ON "transfers" ("account_from", "account_to");


COMMENT ON COLUMN "accounts"."balance" IS 'should be positive';
COMMENT ON COLUMN "transfers"."amount" IS 'should be positive';

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'entries_account_id_fkey') THEN
        ALTER TABLE "entries"
        ADD CONSTRAINT entries_account_id_fkey FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'transfers_account_from_fkey') THEN
        ALTER TABLE "transfers"
        ADD CONSTRAINT transfers_account_from_fkey FOREIGN KEY ("account_from") REFERENCES "accounts" ("id") ON DELETE CASCADE;
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'transfers_account_to_fkey') THEN
        ALTER TABLE "transfers"
        ADD CONSTRAINT transfers_account_to_fkey FOREIGN KEY ("account_to") REFERENCES "accounts" ("id") ON DELETE CASCADE;
    END IF;
END $$;
