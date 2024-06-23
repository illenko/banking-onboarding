package response

import "github.com/google/uuid"

type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	City      string    `json:"city"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Antifraud struct {
	Passed  bool   `json:"passed"`
	Comment string `json:"comment"`
}

type Agreement struct {
	ID   uuid.UUID `json:"id"`
	Link string    `json:"link"`
}

type Signature struct {
	ID      uuid.UUID `json:"id"`
	Valid   bool      `json:"valid"`
	Comment string    `json:"comment"`
}

type Account struct {
	ID       uuid.UUID `json:"id"`
	UserID   uuid.UUID `json:"user_id"`
	Currency string    `json:"currency"`
	Type     string    `json:"type"`
	Iban     string    `json:"iban"`
	Balance  float64   `json:"balance"`
}

type Card struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	Number    string    `json:"number"`
	Expire    string    `json:"expire"`
	Cvv       string    `json:"cvv"`
}
