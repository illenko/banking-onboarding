package server

import (
	"github.com/gin-gonic/gin"
	"github.com/illenko/onboarding-service/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type onboardingHandler interface {
	CreateOnboarding(c *gin.Context)
	GetOnboarding(c *gin.Context)
	SignAgreement(c *gin.Context)
}

func New(handler onboardingHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/onboarding", handler.CreateOnboarding)
	router.POST("/onboarding/:id/signature", handler.SignAgreement)
	router.GET("/onboarding/:id", handler.GetOnboarding)
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
