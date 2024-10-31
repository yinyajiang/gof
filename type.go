package gof

type AuthInfo struct {
	UserID    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	X_BC      string `json:"x_bc"`
	Cookie    string `json:"cookie"`
}

type MPDURLInfo struct {
	MPDURL    string
	Policy    string
	Signature string
	KeyPairID string
	MediaID   string
	PostID    string
}

type PostURLInfo struct {
	PostID   string
	UserName string
}

type Rules struct {
	AppToken         string `json:"app-token"`
	ChecksumConstant int    `json:"checksum_constant"`
	ChecksumIndexes  []int  `json:"checksum_indexes"`
	Prefix           string `json:"prefix"`
	StaticParam      string `json:"static_param"`
	Suffix           string `json:"suffix"`
}
