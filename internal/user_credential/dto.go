package user_credential

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Token    string `json:"token"`
	Username string `json:"username,omitempty"`
	Email    string `json:"email"`
	Uuid     string `json:"uuid"`
}

type RegisterRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PasswordRequestDTO struct {
	Password string `json:"password"`
}

type ForgotPasswordRequestDTO struct {
	Email string `json:"email"`
}
