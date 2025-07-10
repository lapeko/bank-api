DO $$
    BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'currency') THEN
            CREATE TYPE currency AS ENUM ('USD', 'EURO', 'PLN');
        END IF;
    END$$;

CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    full_name varchar NOT NULL,
    email varchar UNIQUE NOT NULL,
    hashed_password varchar NOT NULL,
    password_changes_at timestamp NOT NULL DEFAULT '0001-01-01 00:00:00Z',
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS accounts (
    id bigserial PRIMARY KEY,
    user_id bigint NOT NULL,
    currency currency NOT NULL,
    balance bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS entries (
    id bigserial PRIMARY KEY,
    account_id bigint NOT NULL,
    amount bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS transfers (
    id bigserial PRIMARY KEY,
    account_from bigint NOT NULL,
    account_to bigint NOT NULL,
    amount bigint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts(user_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_accounts_user_currency ON accounts(user_id, currency);
CREATE INDEX IF NOT EXISTS idx_entries_account_id ON entries(account_id);
CREATE INDEX IF NOT EXISTS idx_transfers_from ON transfers(account_from);
CREATE INDEX IF NOT EXISTS idx_transfers_to ON transfers(account_to);
CREATE INDEX IF NOT EXISTS idx_transfers_from_to ON transfers(account_from, account_to);

COMMENT ON COLUMN accounts.balance IS 'should be positive';
COMMENT ON COLUMN transfers.amount IS 'should be positive';

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.table_constraints
            WHERE constraint_name = 'fk_accounts_user_id'
        ) THEN
            ALTER TABLE accounts
                ADD CONSTRAINT fk_accounts_user_id
                    FOREIGN KEY (user_id) REFERENCES users(id);
        END IF;
    END$$;

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.table_constraints
            WHERE constraint_name = 'fk_entries_account_id'
        ) THEN
            ALTER TABLE entries
                ADD CONSTRAINT fk_entries_account_id
                    FOREIGN KEY (account_id) REFERENCES accounts(id);
        END IF;
    END$$;

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.table_constraints
            WHERE constraint_name = 'fk_transfers_from'
        ) THEN
            ALTER TABLE transfers
                ADD CONSTRAINT fk_transfers_from
                    FOREIGN KEY (account_from) REFERENCES accounts(id);
        END IF;
    END$$;

DO $$
    BEGIN
        IF NOT EXISTS (
            SELECT 1 FROM information_schema.table_constraints
            WHERE constraint_name = 'fk_transfers_to'
        ) THEN
            ALTER TABLE transfers
                ADD CONSTRAINT fk_transfers_to
                    FOREIGN KEY (account_to) REFERENCES accounts(id);
        END IF;
    END$$;
