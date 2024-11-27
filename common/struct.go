package common

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/yinyajiang/gof"
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

func MaybeDrmURL(u string) bool {
	return strings.Contains(u, gof.OFDrmMaybe)
}
