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

// CreateOnboarding godoc
// @Summary Create onboarding
// @Description Create onboarding
// @Tags onboarding
// @Accept json
// @Produce json
// @Param request body http.OnboardingRequest true "Onboarding request"
// @Success 200 {object} http.OnboardingStatus
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /onboarding [post]
func (h *OnboardingHandlerImpl) CreateOnboarding(c *gin.Context) {
	var request httpModel.OnboardingRequest
	if !bindRequest(c, &request) {
		return
	}

	response, err := h.service.CreateOnboarding(c, &request)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "Unable to create onboarding")
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetOnboarding godoc
// @Summary Get onboarding
// @Description Get onboarding
// @Tags onboarding
// @Accept json
// @Produce json
// @Param id path string true "Onboarding ID"
// @Success 200 {object} http.OnboardingStatus
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /onboarding/{id} [get]
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

// VerifySignature godoc
// @Summary Verify signature
// @Description Verify signature
// @Tags onboarding
// @Accept json
// @Produce json
// @Param id path string true "Onboarding ID"
// @Param request body http.SignatureRequest true "Signature request"
// @Success 200 {object} http.OnboardingStatus
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /onboarding/{id}/signature [post]
func (h *OnboardingHandlerImpl) VerifySignature(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	var request httpModel.SignatureRequest
	if !bindRequest(c, &request) {
		return
	}

	response, err := h.service.VerifySignature(c, id, &request)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "Unable to verify signature")
		return
	}

	c.JSON(http.StatusOK, response)
}

func bindRequest(c *gin.Context, req interface{}) bool {
	if err := c.Bind(req); err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return false
	}
	return true
}

func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, httpModel.ErrorResponse{
		Code:    http.StatusText(statusCode),
		Message: message,
	})
}
