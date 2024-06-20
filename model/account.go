package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Account struct {
	Id      uuid.UUID       `db:"id" json:"id"`
	UserId  uuid.UUID       `db:"user_id" json:"user_id"`
	Name    string          `db:"name" json:"name"`
	Balance decimal.Decimal `db:"balance" json:"balance"`
	// 'active' | 'inactive'
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
