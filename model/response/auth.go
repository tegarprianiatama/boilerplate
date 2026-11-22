package response

import "time"

type LoginResp struct {
	Token     string    `json:"token"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
	User      UserResp  `json:"user"`
}
