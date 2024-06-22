package activity

import (
	"context"
	"github.com/go-resty/resty/v2"
	"gittub.com/illenko/onboarding-service/internal/model"
)

type FraudCheckError struct{}

func (m *FraudCheckError) Error() string {
	return "Client is fraud"
}

type RequestError struct{}

func (m *RequestError) Error() string {
	return "Request"
}

func AntifraudChecks(ctx context.Context, data model.UserRequest) (model.AntifraudResponse, error) {
	client := resty.New()

	var antifraudResponse model.AntifraudResponse

	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		SetResult(&antifraudResponse).
		Post("http://localhost:8081/antifraud-service/checks")

	if err != nil {
		return model.AntifraudResponse{}, err
	}

	if resp.IsError() {
		return model.AntifraudResponse{}, &RequestError{}
	}

	if antifraudResponse.Passed != true {
		return model.AntifraudResponse{}, &FraudCheckError{}
	}

	return antifraudResponse, nil
}
