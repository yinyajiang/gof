package common

import (
	"encoding/json"
	"fmt"
)

func MustMarshalJSON(v any) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		fmt.Println(err)
		return []byte{}
	}
	return b
}

func MustUnmarshalJSONStr(v any) string {
	return string(MustMarshalJSON(v))
}
