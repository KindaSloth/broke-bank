package repository

import (
	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	Pg *sqlx.DB
}

func (ac *AccountRepository) CreateAccount(user_id string, name string, status string) error {
	_, err := ac.Pg.Exec(
		`INSERT INTO "account" (user_id, name, balance, status)
		VALUES ($1, $2, 0, $3)
		`,
		user_id,
		name,
		status,
	)

	return err
}
