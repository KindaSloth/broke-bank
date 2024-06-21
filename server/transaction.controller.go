package server

import (
	"log"

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
		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

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

func (s *Server) WithdrawalTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}

func (s *Server) TransferTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {}
}
