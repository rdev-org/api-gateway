package users

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Token    string `json:"token"`
	Username string `json:"username,omitempty"`
}

type RegisterRequestDTO struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Mobile       string `json:"mobile,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	MiddleName   string `json:"middle_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}
