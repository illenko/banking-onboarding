package main

import "github.com/google/uuid"

type UserRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	City      string `json:"city"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	City      string    `json:"city"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type AntifraudResponse struct {
	Passed  bool   `json:"passed"`
	Comment string `json:"comment"`
}

type AgreementRequest struct {
	UserID    uuid.UUID `json:"user_id"`
	AccountID uuid.UUID `json:"account_id"`
}

type AgreementResponse struct {
	ID   uuid.UUID `json:"id"`
	Link string    `json:"link"`
}

type SignatureRequest struct {
	AgreementID uuid.UUID `json:"agreement_id"`
	Signature   string    `json:"signature"`
}

type SignatureResponse struct {
	ID      uuid.UUID `json:"id"`
	Valid   bool      `json:"valid"`
	Comment string    `json:"comment"`
}

type AccountRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	Type     string    `json:"type"`
	Currency string    `json:"currency"`
}

type AccountResponse struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Currency string    `json:"currency"`
	Type     string    `json:"type"`
	Iban     string    `json:"iban"`
	Balance  float64   `json:"balance"`
}

type CardRequest struct {
	AccountID uuid.UUID `json:"account_id"`
}

type CardResponse struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Number    string    `json:"number"`
	Expire    string    `json:"expire"`
	Cvv       string    `json:"cvv"`
}
