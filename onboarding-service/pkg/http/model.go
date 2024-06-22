package http

import "github.com/google/uuid"

type OnboardingResponse struct {
	ID      uuid.UUID `json:"id"`
	Status  string    `json:"status"`
	Comment string    `json:"comment"`
	Account Account   `json:"account"`
	Card    Card      `json:"card"`
}

type Account struct {
	ID       uuid.UUID `json:"id"`
	Currency string    `json:"currency"`
	Type     string    `json:"type"`
	Iban     string    `json:"iban"`
	Balance  float64   `json:"balance"`
}

type Card struct {
	ID     uuid.UUID `json:"id"`
	Number string    `json:"number"`
	Expire string    `json:"expire"`
	Cvv    string    `json:"cvv"`
}
