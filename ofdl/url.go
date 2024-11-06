package ofdl

import (
	"regexp"
	"strings"

	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
)

var (
	reSubscriptions = _mustCompile(`/my/collections/user-lists/(?:subscribers|subscriptions)(?:/active)?`)
	reSingleChat    = _mustCompile(`/my/chats/chat/(?P<ID>[0-9]+)$`)
	reChats         = _mustCompile(`/my/chats$`)
	reUserList      = _mustCompile(`/my/collections/user-lists/(?P<ID>[0-9]+)$`)
	reSinglePost    = _mustCompile(`/(?P<PostID>[0-9]+)/(?P<UserName>[A-Za-z0-9\.\-_]+)$`)
	reUser          = _mustCompile(`/(?P<UserName>[A-Za-z0-9\.\-_]+)$`)
)

func _mustCompile(rePath string) *regexp.Regexp {
	re := `(?i)` + regexp.QuoteMeta(gof.OFPostDomain) + rePath
	return regexp.MustCompile(re)
}

func isOFHomeURL(url string) bool {
	url = common.CorrectOFURL(url, true)
	if i := strings.Index(url, "?"); i != -1 {
		url = url[:i]
	}
	return strings.TrimRight(url, "/") == strings.TrimRight(gof.OFPostDomain, "/")
}

func isOFURL(url string) bool {
	url = common.CorrectOFURL(url, true)
	return strings.HasPrefix(url, gof.OFPostDomain)
}

func ofurlMatch(re *regexp.Regexp, url string) bool {
	url = common.CorrectOFURL(url, true)
	return re.MatchString(url)
}

func ofurlFind(re *regexp.Regexp, url, key string) (string, bool) {
	url = common.CorrectOFURL(url, true)
	if m, ok := common.ReGroup(re, url); ok {
		v, ok := m[key]
		return v, ok
	}
	return "", false
}

func ofurlFind2(re *regexp.Regexp, url, key1, key2 string) (string, string, bool) {
	url = common.CorrectOFURL(url, true)
	if m, ok := common.ReGroup(re, url); ok {
		v1, ok1 := m[key1]
		v2, ok2 := m[key2]
		return v1, v2, ok1 && ok2
	}
	return "", "", false
}
