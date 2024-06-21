package server

import (
	"broke-bank/repository"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Repositories repository.Repositories
}

func New() Server {
	repos := repository.New()

	return Server{Repositories: repos}
}

func (s *Server) SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(CorsMiddleware())

	router.GET("/health-check", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Broke Bank"}) })
	router.POST("/register", s.Register())
	router.POST("/login", s.Login())

	router.Use(s.AuthMiddleware())

	// User endpoints
	router.GET("/me", s.Me())
	router.GET("/myAccounts", s.GetMyAccounts())

	// Account endpoints
	router.GET("/account/:id", s.GetAccount())
	router.POST("/account/create", s.CreateAccount())
	router.PATCH("/account/disable/:id", s.DisableAccount())

	// Transaction endpoints
	router.GET("/transaction/:id", s.GetTransaction())
	router.POST("/transaction/deposit", s.DepositTransaction())
	router.POST("/transaction/withdrawal", s.WithdrawalTransaction())
	router.POST("/transaction/transfer", s.TransferTransaction())

	return router
}

func (s *Server) Run(addr string) {
	router := s.SetupRouter()

	router.Run(addr)
}
