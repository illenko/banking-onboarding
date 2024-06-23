package http

import "github.com/google/uuid"

type OnboardingStatus struct {
	ID    uuid.UUID              `json:"id"`
	State string                 `json:"state"`
	Data  map[string]interface{} `json:"data"`
}
