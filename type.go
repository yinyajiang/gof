package gof

type AuthInfo struct {
	UserID    string
	UserAgent string
	X_BC      string
	Cookie    string
}

type MPDInfo struct {
	MPDURL    string
	Policy    string
	Signature string
	KeyPairID string
	MediaID   string
	PostID    string
}

type Rules struct {
	AppToken         string `json:"app-token"`
	ChecksumConstant int    `json:"checksum_constant"`
	ChecksumIndexes  []int  `json:"checksum_indexes"`
	Prefix           string `json:"prefix"`
	StaticParam      string `json:"static_param"`
	Suffix           string `json:"suffix"`
}
