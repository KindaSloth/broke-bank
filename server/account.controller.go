package server

import (
	"broke-bank/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CreateAccountRequest struct {
	Name string `json:"name" validate:"required"`
}

func (s *Server) CreateAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := CreateAccountRequest{}
		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [CreateAccount] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		conn, err := s.Repositories.Pg.Connx(ctx)
		if err != nil {
			log.Println("[ERROR] [CreateAccount] failed to get connection from db pool: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to create account"})
			return
		}
		defer conn.Close()

		err = s.Repositories.AccountRepository.CreateAccount(ctx, conn, user.Id.String(), req.Name, "active")
		if err != nil {
			log.Println("[ERROR] [CreateAccount] failed to create account: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to create account"})
			return
		}

		ctx.Status(200)
	}
}

type GetAccountResponse struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Balance string    `json:"balance"`
	// 'active' | 'inactive'
	Status string `json:"status"`
}

func (s *Server) GetAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account_id := ctx.Param("id")
		if account_id == "" {
			ctx.JSON(400, gin.H{"error": "Missing id param"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [GetAccount] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		conn, err := s.Repositories.Pg.Connx(ctx)
		if err != nil {
			log.Println("[ERROR] [GetAccount] failed to get connection from db pool: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to get account"})
			return
		}
		defer conn.Close()

		account, err := s.Repositories.AccountRepository.GetAccount(ctx, conn, account_id)
		if err != nil {
			log.Printf("[ERROR] [GetAccount] failed to get account: %s, account ID: %s\n", err, account_id)
			ctx.JSON(500, gin.H{"error": "Failed to get account"})
			return
		}

		if account.UserId != user.Id {
			ctx.JSON(500, gin.H{"error": "This account does not belongs to the user"})
			return
		}

		ctx.JSON(200, gin.H{"payload": GetAccountResponse{
			Id:      account.Id,
			Name:    account.Name,
			Balance: account.Balance.StringFixed(2),
			Status:  account.Status,
		}})
	}
}

func (s *Server) DisableAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account_id := ctx.Param("id")
		if account_id == "" {
			ctx.JSON(400, gin.H{"error": "Missing id param"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [DisableAccount] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		conn, err := s.Repositories.Pg.Connx(ctx)
		if err != nil {
			log.Println("[ERROR] [DisableAccount] failed to get connection from db pool: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to disable account"})
			return
		}
		defer conn.Close()

		account, err := s.Repositories.AccountRepository.GetAccount(ctx, conn, account_id)
		if err != nil {
			log.Println("[ERROR] [DisableAccount] failed to get account: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to get account"})
			return
		}

		if account.UserId != user.Id {
			ctx.JSON(500, gin.H{"error": "This account does not belongs to the user"})
			return
		}

		if account.Balance.GreaterThan(decimal.NewFromInt(0)) {
			ctx.JSON(500, gin.H{"error": "Account still has balance and cannot be deleted"})
			return
		}

		err = s.Repositories.AccountRepository.DisableAccount(ctx, conn, account_id)
		if err != nil {
			log.Println("[ERROR] [DisableAccount] failed to disable account: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to disable account"})
			return
		}

		ctx.Status(200)
	}
}
