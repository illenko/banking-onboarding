package output

import "github.com/illenko/onboarding-service/pkg/state"

type Onboarding struct {
	State state.OnboardingState `json:"state"`
	Data  map[string]any        `json:"data"`
}
