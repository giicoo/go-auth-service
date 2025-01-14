package models

type SessionCreateRequest struct {
	UserID    int    `json:"user_id"`
	UserAgent string `json:"user_agent"`
	UserIP    string `json:"user_ip"`
}

type SessionRequest struct {
	ID string `json:"session_id"`
}

type SessionListRequest struct {
	UserID int `json:"user_id"`
}

type SessionResponse struct {
	ID string `json:"session_id"`
}

type SessionResponseFull struct {
	ID        string `json:"session_id"`
	UserID    int    `json:"user_id"`
	UserAgent string `json:"user_agent"`
	UserIP    string `json:"user_ip"`
}
