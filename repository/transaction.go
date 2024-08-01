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

	if _, err := tx.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`); err != nil {
		return err
	}

	account_balance := new(AccountBalance)
	if err = tx.Get(account_balance, `SELECT acc.id, acc.balance FROM "account" acc WHERE acc.id = $1 FOR UPDATE`, to_account_id); err != nil {
		return err
	}

	if _, err = tx.Exec(`UPDATE "account" SET balance = $1 WHERE id = $2`, account_balance.Balance.Add(amount), to_account_id); err != nil {
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

	if _, err := tx.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`); err != nil {
		return err
	}

	account_balance := new(AccountBalance)
	if err = tx.Get(account_balance, `SELECT acc.id, acc.balance FROM "account" acc WHERE acc.id = $1 FOR UPDATE`, from_account_id); err != nil {
		return err
	}

	if _, err = tx.Exec(`UPDATE "account" SET balance = $1 WHERE id = $2`, account_balance.Balance.Sub(amount), from_account_id); err != nil {
		return err
	}

	if _, err = tx.Exec(`INSERT INTO "transaction" (id, type, from_account_id, amount) VALUES ($1, 'withdrawal', $2, $3)`, transaction_id, from_account_id, amount); err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

type AccountBalance struct {
	Id      uuid.UUID       `db:"id" json:"id"`
	Balance decimal.Decimal `db:"balance" json:"balance"`
}

func GetAccountBalance(first_account_balance *AccountBalance, second_account_balance *AccountBalance, account_id string) decimal.Decimal {
	if first_account_balance.Id.String() == account_id {
		return first_account_balance.Balance
	}

	return second_account_balance.Balance
}

func (tr *TransactionRepository) TransferTransaction(transaction_id uuid.UUID, from_account_id string, to_account_id string, amount decimal.Decimal) error {
	tx, err := tr.Pg.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`); err != nil {
		return err
	}

	// Sort the UUIDs here before locking; this will ensure that the locks always happen in the same order to avoid deadlock issues.
	first_id_lock, second_id_lock := utils.SortStringUUIDs(from_account_id, to_account_id)
	first_account_balance := new(AccountBalance)
	if err = tx.Get(first_account_balance, `SELECT acc.id, acc.balance FROM "account" acc WHERE acc.id = $1 FOR UPDATE`, first_id_lock); err != nil {
		return err
	}
	second_account_balance := new(AccountBalance)
	if err = tx.Get(second_account_balance, `SELECT acc.id, acc.balance FROM "account" acc WHERE acc.id = $1 FOR UPDATE`, second_id_lock); err != nil {
		return err
	}

	from_account_balance := GetAccountBalance(first_account_balance, second_account_balance, from_account_id)
	to_account_balance := GetAccountBalance(first_account_balance, second_account_balance, to_account_id)

	if _, err = tx.Exec(`UPDATE "account" SET balance = $1 WHERE id = $2`, from_account_balance.Sub(amount), from_account_id); err != nil {
		return err
	}

	if _, err = tx.Exec(`UPDATE "account" SET balance = $1 WHERE id = $2`, to_account_balance.Add(amount), to_account_id); err != nil {
		return err
	}

	if _, err = tx.Exec(`INSERT INTO "transaction" (id, type, from_account_id, to_account_id, amount) VALUES ($1, 'transfer', $2, $3, $4)`, transaction_id, from_account_id, to_account_id, amount); err != nil {
		return err
	}

	err = tx.Commit()

	return err
}
