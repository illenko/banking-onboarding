package state

type OnboardingState string

const (
	ProcessingState                   OnboardingState = "processing"
	FailedState                       OnboardingState = "failed"
	FraudNotPassedState               OnboardingState = "fraud_not_passed"
	SignatureNotValidSate             OnboardingState = "signature_not_valid"
	WaitingForAgreementSignatureState OnboardingState = "waiting_for_agreement_signature"
	CompletedState                    OnboardingState = "completed"
)
