package server

import (
	"broke-bank/repository"
	"broke-bank/server/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	repositories repository.Repositories
}

func New() Server {
	repos := repository.New()

	return Server{repositories: repos}
}

func (s *Server) SetupRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CorsMiddleWare())

	router.GET("/health-check", func(ctx *gin.Context) { ctx.JSON(200, gin.H{"message": "Broke Bank"}) })

	return router
}

func (s *Server) Run(addr string) {
	router := s.SetupRouter()

	router.Run(addr)
}
