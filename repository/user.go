package repository

import (
	"broke-bank/model"
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	Pg *sqlx.DB
}

func (ur *UserRepository) CreateUser(ctx context.Context, conn *sqlx.Conn, email string, password string) error {
	_, err := conn.ExecContext(
		ctx,
		`INSERT INTO "user" (email, password)
		VALUES ($1, $2)`,
		email,
		password,
	)

	return err
}

func (ur *UserRepository) GetUserById(ctx context.Context, conn *sqlx.Conn, id uuid.UUID) (*model.User, error) {
	user := new(model.User)
	err := conn.GetContext(
		ctx,
		user,
		`SELECT u.id, u.email, u.password, u.created_at, u.updated_at
		FROM "user" u WHERE u.id=$1`,
		id,
	)

	return user, err
}

func (ur *UserRepository) GetUserByEmail(ctx context.Context, conn *sqlx.Conn, email string) (*model.User, error) {
	user := new(model.User)
	err := conn.GetContext(
		ctx,
		user,
		`SELECT u.id, u.email, u.password, u.created_at, u.updated_at
		FROM "user" u WHERE u.email=$1`,
		email,
	)

	return user, err
}
