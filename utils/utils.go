package utils

import (
	"broke-bank/model"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUser(ctx *gin.Context) (*model.User, error) {
	userJson := ctx.GetString("user")
	var user model.User
	err := json.Unmarshal([]byte(userJson), &user)

	return &user, err
}

/*
SortUUIDs sorts two UUIDs and returns them in lexicographical order.

TODO: Maybe sort them based on timestamp would be better? since I'm using UUID v7.
*/
func SortUUIDs(first_id uuid.UUID, second_id uuid.UUID) (uuid.UUID, uuid.UUID) {
	if first_id.String() < second_id.String() {
		return first_id, second_id
	}
	return second_id, first_id
}

// Same as SortUUIDs, but using string type directly.
func SortStringUUIDs(first_id string, second_id string) (string, string) {
	if first_id < second_id {
		return first_id, second_id
	}
	return second_id, first_id
}
