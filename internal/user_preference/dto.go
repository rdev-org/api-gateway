package user_preference

type PreferenceRequestDTO struct {
	UserID     string            `json:"user_id"`
	Language   string            `json:"language"`
	Currency   string            `json:"currency"`
	EmailOptIn bool              `json:"email_opt_in"`
	SmsOptIn   bool              `json:"sms_opt_in"`
	Settings   map[string]string `json:"settings"`
}
