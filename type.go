package gof

type AuthInfo struct {
	UserID    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	X_BC      string `json:"x_bc"`
	Cookie    string `json:"cookie"`
}

type PostURLInfo struct {
	PostID   string
	UserName string
}
