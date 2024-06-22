package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gittub.com/illenko/onboarding-service/internal/activity"
	"gittub.com/illenko/onboarding-service/internal/model"
	"gittub.com/illenko/onboarding-service/internal/queue"
	"gittub.com/illenko/onboarding-service/internal/workflow"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"net/http"
)

func main() {

	go startWorker()

	router := gin.Default()

	temporalClient, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer temporalClient.Close()

	options := client.StartWorkflowOptions{
		ID:        uuid.New().String(),
		TaskQueue: queue.OnboardingTask,
	}

	router.POST("/onboarding", func(c *gin.Context) {

		var input model.UserRequest
		err := c.Bind(&input)

		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Code:    "bad_request",
				Message: "Invalid request",
			})
			return
		}

		we, err := temporalClient.ExecuteWorkflow(context.Background(), options, workflow.Onboarding, input)
		if err != nil {
			log.Fatalln("Unable to start the Workflow:", err)
		}

		log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

		var result string

		err = we.Get(context.Background(), &result)

		if err != nil {
			log.Fatalln("Unable to get Workflow result:", err)
		}
	})

	router.Run(":8080")
}

func startWorker() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client.", err)
	}
	defer c.Close()

	w := worker.New(c, queue.OnboardingTask, worker.Options{})

	w.RegisterWorkflow(workflow.Onboarding)
	w.RegisterActivity(activity.AntifraudChecks)

	// Start listening to the Task Queue.
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("unable to start Worker", err)
	}
}
