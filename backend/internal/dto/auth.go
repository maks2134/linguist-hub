package dto

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname,omitempty"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}
