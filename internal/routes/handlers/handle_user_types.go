package handlers

type UserSignupRequest struct {
	Username string `json:"username" validate:"required,max=255,alphanum"`
	Password string `json:"password" validate:"required,max=72"`
}

type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	AccessToken string `json:"access_token"`
}
