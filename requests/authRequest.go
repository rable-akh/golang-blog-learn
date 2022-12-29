package requests

type SignUpRequest struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type SignInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
