package server

import (
	"broke-bank/utils"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

func (s *Server) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := RegisterRequest{}

		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		encrypted_password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			ctx.JSON(500, gin.H{"error": "Failed to hash password"})
			return
		}
		_, err = s.Repositories.UserRepository.GetUserByEmail(req.Email)
		if err != nil && err == sql.ErrNoRows {
			err := s.Repositories.UserRepository.CreateUser(req.Email, string(encrypted_password))
			if err != nil {
				log.Println("[ERROR] [Register] failed to create user: ", err)
				ctx.JSON(500, gin.H{"error": "Failed to create user"})
				return
			}

			ctx.Status(200)
			return
		}
		if err != nil {
			log.Println("[ERROR] [Register] an unexpected error occurred: ", err)
			ctx.JSON(500, gin.H{"error": "Unexpected error :("})
			return
		}

		ctx.JSON(409, gin.H{"error": "Email already registered"})
	}
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required,min=8,max=255"`
}

func (s *Server) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := LoginRequest{}

		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		user, err := s.Repositories.UserRepository.GetUserByEmail(req.Email)
		if err != nil {
			ctx.JSON(409, gin.H{"error": "Email not registered"})
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(req.Password))
		if err != nil {
			ctx.JSON(409, gin.H{"error": "Wrong password"})
			return
		}

		session_id, err := uuid.NewV7()
		if err != nil {
			log.Println("[ERROR] [Login] an unexpected error occurred while creating session ID: ", err)
			ctx.JSON(500, gin.H{"error": "Unexpected error :("})
			return
		}

		valkey := s.Repositories.Valkey
		err = valkey.Do(ctx, valkey.B().Set().Key(session_id.String()).Value(user.Id.String()).Nx().Build()).Error()
		if err != nil {
			log.Println("[ERROR] [Login] an unexpected error occurred while storing user session: ", err)
			ctx.JSON(500, gin.H{"error": "Unexpected error :("})
			return
		}

		ctx.SetCookie("sessionId", session_id.String(), 3600*24, "/", "localhost", true, true)
		ctx.Status(200)
	}
}

type MeResponse struct {
	Email string `json:"user_email"`
}

func (s *Server) Me() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [CreateAccount] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		ctx.JSON(200, gin.H{"payload": MeResponse{Email: user.Email}})
	}
}

type GetMyAccountsRequest struct {
	Limit  int `db:"limit" form:"limit"`
	Offset int `db:"offset" form:"offset"`
}

type GetAccountsResponse struct {
	Id      uuid.UUID `json:"id"`
	Name    string    `json:"name"`
	Balance string    `json:"balance"`
	// 'active' | 'inactive'
	Status string `json:"status"`
}

func (s *Server) GetMyAccounts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := GetMyAccountsRequest{}

		if ctx.ShouldBindQuery(&req) != nil {
			ctx.JSON(400, gin.H{"error": "Invalid input"})
			ctx.Abort()
			return
		}

		if req.Limit == 0 {
			req.Limit = 10
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [GetMyAccounts] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		raw_accounts, err := s.Repositories.AccountRepository.GetMyAccounts(user.Id.String(), req.Limit, req.Offset)
		if err != nil {
			log.Println("[ERROR] [GetMyAccounts] failed to retrieve accounts: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to retrieve accounts"})
			return
		}

		accounts := []GetAccountsResponse{}
		for _, value := range *raw_accounts {
			accounts = append(accounts, GetAccountsResponse{
				Id:      value.Id,
				Name:    value.Name,
				Balance: value.Balance.StringFixed(2),
				Status:  value.Status,
			})
		}

		ctx.JSON(200, gin.H{"payload": accounts})
	}
}
