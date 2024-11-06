package ofdl

import (
	"regexp"
	"strings"

	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
	"github.com/yinyajiang/gof/ofapi"
)

var (
	reSubscriptions             = _mustCompile(`/my/collections/user-lists/(?:subscribers|subscriptions)(?:/active)?`)
	reSingleChat                = _mustCompile(`/my/chats/chat/(?P<ID>[0-9]+)$`)
	reChats                     = _mustCompile(`/my/chats$`)
	reUserList                  = _mustCompile(`/my/collections/user-lists/(?P<ID>[0-9]+)$`)
	reSinglePost                = _mustCompile(`/(?P<PostID>[0-9]+)/(?P<UserName>[A-Za-z0-9\.\-_]+)$`)
	reUser                      = _mustCompile(`/(?P<UserName>[A-Za-z0-9\.\-_]+)$`)
	reUserByMediaType           = _mustCompile(`/(?P<UserName>[A-Za-z0-9\.\-_]+)/(?P<MediaType>media|videos|photos)$`)
	reAllBookmarks              = _mustCompile(`/my/collections/bookmarks(?:/all)?$`)
	reAllBookmarksByMediaType   = _mustCompile(`/my/collections/bookmarks/all/(?P<MediaType>photos|videos|audios|other|locked)$`)
	reSingleBookmark            = _mustCompile(`/my/collections/bookmarks/(?P<ID>[0-9]+)$`)
	reSingleBookmarkByMediaType = _mustCompile(`/my/collections/bookmarks/(?P<ID>[0-9]+)/(?P<MediaType>photos|videos|audios|other|locked)$`)
)

func bookmarkMediaType(s string) ofapi.BookmarkMedia {
	switch s {
	case "photos":
		return ofapi.BookmarkPhotos
	case "videos":
		return ofapi.BookmarkVideos
	case "audios":
		return ofapi.BookmarkAudios
	case "other":
		return ofapi.BookmarkOther
	case "locked":
		return ofapi.BookmarkLocked
	}
	return ofapi.BookmarkAll
}

func userMediasType(s string) ofapi.UserMedias {
	switch s {
	case "media":
		return ofapi.UserMediasAll
	case "videos":
		return ofapi.UserMediasVideos
	case "photos":
		return ofapi.UserMediasPhotos
	}
	return ofapi.UserMediasAll
}

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

func ofurlMatchs(url string, res ...*regexp.Regexp) bool {
	if len(res) == 0 {
		return false
	}
	url = common.CorrectOFURL(url, true)
	for _, re := range res {
		if re.MatchString(url) {
			return true
		}
	}
	return false
}

func ofurlFinds(url, key string, res ...*regexp.Regexp) (string, bool) {
	if len(res) == 0 {
		return "", false
	}
	url = common.CorrectOFURL(url, true)

	for _, re := range res {
		if m, ok := common.ReGroup(re, url); ok {
			v, ok := m[key]
			if ok {
				return v, ok
			}
		}
	}
	return "", false
}

func ofurlFinds2(url, key1, key2 string, res ...*regexp.Regexp) (string, string, bool) {
	if len(res) == 0 {
		return "", "", false
	}
	url = common.CorrectOFURL(url, true)
	for _, re := range res {
		if m, ok := common.ReGroup(re, url); ok {
			v1, ok1 := m[key1]
			v2, ok2 := m[key2]
			if ok1 && ok2 {
				return v1, v2, true
			}
		}
	}
	return "", "", false
}
