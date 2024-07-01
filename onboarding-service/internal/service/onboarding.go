package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/illenko/onboarding-common/input"
	"github.com/illenko/onboarding-common/output"
	"github.com/illenko/onboarding-common/query"
	"github.com/illenko/onboarding-common/queue"
	"github.com/illenko/onboarding-common/signal"
	"github.com/illenko/onboarding-common/state"
	"github.com/illenko/onboarding-common/workflow"
	"github.com/illenko/onboarding-service/pkg/http"
	"go.temporal.io/sdk/client"
)

type OnboardingService struct {
	temporalClient client.Client
}

func NewOnboardingService(temporalClient client.Client) *OnboardingService {
	return &OnboardingService{temporalClient: temporalClient}
}

func (s *OnboardingService) CreateOnboarding(ctx context.Context, req *http.OnboardingRequest) (http.OnboardingStatus, error) {
	workflowId := uuid.New()

	options := client.StartWorkflowOptions{
		ID:        workflowId.String(),
		TaskQueue: queue.OnboardingTask,
	}

	_, err := s.temporalClient.ExecuteWorkflow(ctx, options, workflow.Onboarding, input.Onboarding{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		City:        req.City,
		AccountType: req.AccountType,
		Currency:    req.Currency,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Unable to start the Workflow:", slog.String("Error", err.Error()))

		return http.OnboardingStatus{
			ID:    workflowId,
			State: state.FailedState,
		}, nil

	}

	return http.OnboardingStatus{
		ID:    workflowId,
		State: state.ProcessingState,
	}, nil
}

func (s *OnboardingService) GetOnboarding(ctx context.Context, id uuid.UUID) (http.OnboardingStatus, error) {
	currentState, err := s.getCurrentState(ctx, id)

	if err != nil {
		return http.OnboardingStatus{}, err
	}

	return http.OnboardingStatus{
		ID:    id,
		State: currentState.State,
		Data:  currentState.Data,
	}, nil
}

func (s *OnboardingService) VerifySignature(ctx context.Context, id uuid.UUID, req *http.SignatureRequest) (http.OnboardingStatus, error) {
	currentState, err := s.getCurrentState(ctx, id)

	if err != nil {
		return http.OnboardingStatus{}, err
	}

	if currentState.State != state.WaitingForAgreementSignatureState {
		slog.WarnContext(ctx, "Invalid state for signature verification")
		return http.OnboardingStatus{
			ID:    id,
			State: currentState.State,
			Data:  currentState.Data,
		}, nil
	}

	signatureSignal := signal.Signature{
		Signature: req.Signature,
	}

	err = s.temporalClient.SignalWorkflow(ctx, id.String(), "", signal.SignatureSignal, signatureSignal)
	if err != nil {
		return http.OnboardingStatus{}, err
	}

	return http.OnboardingStatus{
		ID:    id,
		State: state.ProcessingState,
	}, nil
}

func (s *OnboardingService) getCurrentState(ctx context.Context, id uuid.UUID) (output.Onboarding, error) {
	response, err := s.temporalClient.QueryWorkflow(ctx, id.String(), "", query.CurrentState)
	if err != nil {
		return output.Onboarding{}, err
	}

	var currentState output.Onboarding
	err = response.Get(&currentState)
	if err != nil {
		return output.Onboarding{}, err
	}

	return currentState, nil
}
