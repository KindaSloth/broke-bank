package server

import (
	"broke-bank/utils"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Server) GetTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		transaction_id := ctx.Param("id")
		if transaction_id == "" {
			ctx.JSON(400, gin.H{"error": "Missing id param"})
			return
		}

		transaction, err := s.Repositories.TransactionRepository.GetTransaction(transaction_id)
		if err != nil {
			log.Printf("[ERROR] [GetTransaction] failed to get transaction: %s, transaction ID: %s\n", err, transaction_id)
			ctx.JSON(500, gin.H{"error": "Failed to get transaction"})
			return
		}

		ctx.JSON(200, gin.H{"payload": transaction})
	}
}

type DepositTransactionRequest struct {
	Amount      decimal.Decimal `json:"amount"`
	ToAccountId string          `json:"to_account_id"`
}

func (s *Server) DepositTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := DepositTransactionRequest{}
		if ctx.ShouldBindJSON(&req) != nil || req.Amount.LessThan(decimal.NewFromInt(0)) {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		if _, err := s.Repositories.Pg.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`); err != nil {
			log.Println("[ERROR] [DepositTransaction] failed to set SERIALIZABLE transaction level: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete deposit transaction"})
			return
		}
		// TODO: handle err
		defer s.Repositories.Pg.Exec(`SET TRANSACTION ISOLATION LEVEL READ COMMITTED`)

		transaction_id, err := uuid.NewV7()
		if err != nil {
			log.Println("[ERROR] [DepositTransaction] an unexpected error occurred while creating transaction ID: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete deposit transaction"})
			return
		}

		tx, err := s.Repositories.TransactionRepository.GetTransaction(transaction_id.String())
		if tx != nil && err == nil {
			log.Println("[ERROR] [DepositTransaction] duplicated transaction: ", err)
			ctx.JSON(500, gin.H{"error": "Duplicated transaction"})
			return
		}

		err = s.Repositories.TransactionRepository.DepositTransaction(transaction_id, req.ToAccountId, req.Amount)
		if err != nil {
			log.Println("[ERROR] [DepositTransaction] failed to complete deposit transaction: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete deposit transaction"})
			return
		}

		ctx.Status(200)
	}
}

type WithdrawalTransactionRequest struct {
	Amount        decimal.Decimal `json:"amount"`
	FromAccountId string          `json:"from_account_id"`
}

func (s *Server) WithdrawalTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := WithdrawalTransactionRequest{}
		if ctx.ShouldBindJSON(&req) != nil || req.Amount.LessThan(decimal.NewFromInt(0)) {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [WithdrawalTransaction] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		if _, err = s.Repositories.Pg.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`); err != nil {
			log.Println("[ERROR] [WithdrawalTransaction] failed to set SERIALIZABLE transaction level: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete withdrawal transaction"})
			return
		}
		// TODO: handle err
		defer s.Repositories.Pg.Exec(`SET TRANSACTION ISOLATION LEVEL READ COMMITTED`)

		account, err := s.Repositories.AccountRepository.GetAccount(req.FromAccountId)
		if err != nil {
			log.Printf("[ERROR] [WithdrawalTransaction] failed to get account: %s, account ID: %s\n", err, req.FromAccountId)
			ctx.JSON(500, gin.H{"error": "Failed to get account"})
			return
		}

		if account.UserId != user.Id {
			ctx.JSON(500, gin.H{"error": "This account does not belongs to the user"})
			return
		}

		if account.Balance.LessThan(req.Amount) {
			ctx.JSON(500, gin.H{"error": "Insufficient account balance"})
			return
		}

		transaction_id, err := uuid.NewV7()
		if err != nil {
			log.Println("[ERROR] [WithdrawalTransaction] an unexpected error occurred while creating transaction ID: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete withdrawal transaction"})
			return
		}

		tx, err := s.Repositories.TransactionRepository.GetTransaction(transaction_id.String())
		if tx != nil && err == nil {
			log.Println("[ERROR] [WithdrawalTransaction] duplicated transaction: ", err)
			ctx.JSON(500, gin.H{"error": "Duplicated transaction"})
			return
		}

		err = s.Repositories.TransactionRepository.WithdrawalTransaction(transaction_id, req.FromAccountId, req.Amount)
		if err != nil {
			log.Println("[ERROR] [WithdrawalTransaction] failed to complete withdrawal transaction: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete withdrawal transaction"})
			return
		}

		ctx.Status(200)
	}
}

type TransferTransactionRequest struct {
	Amount        decimal.Decimal `json:"amount"`
	FromAccountId string          `json:"from_account_id"`
	ToAccountId   string          `json:"to_account_id"`
}

func (s *Server) TransferTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := TransferTransactionRequest{}
		if ctx.ShouldBindJSON(&req) != nil || req.Amount.LessThan(decimal.NewFromInt(0)) || req.FromAccountId == req.ToAccountId {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [TransferTransaction] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		if _, err = s.Repositories.Pg.Exec(`SET TRANSACTION ISOLATION LEVEL SERIALIZABLE`); err != nil {
			log.Println("[ERROR] [TransferTransaction] failed to set SERIALIZABLE transaction level: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete transfer transaction"})
			return
		}
		// TODO: handle err
		defer s.Repositories.Pg.Exec(`SET TRANSACTION ISOLATION LEVEL READ COMMITTED`)

		from_account, err := s.Repositories.AccountRepository.GetAccount(req.FromAccountId)
		if err != nil {
			log.Printf("[ERROR] [TransferTransaction] failed to get sender account: %s, account ID: %s\n", err, req.FromAccountId)
			ctx.JSON(500, gin.H{"error": "Failed to get sender account"})
			return
		}

		if from_account.UserId != user.Id {
			ctx.JSON(500, gin.H{"error": "This account does not belongs to the user"})
			return
		}

		if from_account.Balance.LessThan(req.Amount) {
			ctx.JSON(500, gin.H{"error": "Insufficient account balance"})
			return
		}

		_, err = s.Repositories.AccountRepository.GetAccount(req.ToAccountId)
		if err != nil {
			log.Printf("[ERROR] [TransferTransaction] failed to get receiver account: %s, account ID: %s\n", err, req.ToAccountId)
			ctx.JSON(500, gin.H{"error": "Failed to get receiver account"})
			return
		}

		transaction_id, err := uuid.NewV7()
		if err != nil {
			log.Println("[ERROR] [TransferTransaction] an unexpected error occurred while creating transaction ID: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete transfer transaction"})
			return
		}

		tx, err := s.Repositories.TransactionRepository.GetTransaction(transaction_id.String())
		if tx != nil && err == nil {
			log.Println("[ERROR] [TransferTransaction] duplicated transaction: ", err)
			ctx.JSON(500, gin.H{"error": "Duplicated transaction"})
			return
		}

		max_retries := 3

		for i := 0; i < max_retries; i++ {
			err = s.Repositories.TransactionRepository.TransferTransaction(transaction_id, req.FromAccountId, req.ToAccountId, req.Amount)
			if err == nil {
				return
			}

			if (i + 1) == max_retries {
				log.Println("[ERROR] [TransferTransaction] failed to complete transfer transaction: ", err)
				ctx.JSON(500, gin.H{"error": "Failed to complete transfer transaction"})
				return
			}

			time.Sleep(time.Millisecond * time.Duration(100*i))
		}

		ctx.Status(200)
	}
}
