package repository

import (
	"broke-bank/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type TransactionRepository struct {
	Pg *sqlx.DB
}

func (tr *TransactionRepository) GetTransaction(transaction_id string) (*model.Transaction, error) {
	transaction := new(model.Transaction)
	err := tr.Pg.Get(
		transaction,
		`SELECT * FROM "transaction" tx WHERE tx.id = $1`,
		transaction_id,
	)

	return transaction, err
}

func (tr *TransactionRepository) DepositTransaction(transaction_id uuid.UUID, to_account_id string, amount decimal.Decimal) error {
	tx, err := tr.Pg.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`); err != nil {
		return err
	}

	if _, err = tx.Exec(`SELECT acc.balance FROM "account" acc WHERE acc.id = $1 FOR UPDATE`, to_account_id); err != nil {
		return err
	}

	if _, err = tx.Exec(`UPDATE "account" SET balance = balance + $2 WHERE id = $1`, to_account_id, amount); err != nil {
		return err
	}

	if _, err = tx.Exec(`INSERT INTO "transaction" (id, type, to_account_id, amount) VALUES ($1, 'deposit', $2, $3)`, transaction_id, to_account_id, amount); err != nil {
		return err
	}

	err = tx.Commit()

	return err
}