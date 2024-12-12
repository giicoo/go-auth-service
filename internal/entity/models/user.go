package models

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type UserCreateRequest struct {
	UserRequest
}

type UserDeleteRequest struct {
	UserRequest
}

type UserUpdateRequest struct {
	ID int `json:"id"`
	UserRequest
}

type UserGetByEmailRequest struct {
	Email string `json:"email"`
}

type UserResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
