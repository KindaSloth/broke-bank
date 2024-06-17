package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/valkey-io/valkey-go"
)

type UserRepository struct {
	pg     *sqlx.DB
	valkey *valkey.Client
}
