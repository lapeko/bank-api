CREATE TYPE "currency" AS ENUM (
  'USD',
  'EURO',
  'PLN'
);

CREATE TABLE "account" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "currency" zurrency NOT NULL,
  "balance" bigint NOT NULL,
  "created_at" timestamptz NOT NULL
);

CREATE TABLE "entry" (
  "id" bigserial PRIMARY KEY,
  "account_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL
);

CREATE TABLE "transfer" (
  "id" bigserial PRIMARY KEY,
  "account_from" bigint NOT NULL,
  "account_to" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz NOT NULL
);

CREATE INDEX ON "account" ("owner");

CREATE INDEX ON "entry" ("account_id");

CREATE INDEX ON "transfer" ("account_from");

CREATE INDEX ON "transfer" ("account_to");

CREATE INDEX ON "transfer" ("account_from", "account_to");

COMMENT ON COLUMN "account"."balance" IS 'should be positive';

COMMENT ON COLUMN "entry"."amount" IS 'should be positive';

COMMENT ON COLUMN "transfer"."amount" IS 'should be positive';

ALTER TABLE "entry" ADD FOREIGN KEY ("account_id") REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("account_from") REFERENCES "account" ("id");

ALTER TABLE "transfer" ADD FOREIGN KEY ("account_to") REFERENCES "account" ("id");
