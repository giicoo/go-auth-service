package models

type UserCreateRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDeleteRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateRequest struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserGetByEmailRequest struct {
	Email string `json:"email"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
