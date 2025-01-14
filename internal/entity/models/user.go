package models

type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateEmailRequest struct {
	UserID int    `json:"id"`
	Email  string `json:"new_email"`
}

type UserUpdatePasswordRequest struct {
	UserID   int    `json:"id"`
	Password string `json:"new_password"`
}

type UserIdRequest struct {
	UserID int `json:"id"`
}
type UserCheckRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	UserID int    `json:"id"`
	Email  string `json:"email"`
}
