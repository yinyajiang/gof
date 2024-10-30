package gof

type AuthInfo struct {
	UserID    string
	UserAgent string
	X_BC      string
	Cookie    string
}

type VideoMPDInfo struct {
	MPDURL    string
	Policy    string
	Signature string
	KeyPairID string
	MediaID   string
	PostID    string
}
