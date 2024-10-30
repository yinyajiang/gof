package ofdrm

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof"
)

func GenAuthHeader(urlpath string, auth gof.AuthInfo, rules Rules) map[string]string {
	timestamp := time.Now().UTC().UnixMilli()
	hashBytes := sha1.Sum([]byte(strings.Join([]string{rules.StaticParam, fmt.Sprintf("%d", timestamp), urlpath, auth.UserID}, "\n")))
	hashString := strings.ToLower(hex.EncodeToString(hashBytes[:]))
	checksum := slice.Reduce(rules.ChecksumIndexes, func(_ int, number int, accumulator int) int {
		return accumulator + int(hashString[number])
	}, 0) + rules.ChecksumConstant
	sign := rules.Prefix + ":" + hashString + ":" + strings.ToLower(fmt.Sprintf("%X", checksum)) + ":" + rules.Suffix
	header := map[string]string{
		"accept":     "application/json, text/plain",
		"app-token":  rules.AppToken,
		"cookie":     auth.Cookie,
		"sign":       sign,
		"time":       fmt.Sprintf("%d", timestamp),
		"user-id":    auth.UserID,
		"user-agent": auth.UserAgent,
		"x-bc":       auth.X_BC,
	}
	return header
}
