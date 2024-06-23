package activity

import (
	"context"
	"github.com/illenko/onboarding-service/internal/configuration"

	"github.com/go-resty/resty/v2"
	"github.com/illenko/onboarding-service/internal/request"
	"github.com/illenko/onboarding-service/internal/response"
	"github.com/illenko/onboarding-service/internal/util"
)

type RequestError struct{}

func (m *RequestError) Error() string {
	return "Request"
}

var client = resty.New()

func makeRequest(url string, requestBody any, responseBody any) error {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		SetResult(responseBody).
		Post(url)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return &RequestError{}
	}

	return nil
}

func AntifraudChecks(ctx context.Context, data request.User) (response.Antifraud, error) {
	var antifraudResponse response.Antifraud
	err := makeRequest(configuration.Get("ANTIFRAUD_SERVICE_URL")+"/checks", data, &antifraudResponse)
	return antifraudResponse, err
}
func CreateUser(ctx context.Context, data request.User) (response.User, error) {
	var userResponse response.User
	err := makeRequest(configuration.Get("USER_SERVICE_URL")+"/users", data, &userResponse)
	return userResponse, err
}

func CreateAccount(ctx context.Context, data request.Account) (response.Account, error) {
	var accountResponse response.Account
	err := makeRequest(configuration.Get("ACCOUNT_SERVICE_URL")+"/accounts", data, &accountResponse)
	return accountResponse, err
}

func CreateAgreement(ctx context.Context, data request.Agreement) (response.Agreement, error) {
	var agreementResponse response.Agreement
	err := makeRequest(configuration.Get("AGREEMENT_SERVICE_URL")+"/agreements", data, &agreementResponse)
	return agreementResponse, err
}

func ValidateSignature(ctx context.Context, data request.Signature) (response.Signature, error) {
	var signatureResponse response.Signature
	err := makeRequest(configuration.Get("SIGNATURE_SERVICE_URL")+"/signatures", data, &signatureResponse)
	return signatureResponse, err
}

func CreateCard(ctx context.Context, data request.Card) (response.Card, error) {
	var cardResponse response.Card
	err := makeRequest(configuration.Get("CARD_SERVICE_URL")+"/cards", data, &cardResponse)
	if err == nil {
		cardResponse.Cvv = "***" // Mask CVV
		cardResponse.Number = util.MaskPan(cardResponse.Number)
	}
	return cardResponse, err
}
