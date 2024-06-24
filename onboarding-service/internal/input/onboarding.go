package input

type Onboarding struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	City        string `json:"city"`
	AccountType string `json:"account_type"`
	Currency    string `json:"currency"`
}
