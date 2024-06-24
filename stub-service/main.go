package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(RequestLogger())
	router.Use(ResponseLogger())
	router.Use(RandomErrorMiddleware())

	router.POST("/antifraud-service/checks", antifraudChecksHandler)
	router.POST("/user-service/users", usersHandler)
	router.POST("/agreement-service/agreements", agreementsHandler)
	router.POST("/signature-service/signatures", signaturesHandler)
	router.POST("/account-service/accounts", accountsHandler)
	router.POST("/card-service/cards", cardsHandler)

	err := router.Run(":8081")
	if err != nil {
		return
	}
}
