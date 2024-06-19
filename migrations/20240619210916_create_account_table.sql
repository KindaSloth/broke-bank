CREATE TYPE account_status AS ENUM ('active', 'inactive');

CREATE TABLE "account" (
  id UUID PRIMARY KEY DEFAULT uuid_generate_v7(),
  user_id UUID,
  name VARCHAR(100) NOT NULL,
  balance DECIMAL(15, 2) NOT NULL,
  status account_status NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES "user"(id)
);