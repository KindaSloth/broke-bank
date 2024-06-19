CREATE TYPE transaction_type AS ENUM ('deposit', 'withdrawal', 'transfer');

CREATE TABLE "transaction" (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
    type transaction_type NOT NULL,
    from_account_id UUID,
    to_account_id UUID,
    date_issued TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    amount DECIMAL(15, 2) NOT NULL,

    CONSTRAINT fk_from_account FOREIGN KEY(from_account_id) REFERENCES "account"(id),
    CONSTRAINT fk_to_account FOREIGN KEY(to_account_id) REFERENCES "account"(id)
);