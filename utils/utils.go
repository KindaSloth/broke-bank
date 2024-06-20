package utils

import (
	"broke-bank/model"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func GetUser(ctx *gin.Context) (*model.User, error) {
	userJson := ctx.GetString("user")
	var user model.User
	err := json.Unmarshal([]byte(userJson), &user)

	return &user, err
}
