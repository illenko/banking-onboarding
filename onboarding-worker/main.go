package main

import (
	"log/slog"

	"github.com/illenko/onboarding-common/activity"
	"github.com/illenko/onboarding-common/configuration"
	"github.com/illenko/onboarding-common/queue"
	"github.com/illenko/onboarding-common/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	configuration.LoadEnv()

	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		slog.Error("Unable to create Temporal client:", slog.String("error", err.Error()))
		return
	}
	defer temporalClient.Close()

	w := worker.New(temporalClient, queue.OnboardingTask, worker.Options{})

	w.RegisterWorkflow(workflow.Onboarding)
	w.RegisterActivity(activity.AntifraudChecks)
	w.RegisterActivity(activity.CreateUser)
	w.RegisterActivity(activity.CreateAccount)
	w.RegisterActivity(activity.CreateAgreement)
	w.RegisterActivity(activity.ValidateSignature)
	w.RegisterActivity(activity.CreateCard)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		slog.Error("Unable to start worker", slog.String("error", err.Error()))
		return
	}
}
