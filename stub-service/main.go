package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(RandomErrorMiddleware())

	router.POST("/antifraud-service/checks", antifraudChecksHandler)
	router.POST("/user-service/users", usersHandler)
	router.POST("/agreement-service/agreements", agreementsHandler)
	router.POST("/signature-service/signatures", signaturesHandler)
	router.POST("/account-service/accounts", accountsHandler)
	router.POST("/card-service/cards", cardsHandler)

	router.Run(":8081")
}
