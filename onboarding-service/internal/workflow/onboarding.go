package workflow

import (
	"time"

	"github.com/illenko/onboarding-service/internal/activity"
	"github.com/illenko/onboarding-service/internal/input"
	"github.com/illenko/onboarding-service/internal/output"
	"github.com/illenko/onboarding-service/internal/query"
	"github.com/illenko/onboarding-service/internal/request"
	"github.com/illenko/onboarding-service/internal/response"
	"github.com/illenko/onboarding-service/internal/signal"
	"github.com/illenko/onboarding-service/pkg/state"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

const (
	AccountTypePersonal = "personal"
	CurrencyUAH         = "UAH"
)

func Onboarding(ctx workflow.Context, input input.Onboarding) (output.Onboarding, error) {
	logger := workflow.GetLogger(ctx)

	options := getDefaultActivityOptions()

	currentState := output.Onboarding{State: state.ProcessingState}

	err := workflow.SetQueryHandler(ctx, query.CurrentState, func() (output.Onboarding, error) {
		return currentState, nil
	})

	if err != nil {
		currentState = output.Onboarding{State: state.FailedState}
		return currentState, err
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	userInput := request.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		City:      input.City,
	}

	// 1. Execute antifraud checks
	logger.Info("Executing antifraud checks for user", "Email", userInput.Email)
	antifraudChecksResult, err := executeActivity[request.User, response.Antifraud](ctx, activity.AntifraudChecks, userInput)

	if err != nil {
		currentState = output.Onboarding{State: state.FailedState}
		return currentState, err
	}

	// If antifraud checks failed, return fraud_not_passed state
	if !antifraudChecksResult.Passed {
		logger.Warn("User did not pass antifraud checks", "Email", userInput.Email)
		currentState = output.Onboarding{State: state.FraudNotPassedState}
		return currentState, nil
	}

	// 2. Create user
	logger.Info("Creating user", "Email", userInput.Email)
	createUserResult, err := executeActivity[request.User, response.User](ctx, activity.CreateUser, userInput)
	if err != nil {
		logger.Error("Unable to create user", "Email", userInput.Email)
		currentState = output.Onboarding{State: state.FailedState}
		return currentState, err
	}

	// 3. Create account
	logger.Info("Creating account for user", "UserID", createUserResult.ID)
	accountInput := request.Account{
		UserID:   createUserResult.ID,
		Type:     AccountTypePersonal,
		Currency: CurrencyUAH,
	}

	createAccountResult, err := executeActivity[request.Account, response.Account](ctx, activity.CreateAccount, accountInput)
	if err != nil {
		logger.Error("Unable to create account", "UserID", createUserResult.ID)
		currentState = output.Onboarding{State: state.FailedState}
		return currentState, err
	}

	// 4. Create agreement
	logger.Info("Creating agreement for user", "UserID", createUserResult.ID, "AccountID", createAccountResult.ID)
	agreementInput := request.Agreement{
		UserID:    createUserResult.ID,
		AccountID: createAccountResult.ID,
	}

	createAgreementResult, err := executeActivity[request.Agreement, response.Agreement](ctx, activity.CreateAgreement, agreementInput)
	if err != nil {
		currentState = output.Onboarding{State: state.FailedState}
		return currentState, err
	}

	// 5. Wait for signature
	logger.Info("Waiting for agreement signature", "AgreementID", createAgreementResult.ID)
	currentState = output.Onboarding{
		State: state.WaitingForAgreementSignatureState,
		Data:  map[string]any{"link": createAgreementResult.Link},
	}

	var signatureSignal signal.Signature

	signalChan := workflow.GetSignalChannel(ctx, signal.SignatureSignal)
	signalChan.Receive(ctx, &signatureSignal)

	// Send signature for validation
	signatureInput := request.Signature{
		AgreementID: createAgreementResult.ID,
		Signature:   signatureSignal.Signature,
	}

	signatureResult, err := executeActivity[request.Signature, response.Signature](ctx, activity.ValidateSignature, signatureInput)
	if err != nil {
		logger.Error("Unable to validate signature", "AgreementID", createAgreementResult.ID)
		currentState = output.Onboarding{State: state.FailedState}
		return currentState, err
	}

	// If signature is not valid, return signature_not_valid state
	if !signatureResult.Valid {
		logger.Warn("Signature is not valid", "AgreementID", createAgreementResult.ID)
		currentState = output.Onboarding{State: state.SignatureNotValidSate}
		return currentState, nil
	}

	currentState = output.Onboarding{State: state.ProcessingState}

	// 6. Create card
	logger.Info("Creating card for account", "AccountID", createAccountResult.ID)
	cardInput := request.Card{
		AccountID: createAccountResult.ID,
	}

	createCardResult, err := executeActivity[request.Card, response.Card](ctx, activity.CreateCard, cardInput)
	if err != nil {
		logger.Error("Unable to create card", "AccountID", createAccountResult.ID)
		currentState = output.Onboarding{State: state.FailedState}
		return currentState, err
	}

	// 7. Onboarding completed
	logger.Info("Onboarding completed", "UserID", createUserResult.ID, "AccountID", createAccountResult.ID, "CardID", createCardResult.ID)
	currentState = toFinalState(createAccountResult, createCardResult)

	return currentState, nil
}

func getDefaultActivityOptions() workflow.ActivityOptions {
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

	return options
}

func executeActivity[I any, R any](ctx workflow.Context, activityFunc interface{}, input I) (R, error) {
	var res R
	err := workflow.ExecuteActivity(ctx, activityFunc, input).Get(ctx, &res)
	return res, err
}

func toFinalState(accountResult response.Account, cardResult response.Card) output.Onboarding {
	return output.Onboarding{
		State: state.CompletedState,
		Data: map[string]any{
			"account": accountResult,
			"card":    cardResult,
		},
	}
}
