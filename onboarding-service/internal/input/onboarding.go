package input

import "github.com/google/uuid"

type Onboarding struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	City      string    `json:"city"`
}
