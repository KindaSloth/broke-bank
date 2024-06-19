package repository

import (
	"broke-bank/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/valkey-io/valkey-go"
)

type UserRepository struct {
	Pg     *sqlx.DB
	Valkey *valkey.Client
}

func (ur *UserRepository) GetUserById(id uuid.UUID) (*model.User, error) {
	user := new(model.User)
	err := ur.Pg.Get(
		user,
		`SELECT u.id, u.email, u.created_at
		FROM "user" u WHERE u.id=$1`,
		id,
	)

	return user, err
}
