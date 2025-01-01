package models

type SessionCreateRequest struct {
	UserID    int    `json:"user_id"`
	UserAgent string `json:"user_agent"`
	UserIP    string `json:"user_ip"`
}

type SessionGetRequest struct {
	ID     string `json:"id"`
	UserID int    `json:"user_id"`
}

type SessionDeleteRequest struct {
	ID     string `json:"id"`
	UserID int    `json:"user_id"`
}

type SessionGetListRequest struct {
	UserID int `json:"user_id"`
}

type SessionDeleteListRequest struct {
	UserID int `json:"user_id"`
}

type SessionResponse struct {
	ID string `json:"id"`
}

type SessionResponseFull struct {
	ID        string `json:"id"`
	UserID    int    `json:"user_id"`
	UserAgent string `json:"user_agent"`
	UserIP    string `json:"user_ip"`
}
