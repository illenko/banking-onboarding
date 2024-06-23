package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gittub.com/illenko/onboarding-service/internal/model"
	"gittub.com/illenko/onboarding-service/internal/queue"
	"gittub.com/illenko/onboarding-service/internal/worker"
	"gittub.com/illenko/onboarding-service/internal/workflow"
	httpModel "gittub.com/illenko/onboarding-service/pkg/http"
	"go.temporal.io/sdk/client"
	"log"
	"net/http"
)

func main() {

	go worker.Run()

	router := gin.Default()

	temporalClient, err := client.Dial(client.Options{})

	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}

	defer temporalClient.Close()

	router.POST("/onboarding", func(c *gin.Context) {

		id := uuid.New()

		var input model.UserRequest
		err := c.Bind(&input)

		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Code:    "bad_request",
				Message: "Invalid request",
			})
			return
		}

		options := client.StartWorkflowOptions{
			ID:        id.String(),
			TaskQueue: queue.OnboardingTask,
		}

		go func() {
			we, err := temporalClient.ExecuteWorkflow(context.Background(), options, workflow.Onboarding, input)
			if err != nil {
				log.Fatalln("Unable to start the Workflow:", err)
			}

			log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

			var result httpModel.OnboardingResponse

			err = we.Get(context.Background(), &result)
		}()

		c.JSON(http.StatusOK, httpModel.OnboardingStatus{
			ID:    id,
			State: "processing",
		})
	})

	router.POST("/onboarding/:id/signature", func(c *gin.Context) {

		id, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Code:    "bad_request",
				Message: "Invalid request",
			})
			return
		}

		var input workflow.SignatureSignal
		err = c.Bind(&input)

		if err != nil {
			c.JSON(http.StatusBadRequest, model.ErrorResponse{
				Code:    "bad_request",
				Message: "Invalid request",
			})
			return
		}

		err = temporalClient.SignalWorkflow(context.Background(), id.String(), "", "signature-signal", input)

		if err != nil {
			c.JSON(http.StatusInternalServerError, model.ErrorResponse{
				Code:    "internal_error",
				Message: "Unable to signal the Workflow",
			})
			return
		}

		c.JSON(http.StatusOK, httpModel.OnboardingStatus{
			ID:    id,
			State: "processing",
		})
	})

	router.Run(":8080")
}
