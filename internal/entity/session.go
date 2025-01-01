package entity

type Session struct {
	ID        string `redis:"id"`
	UserID    int    `redis:"user_id"`
	UserAgent string `redis:"user_agent"`
	UserIP    string `redis:"user_ip"`
}

// func (s Session) MarshalBinary() ([]byte, error) {
// 	return json.Marshal(s)
// }
