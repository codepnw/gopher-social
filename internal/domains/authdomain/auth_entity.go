package authdomain

type RegisterUserPayload struct {
	Username string `json:"username" binding:"required,max=100"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

type LoginUserPayload struct {
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}