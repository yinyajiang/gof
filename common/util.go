package common

import (
	"encoding/json"

	"github.com/yinyajiang/gof"
)

func MustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func PanicAuthInfo(authInfo gof.AuthInfo) {
	if authInfo.Cookie == "" || authInfo.X_BC == "" || authInfo.UserAgent == "" {
		panic("AuthInfo is invalid")
	}
}
