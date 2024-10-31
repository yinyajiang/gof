package common

import (
	"encoding/json"
)

func MustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
