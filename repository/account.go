package repository

import (
	"broke-bank/model"
	"context"

	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	Pg *sqlx.DB
}

func (ac *AccountRepository) CreateAccount(ctx context.Context, conn *sqlx.Conn, user_id string, name string, status string) error {
	_, err := conn.ExecContext(
		ctx,
		`INSERT INTO "account" (user_id, name, balance, status)
		VALUES ($1, $2, 0, $3)
		`,
		user_id,
		name,
		status,
	)

	return err
}

func (ac *AccountRepository) GetAccount(ctx context.Context, conn *sqlx.Conn, acc_id string) (*model.Account, error) {
	account := new(model.Account)
	err := conn.GetContext(
		ctx,
		account,
		`SELECT acc.id, acc.user_id, acc.name, acc.balance, acc.status, acc.created_at, acc.updated_at 
		FROM "account" acc WHERE acc.id = $1`,
		acc_id,
	)

	return account, err
}

func (ac *AccountRepository) GetMyAccounts(ctx context.Context, conn *sqlx.Conn, user_id string, limit int, offset int) (*[]model.Account, error) {
	accounts := new([]model.Account)
	err := conn.SelectContext(
		ctx,
		accounts,
		`
		SELECT 
			acc.id, acc.user_id, acc.name, acc.balance, acc.status, acc.created_at, acc.updated_at 
		FROM 
			"account" acc 
		WHERE 
			acc.user_id = $1
		ORDER BY
			acc.status
		LIMIT 
			$2
		OFFSET 
			$3
		`,
		user_id,
		limit,
		offset,
	)

	return accounts, err
}

func (ac *AccountRepository) DisableAccount(ctx context.Context, conn *sqlx.Conn, acc_id string) error {
	_, err := conn.ExecContext(
		ctx,
		`UPDATE "account"
		SET status = 'inactive'
		WHERE id = $1`,
		acc_id,
	)

	return err
}
