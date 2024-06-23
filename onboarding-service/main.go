package main

import (
	"context"
	"github.com/illenko/onboarding-service/internal/configuration"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/illenko/onboarding-service/internal/input"
	"github.com/illenko/onboarding-service/internal/output"
	"github.com/illenko/onboarding-service/internal/queue"
	"github.com/illenko/onboarding-service/internal/signal"
	"github.com/illenko/onboarding-service/internal/worker"
	"github.com/illenko/onboarding-service/internal/workflow"
	httpModel "github.com/illenko/onboarding-service/pkg/http"
	"go.temporal.io/sdk/client"
)

func main() {
	configuration.LoadEnv()

	go worker.Run()

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

		var request httpModel.OnboardingRequest
		err := c.Bind(&request)
		if err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
			return
		}

		options := client.StartWorkflowOptions{
			ID:        id.String(),
			TaskQueue: queue.OnboardingTask,
		}

		onboardingInput := input.Onboarding{
			ID:        id,
			FirstName: request.FirstName,
			LastName:  request.LastName,
			Email:     request.Email,
			City:      request.City,
		}

		go func() {
			we, err := temporalClient.ExecuteWorkflow(context.Background(), options, workflow.Onboarding, onboardingInput)
			if err != nil {
				log.Fatalln("Unable to start the Workflow:", err)
			}

			log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

			var result output.Onboarding

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

		var request httpModel.SignatureRequest
		err = c.Bind(&request)
		if err != nil {
			sendErrorResponse(c, http.StatusBadRequest, "Invalid request")
			return
		}

		signatureSignal := signal.Signature{
			Signature: request.Signature,
		}

		err = temporalClient.SignalWorkflow(context.Background(), id.String(), "", signal.SignatureSignal, signatureSignal)
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

		var currentState output.Onboarding
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
	c.JSON(statusCode, httpModel.ErrorResponse{
		Code:    "bad_request",
		Message: message,
	})
}
