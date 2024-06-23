package request

import "github.com/google/uuid"

type User struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	City      string `json:"city"`
}

type Agreement struct {
	UserID    uuid.UUID `json:"user_id"`
	AccountID uuid.UUID `json:"account_id"`
}

type Signature struct {
	AgreementID uuid.UUID `json:"agreement_id"`
	Signature   string    `json:"signature"`
}

type Account struct {
	UserID   uuid.UUID `json:"user_id"`
	Type     string    `json:"type"`
	Currency string    `json:"currency"`
}

type Card struct {
	AccountID uuid.UUID `json:"account_id"`
}
