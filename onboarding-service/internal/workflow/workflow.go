package workflow

import (
	"gittub.com/illenko/onboarding-service/internal/activity"
	"gittub.com/illenko/onboarding-service/internal/model"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

func Onboarding(ctx workflow.Context, input model.UserRequest) (string, error) {
	retryPolicy := &temporal.RetryPolicy{
		InitialInterval:        time.Second,
		BackoffCoefficient:     2.0,
		MaximumInterval:        100 * time.Second,
		MaximumAttempts:        500,
		NonRetryableErrorTypes: []string{"FraudCheckError", "InvalidSignatureError"},
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
		return "", antifraudErr
	}

	return "", nil
}
