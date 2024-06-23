package server

import (
	"github.com/gin-gonic/gin"
	"github.com/illenko/onboarding-service/internal/handler"
)

func New(onboardingHandler handler.OnboardingHandler) *gin.Engine {
	router := gin.Default()

	router.POST("/onboarding", onboardingHandler.CreateOnboarding)
	router.POST("/onboarding/:id/signature", onboardingHandler.VerifySignature)
	router.GET("/onboarding/:id", onboardingHandler.GetOnboarding)

	return router
}
