package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id uuid.UUID `db:"id" json:"id"`
	// 'deposit' | 'withdrawal' | 'transfer'
	Type          string          `db:"type" json:"type"`
	FromAccountId *uuid.UUID      `db:"from_account_id" json:"from_account_id"`
	ToAccountId   *uuid.UUID      `db:"to_account_id" json:"to_account_id"`
	DateIssued    time.Time       `db:"date_issued" json:"date_issued"`
	Amount        decimal.Decimal `db:"amount" json:"amount"`
}
