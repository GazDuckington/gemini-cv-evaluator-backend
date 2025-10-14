package dto

type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type ProfileResponse struct {
	ID        string `json:"id"`
	FullName  string `json:"full_name,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Phone     string `json:"phone,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type UserResponse struct {
	ID      string           `json:"id"`
	Email   string           `json:"email"`
	Role    string           `json:"role"`
	Profile *ProfileResponse `json:"profile,omitempty"`
}
