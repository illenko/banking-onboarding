package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gittub.com/illenko/onboarding-service/internal/model"
	"gittub.com/illenko/onboarding-service/internal/queue"
	"gittub.com/illenko/onboarding-service/internal/worker"
	"gittub.com/illenko/onboarding-service/internal/workflow"
	httpModel "gittub.com/illenko/onboarding-service/pkg/http"
	"go.temporal.io/sdk/client"
)

func main() {
	worker.Run()

	router := gin.Default()

	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client:", err)
	}
	defer temporalClient.Close()

	router.POST("/onboarding", handleOnboarding(temporalClient))
	router.POST("/onboarding/:id/signature", handleSignature(temporalClient))
	router.GET("/onboarding/:id", handleGetOnboarding(temporalClient))

	err = router.Run(":8080")
	if err != nil {
		return
	}
}

func handleOnboarding(temporalClient client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New()

		var input model.UserRequest
		err := c.Bind(&input)
		if err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
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

			var result workflow.CurrentState

			err = we.Get(context.Background(), &result)

			if err != nil {
				log.Fatalln("Unable to get the Workflow result:", err)
			}
		}()

		c.JSON(http.StatusOK, httpModel.OnboardingStatus{
			ID:    id,
			State: "processing",
		})
	}
}

func handleSignature(temporalClient client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
			return
		}

		var input workflow.SignatureSignal
		err = c.Bind(&input)
		if err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
			return
		}

		err = temporalClient.SignalWorkflow(context.Background(), id.String(), "", "signature-signal", input)
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Unable to signal the Workflow")
			return
		}

		c.JSON(http.StatusOK, httpModel.OnboardingStatus{
			ID:    id,
			State: "processing",
		})
	}
}

func handleGetOnboarding(temporalClient client.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))
		if err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
			return
		}

		queryType := "current_state"

		response, err := temporalClient.QueryWorkflow(context.Background(), id.String(), "", queryType)
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Unable to query the Workflow")
			return
		}

		var currentState workflow.CurrentState
		err = response.Get(&currentState)
		if err != nil {
			sendErrorResponse(c, http.StatusInternalServerError, "Unable to get the current state")
			return
		}

		c.JSON(http.StatusOK, httpModel.OnboardingStatus{
			ID:    id,
			State: currentState.State,
			Data:  currentState.Data,
		})
	}
}

func sendErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, model.ErrorResponse{
		Code:    "bad_request",
		Message: message,
	})
}
