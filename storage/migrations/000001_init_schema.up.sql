DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'currency') THEN
        CREATE TYPE "currency" AS ENUM ('USD', 'EURO', 'PLN');
    END IF;
END $$;


CREATE TABLE IF NOT EXISTS "account" (
    "id" bigserial PRIMARY KEY,
    "owner" varchar NOT NULL,
    "currency" currency NOT NULL,
    "balance" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "entry" (
    "id" bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE IF NOT EXISTS "transfer" (
    "id" bigserial PRIMARY KEY,
    "account_from" bigint NOT NULL,
    "account_to" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE INDEX IF NOT EXISTS "idx_account_owner " ON "account" ("owner");
CREATE INDEX IF NOT EXISTS idx_entry_account_id ON "entry" ("account_id");
CREATE INDEX IF NOT EXISTS idx_transfer_account_from ON "transfer" ("account_from");
CREATE INDEX IF NOT EXISTS idx_transfer_account_to ON "transfer" ("account_to");
CREATE INDEX IF NOT EXISTS idx_transfer_account_from_to ON "transfer" ("account_from", "account_to");


COMMENT ON COLUMN "account"."balance" IS 'should be positive';
COMMENT ON COLUMN "transfer"."amount" IS 'should be positive';

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'entry_account_id_fkey') THEN
        ALTER TABLE "entry" ADD CONSTRAINT entry_account_id_fkey FOREIGN KEY ("account_id") REFERENCES "account" ("id");
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'transfer_account_from_fkey') THEN
        ALTER TABLE "transfer" ADD CONSTRAINT transfer_account_from_fkey FOREIGN KEY ("account_from") REFERENCES "account" ("id");
    END IF;
END $$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_constraint WHERE conname = 'transfer_account_to_fkey') THEN
        ALTER TABLE "transfer" ADD CONSTRAINT transfer_account_to_fkey FOREIGN KEY ("account_to") REFERENCES "account" ("id");
    END IF;
END $$;
