package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	httpModel "github.com/illenko/onboarding-service/pkg/http"
)

type onboardingService interface {
	CreateOnboarding(ctx context.Context, req *httpModel.OnboardingRequest) (httpModel.OnboardingStatus, error)
	GetOnboarding(ctx context.Context, id uuid.UUID) (httpModel.OnboardingStatus, error)
	SignAgreement(ctx context.Context, id uuid.UUID, req *httpModel.SignatureRequest) (httpModel.OnboardingStatus, error)
}

type OnboardingHandler struct {
	service onboardingService
}

func NewOnboardingHandler(service onboardingService) *OnboardingHandler {
	return &OnboardingHandler{service: service}
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
func (h *OnboardingHandler) CreateOnboarding(c *gin.Context) {
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
func (h *OnboardingHandler) GetOnboarding(c *gin.Context) {
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

// SignAgreement godoc
// @Summary Sign agreement
// @Description Sign agreement
// @Tags onboarding
// @Accept json
// @Produce json
// @Param id path string true "Onboarding ID"
// @Param request body http.SignatureRequest true "Signature request"
// @Success 200 {object} http.OnboardingStatus
// @Failure 400 {object} http.ErrorResponse
// @Failure 500 {object} http.ErrorResponse
// @Router /onboarding/{id}/signature [post]
func (h *OnboardingHandler) SignAgreement(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	var request httpModel.SignatureRequest
	if !bindRequest(c, &request) {
		return
	}

	response, err := h.service.SignAgreement(c, id, &request)
	if err != nil {
		sendErrorResponse(c, http.StatusInternalServerError, "Unable to sign agreement")
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
