package server

import (
	"github.com/gin-gonic/gin"
	"github.com/illenko/onboarding-service/docs"
	"github.com/illenko/onboarding-service/internal/handler"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func New(onboardingHandler handler.OnboardingHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/onboarding", onboardingHandler.CreateOnboarding)
	router.POST("/onboarding/:id/signature", onboardingHandler.VerifySignature)
	router.GET("/onboarding/:id", onboardingHandler.GetOnboarding)
	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router
}
