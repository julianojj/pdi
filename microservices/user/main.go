package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func getUser(users []*User, userID string) *User {
	for _, user := range users {
		if user.ID == userID {
			return user
		}
	}
	return nil
}

func main() {
	r := gin.Default()

	users := []*User{
		{
			ID:    "f5a0a785-75d2-4d8b-ae33-0039ee24216f",
			Name:  "Juliano Silva",
			Email: "juliano.silva@stone.com.br",
		},
	}

	r.GET("/users/:id", func(ctx *gin.Context) {
		user := getUser(users, ctx.Param("id"))
		if user == nil {
			ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]any{
				"message": "user not found",
			})
			return
		}
		ctx.JSON(http.StatusOK, user)
	})

	r.Run(":8080")
}
