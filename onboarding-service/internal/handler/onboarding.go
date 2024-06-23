package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/illenko/onboarding-service/internal/service"
	httpModel "github.com/illenko/onboarding-service/pkg/http"
)

type OnboardingHandler interface {
	CreateOnboarding(c *gin.Context)
	GetOnboarding(c *gin.Context)
	VerifySignature(c *gin.Context)
}

type OnboardingHandlerImpl struct {
	service service.OnboardingService
}

func NewOnboardingHandler(service service.OnboardingService) *OnboardingHandlerImpl {
	return &OnboardingHandlerImpl{service: service}
}

func (h *OnboardingHandlerImpl) CreateOnboarding(c *gin.Context) {
	var request httpModel.OnboardingRequest
	err := c.Bind(&request)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	response, err := h.service.CreateOnboarding(c, &request)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "Unable to create onboarding")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *OnboardingHandlerImpl) GetOnboarding(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	response, err := h.service.GetOnboarding(c, id)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "Unable to get onboarding")
		return
	}

	c.JSON(http.StatusOK, response)

}

func (h *OnboardingHandlerImpl) VerifySignature(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	var request httpModel.SignatureRequest
	err = c.Bind(&request)
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	response, err := h.service.VerifySignature(c, id, &request)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "Unable to verify signature")
		return
	}

	c.JSON(http.StatusOK, response)
}

func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, httpModel.ErrorResponse{
		Code:    "bad_request",
		Message: message,
	})
}
