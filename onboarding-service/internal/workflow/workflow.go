package workflow

import (
	"time"

	"gittub.com/illenko/onboarding-service/internal/activity"
	"gittub.com/illenko/onboarding-service/internal/model"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type SignatureSignal struct {
	Signature string
}

type CurrentState struct {
	State string                 `json:"state"`
	Data  map[string]interface{} `json:"data"`
}

func Onboarding(ctx workflow.Context, input model.UserRequest) (CurrentState, error) {
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

	currentState := CurrentState{State: "processing"}
	queryType := "current_state"

	err := workflow.SetQueryHandler(ctx, queryType, func() (CurrentState, error) {
		return currentState, nil
	})

	if err != nil {
		currentState = CurrentState{State: "failed"}
		return currentState, err
	}

	ctx = workflow.WithActivityOptions(ctx, options)

	// Withdraw money.
	var antifraudChecksOutput model.AntifraudResponse

	antifraudErr := workflow.ExecuteActivity(ctx, activity.AntifraudChecks, input).Get(ctx, &antifraudChecksOutput)

	if antifraudErr != nil {
		currentState = CurrentState{State: "failed"}
		return currentState, antifraudErr
	}

	if !antifraudChecksOutput.Passed {
		currentState = CurrentState{
			State: "fraud_not_passed",
			Data:  map[string]interface{}{"comment": antifraudChecksOutput.Comment}}
		return currentState, nil
	}

	// Create user.
	var userOutput model.UserResponse

	userErr := workflow.ExecuteActivity(ctx, activity.CreateUser, input).Get(ctx, &userOutput)

	if userErr != nil {
		currentState = CurrentState{State: "failed"}
		return currentState, userErr
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
		currentState = CurrentState{State: "failed"}
		return currentState, accountErr
	}

	// Create agreement.
	agreementInput := model.AgreementRequest{
		UserID:    userOutput.ID,
		AccountID: accountOutput.ID,
	}

	var agreementOutput model.AgreementResponse

	agreementErr := workflow.ExecuteActivity(ctx, activity.CreateAgreement, agreementInput).Get(ctx, &agreementOutput)

	if agreementErr != nil {
		currentState = CurrentState{State: "failed"}
		return currentState, agreementErr
	}

	currentState = CurrentState{
		State: "waiting_for_signature",
		Data:  map[string]interface{}{"link": agreementOutput.Link},
	}

	var signal SignatureSignal

	signalChan := workflow.GetSignalChannel(ctx, "signature-signal")
	signalChan.Receive(ctx, &signal)

	// Create signature.
	signatureInput := model.SignatureRequest{
		AgreementID: agreementOutput.ID,
		Signature:   signal.Signature,
	}

	var signatureOutput model.SignatureResponse

	signatureErr := workflow.ExecuteActivity(ctx, activity.CreateSignature, signatureInput).Get(ctx, &signatureOutput)

	if signatureErr != nil {
		currentState = CurrentState{State: "failed"}
		return currentState, signatureErr
	}

	currentState = CurrentState{State: "processing"}

	// Create card.
	cardInput := model.CardRequest{
		AccountID: accountOutput.ID,
	}

	var cardOutput model.CardResponse

	cardErr := workflow.ExecuteActivity(ctx, activity.CreateCard, cardInput).Get(ctx, &cardOutput)

	if cardErr != nil {
		currentState = CurrentState{State: "failed"}
		return currentState, cardErr
	}

	currentState = CurrentState{
		State: "completed",
		Data:  map[string]interface{}{"account": accountOutput, "card": cardOutput},
	}

	return currentState, nil
}
