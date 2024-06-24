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

func AntifraudChecks(ctx context.Context, req request.Base) (response.Antifraud, error) {
	return httpActivity[response.Antifraud](ctx, configuration.Get("ANTIFRAUD_SERVICE_URL")+"/checks", req)
}

func CreateUser(ctx context.Context, req request.Base) (response.User, error) {
	return httpActivity[response.User](ctx, configuration.Get("USER_SERVICE_URL")+"/users", req)
}

func CreateAccount(ctx context.Context, req request.Base) (response.Account, error) {
	return httpActivity[response.Account](ctx, configuration.Get("ACCOUNT_SERVICE_URL")+"/accounts", req)
}

func CreateAgreement(ctx context.Context, req request.Base) (response.Agreement, error) {
	return httpActivity[response.Agreement](ctx, configuration.Get("AGREEMENT_SERVICE_URL")+"/agreements", req)
}

func ValidateSignature(ctx context.Context, req request.Base) (response.Signature, error) {
	return httpActivity[response.Signature](ctx, configuration.Get("SIGNATURE_SERVICE_URL")+"/signatures", req)
}

func CreateCard(ctx context.Context, req request.Base) (response.Card, error) {
	cardResponse, err := httpActivity[response.Card](ctx, configuration.Get("CARD_SERVICE_URL")+"/cards", req)
	if err == nil {
		cardResponse.Cvv = "***" // Mask CVV
		cardResponse.Number = util.MaskPan(cardResponse.Number)
	}
	return cardResponse, err
}

var client = resty.New().SetDebug(true)

func makeRequest[T any](url string, req request.Base, res *T) error {
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetHeaders(req.Headers).
		SetBody(req.Body).
		SetResult(res).
		Post(url)

	if err != nil {
		return err
	}

	if resp.IsError() {
		return &RequestError{}
	}

	return nil
}

func httpActivity[T any](ctx context.Context, url string, req request.Base) (T, error) {
	var res T
	err := makeRequest(url, req, &res)
	return res, err
}
