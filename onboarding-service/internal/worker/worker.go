package worker

import (
	"gittub.com/illenko/onboarding-service/internal/activity"
	"gittub.com/illenko/onboarding-service/internal/queue"
	"gittub.com/illenko/onboarding-service/internal/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
)

func Run() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, queue.OnboardingTask, worker.Options{})

	w.RegisterWorkflow(workflow.Onboarding)
	w.RegisterActivity(activity.AntifraudChecks)
	w.RegisterActivity(activity.CreateUser)
	w.RegisterActivity(activity.CreateAccount)
	w.RegisterActivity(activity.CreateAgreement)
	w.RegisterActivity(activity.CreateSignature)
	w.RegisterActivity(activity.CreateCard)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
