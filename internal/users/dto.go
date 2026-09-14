package users

type LoginRequestDTO struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponseDTO struct {
	Token    string `json:"token"`
	Username string `json:"username,omitempty"`
}
