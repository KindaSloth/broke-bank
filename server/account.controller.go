package server

import (
	"broke-bank/utils"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

		err = s.Repositories.AccountRepository.CreateAccount(user.Id.String(), req.Name, "active")
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
		accountId := ctx.Param("id")
		if accountId == "" {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			return
		}

		account, err := s.Repositories.AccountRepository.GetMyAccount(accountId)
		if err != nil {
			log.Println("[ERROR] [GetAccount] failed to get account: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to get account"})
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
