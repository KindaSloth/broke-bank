package repository

import (
	"broke-bank/model"
	"broke-bank/utils"

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

func (tr *TransactionRepository) WithdrawalTransaction(transaction_id uuid.UUID, from_account_id string, amount decimal.Decimal) error {
	tx, err := tr.Pg.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`); err != nil {
		return err
	}

	if _, err = tx.Exec(`SELECT acc.balance FROM "account" acc WHERE acc.id = $1 FOR UPDATE`, from_account_id); err != nil {
		return err
	}

	if _, err = tx.Exec(`UPDATE "account" SET balance = balance - $2 WHERE id = $1`, from_account_id, amount); err != nil {
		return err
	}

	if _, err = tx.Exec(`INSERT INTO "transaction" (id, type, from_account_id, amount) VALUES ($1, 'withdrawal', $2, $3)`, transaction_id, from_account_id, amount); err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func (tr *TransactionRepository) TransferTransaction(transaction_id uuid.UUID, from_account_id string, to_account_id string, amount decimal.Decimal) error {
	tx, err := tr.Pg.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err = tx.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`); err != nil {
		return err
	}

	// Sort the UUIDs here before locking; this will ensure that the locks always happen in the same order to avoid deadlock issues.
	first_id_lock, second_id_lock := utils.SortStringUUIDs(from_account_id, to_account_id)
	if _, err = tx.Exec(`SELECT acc.balance FROM "account" acc WHERE acc.id = $1 FOR UPDATE`, first_id_lock); err != nil {
		return err
	}
	if _, err = tx.Exec(`SELECT acc.balance FROM "account" acc WHERE acc.id = $1 FOR UPDATE`, second_id_lock); err != nil {
		return err
	}

	if _, err = tx.Exec(`UPDATE "account" SET balance = balance - $2 WHERE id = $1`, from_account_id, amount); err != nil {
		return err
	}

	if _, err = tx.Exec(`UPDATE "account" SET balance = balance + $2 WHERE id = $1`, to_account_id, amount); err != nil {
		return err
	}

	if _, err = tx.Exec(`INSERT INTO "transaction" (id, type, from_account_id, to_account_id, amount) VALUES ($1, 'transfer', $2, $3, $4)`, transaction_id, from_account_id, to_account_id, amount); err != nil {
		return err
	}

	err = tx.Commit()

	return err
}
