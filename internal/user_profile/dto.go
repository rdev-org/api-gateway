package user_profile

type ProfileRequestDTO struct {
	FirstName    string `json:"first_name"`
	MiddleName   string `json:"middle_name"`
	LastName     string `json:"last_name"`
	MobileNumber string `json:"mobile_number"`
	DateOfBirth  string `json:"date_of_birth"`
	Gender       string `json:"gender"`
	Avatar       string `json:"avatar"`
}
