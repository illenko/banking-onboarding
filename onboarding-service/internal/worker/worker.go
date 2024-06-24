package worker

import (
	"log/slog"

	"github.com/illenko/onboarding-service/internal/activity"
	"github.com/illenko/onboarding-service/internal/queue"
	"github.com/illenko/onboarding-service/internal/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func Run(temporalClient client.Client) {

	w := worker.New(temporalClient, queue.OnboardingTask, worker.Options{})

	w.RegisterWorkflow(workflow.Onboarding)
	w.RegisterActivity(activity.AntifraudChecks)
	w.RegisterActivity(activity.CreateUser)
	w.RegisterActivity(activity.CreateAccount)
	w.RegisterActivity(activity.CreateAgreement)
	w.RegisterActivity(activity.ValidateSignature)
	w.RegisterActivity(activity.CreateCard)

	err := w.Run(worker.InterruptCh())
	if err != nil {
		slog.Error("Unable to start worker", slog.String("error", err.Error()))
		return
	}

}
