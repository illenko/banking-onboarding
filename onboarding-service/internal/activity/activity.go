package activity

import (
	"context"
	"github.com/go-resty/resty/v2"
	"gittub.com/illenko/onboarding-service/internal/model"
)

type RequestError struct{}

func (m *RequestError) Error() string {
	return "Request"
}

var client = resty.New()

func makeRequest(url string, requestBody interface{}, responseBody interface{}) error {
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

func AntifraudChecks(ctx context.Context, data model.UserRequest) (model.AntifraudResponse, error) {
	var antifraudResponse model.AntifraudResponse
	err := makeRequest("http://localhost:8081/antifraud-service/checks", data, &antifraudResponse)
	return antifraudResponse, err
}

func CreateUser(ctx context.Context, data model.UserRequest) (model.UserResponse, error) {
	var userResponse model.UserResponse
	err := makeRequest("http://localhost:8081/user-service/users", data, &userResponse)
	return userResponse, err
}

func CreateAccount(ctx context.Context, data model.AccountRequest) (model.AccountResponse, error) {
	var accountResponse model.AccountResponse
	err := makeRequest("http://localhost:8081/account-service/accounts", data, &accountResponse)
	return accountResponse, err
}

func CreateAgreement(ctx context.Context, data model.AgreementRequest) (model.AgreementResponse, error) {
	var agreementResponse model.AgreementResponse
	err := makeRequest("http://localhost:8081/agreement-service/agreements", data, &agreementResponse)
	return agreementResponse, err
}

func CreateSignature(ctx context.Context, data model.SignatureRequest) (model.SignatureResponse, error) {
	var signatureResponse model.SignatureResponse
	err := makeRequest("http://localhost:8081/signature-service/signatures", data, &signatureResponse)
	return signatureResponse, err
}

func CreateCard(ctx context.Context, data model.CardRequest) (model.CardResponse, error) {
	var cardResponse model.CardResponse
	err := makeRequest("http://localhost:8081/card-service/cards", data, &cardResponse)
	return cardResponse, err
}
