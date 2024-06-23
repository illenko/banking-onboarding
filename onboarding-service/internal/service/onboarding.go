package service

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/illenko/onboarding-service/internal/input"
	"github.com/illenko/onboarding-service/internal/output"
	"github.com/illenko/onboarding-service/internal/query"
	"github.com/illenko/onboarding-service/internal/queue"
	"github.com/illenko/onboarding-service/internal/signal"
	"github.com/illenko/onboarding-service/internal/workflow"
	"github.com/illenko/onboarding-service/pkg/http"
	"github.com/illenko/onboarding-service/pkg/state"
	"go.temporal.io/sdk/client"
)

type OnboardingService interface {
	CreateOnboarding(ctx context.Context, req *http.OnboardingRequest) (*http.OnboardingStatus, error)
	GetOnboarding(ctx context.Context, id uuid.UUID) (*http.OnboardingStatus, error)
	VerifySignature(ctx context.Context, id uuid.UUID, req *http.SignatureRequest) (*http.OnboardingStatus, error)
}

type OnboardingServiceImpl struct {
	temporalClient client.Client
}

func NewOnboardingService(temporalClient client.Client) *OnboardingServiceImpl {
	return &OnboardingServiceImpl{temporalClient: temporalClient}
}

func (s *OnboardingServiceImpl) CreateOnboarding(ctx context.Context, req *http.OnboardingRequest) (*http.OnboardingStatus, error) {
	id := uuid.New()

	onboardingInput := input.Onboarding{
		ID:        id,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		City:      req.City,
	}

	go s.executeOnboardingWorkflow(ctx, id.String(), onboardingInput)

	return &http.OnboardingStatus{
		ID:    id,
		State: state.ProcessingState,
	}, nil
}

func (s *OnboardingServiceImpl) executeOnboardingWorkflow(ctx context.Context, workflowID string, onboardingInput input.Onboarding) {
	options := client.StartWorkflowOptions{
		ID:        workflowID,
		TaskQueue: queue.OnboardingTask,
	}

	we, err := s.temporalClient.ExecuteWorkflow(ctx, options, workflow.Onboarding, onboardingInput)
	if err != nil {
		log.Println("Unable to start the Workflow:", err)
		return
	}

	log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

	var result output.Onboarding
	err = we.Get(ctx, &result)

	if err != nil {
		log.Println("Unable to get the Workflow result:", err)
		return
	}
}

func (s *OnboardingServiceImpl) GetOnboarding(ctx context.Context, id uuid.UUID) (*http.OnboardingStatus, error) {
	response, err := s.temporalClient.QueryWorkflow(context.Background(), id.String(), "", query.CurrentState)
	if err != nil {
		return nil, err
	}

	var currentState output.Onboarding
	err = response.Get(&currentState)
	if err != nil {
		return nil, err
	}

	return &http.OnboardingStatus{
		ID:    id,
		State: currentState.State,
		Data:  currentState.Data,
	}, nil
}

func (s *OnboardingServiceImpl) VerifySignature(ctx context.Context, id uuid.UUID, req *http.SignatureRequest) (*http.OnboardingStatus, error) {
	signatureSignal := signal.Signature{
		Signature: req.Signature,
	}

	err := s.temporalClient.SignalWorkflow(context.Background(), id.String(), "", signal.SignatureSignal, signatureSignal)
	if err != nil {
		return nil, err
	}

	return &http.OnboardingStatus{
		ID:    id,
		State: state.ProcessingState,
	}, nil
}
