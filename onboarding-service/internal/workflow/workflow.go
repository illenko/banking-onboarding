package workflow

import (
	"gittub.com/illenko/onboarding-service/internal/activity"
	"gittub.com/illenko/onboarding-service/internal/model"
	"gittub.com/illenko/onboarding-service/pkg/http"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func Onboarding(ctx workflow.Context, input model.UserRequest) (http.OnboardingResponse, error) {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    100 * time.Second,
		MaximumAttempts:    500,
	}

	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retryPolicy,
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// Withdraw money.
	var antifraudChecksOutput model.AntifraudResponse

	antifraudErr := workflow.ExecuteActivity(ctx, activity.AntifraudChecks, input).Get(ctx, &antifraudChecksOutput)

	if antifraudErr != nil {
		return http.OnboardingResponse{}, antifraudErr
	}

	if !antifraudChecksOutput.Passed {
		return http.OnboardingResponse{Status: "fraud_not_passed"}, nil
	}

	// Create user.
	var userOutput model.UserResponse

	userErr := workflow.ExecuteActivity(ctx, activity.CreateUser, input).Get(ctx, &userOutput)

	if userErr != nil {
		return http.OnboardingResponse{}, userErr
	}

	// Create account.
	accountInput := model.AccountRequest{
		UserID:   userOutput.ID,
		Type:     "personal",
		Currency: "USD",
	}

	var accountOutput model.AccountResponse

	accountErr := workflow.ExecuteActivity(ctx, activity.CreateAccount, accountInput).Get(ctx, &accountOutput)

	if accountErr != nil {
		return http.OnboardingResponse{}, accountErr
	}

	// Create agreement.
	agreementInput := model.AgreementRequest{
		UserID:    userOutput.ID,
		AccountID: accountOutput.ID,
	}

	var agreementOutput model.AgreementResponse

	agreementErr := workflow.ExecuteActivity(ctx, activity.CreateAgreement, agreementInput).Get(ctx, &agreementOutput)

	if agreementErr != nil {
		return http.OnboardingResponse{}, agreementErr
	}

	// Create signature.
	signatureInput := model.SignatureRequest{
		AgreementID: agreementOutput.ID,
		Signature:   "signature",
	}

	var signatureOutput model.SignatureResponse

	signatureErr := workflow.ExecuteActivity(ctx, activity.CreateSignature, signatureInput).Get(ctx, &signatureOutput)

	if signatureErr != nil {
		return http.OnboardingResponse{}, signatureErr
	}

	// Create card.
	cardInput := model.CardRequest{
		AccountID: accountOutput.ID,
	}

	var cardOutput model.CardResponse

	cardErr := workflow.ExecuteActivity(ctx, activity.CreateCard, cardInput).Get(ctx, &cardOutput)

	if cardErr != nil {
		return http.OnboardingResponse{}, cardErr
	}

	return http.OnboardingResponse{
		ID:      userOutput.ID,
		Status:  "success",
		Comment: "User created",
		Account: http.Account{
			ID:       accountOutput.ID,
			Currency: accountOutput.Currency,
			Type:     accountOutput.Type,
			Iban:     accountOutput.Iban,
			Balance:  accountOutput.Balance,
		},
		Card: http.Card{
			ID:     cardOutput.ID,
			Number: cardOutput.Number,
			Expire: cardOutput.Expire,
			Cvv:    cardOutput.Cvv,
		},
	}, nil
}
