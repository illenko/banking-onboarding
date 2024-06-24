package output

import "github.com/illenko/onboarding-common/state"

type Onboarding struct {
	State state.OnboardingState
	Data  map[string]any
}
