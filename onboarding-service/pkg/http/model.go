package http

import (
	"github.com/google/uuid"
	"github.com/illenko/onboarding-service/pkg/state"
)

type OnboardingRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	City      string `json:"city"`
}

type OnboardingStatus struct {
	ID    uuid.UUID             `json:"id"`
	State state.OnboardingState `json:"state"`
	Data  map[string]any        `json:"data,omitempty"`
}

type SignatureRequest struct {
	Signature string `json:"signature"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
