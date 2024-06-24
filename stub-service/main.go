package main

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	sloggin "github.com/samber/slog-gin"
)

func main() {
	router := gin.New()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	config := sloggin.Config{
		WithRequestBody:    true,
		WithResponseBody:   true,
		WithRequestHeader:  true,
		WithResponseHeader: true,
	}

	router.Use(sloggin.NewWithConfig(logger, config))
	router.Use(RandomErrorMiddleware())
	router.Use(gin.Recovery())

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
