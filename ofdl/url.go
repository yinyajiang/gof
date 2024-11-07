package ofdl

import (
	"regexp"
	"strings"

	"github.com/yinyajiang/gof"
	"github.com/yinyajiang/gof/common"
)

var (
	reHome                   = _mustCompile(gof.OFPostDomain + "$")
	reSubscriptions          = _mustCompile(`/my/collections/user-lists/(?:subscribers|subscriptions|restricted|blocked)`)
	reChat                   = _mustCompile(`/my/chats(?:/chat/(?P<ID>[0-9]+))?$`)
	reUserList               = _mustCompile(`/my/collections/user-lists(?:/(?P<ID>[0-9]+))?$`)
	reSinglePost             = _mustCompile(`/(?P<PostID>[0-9]+)/(?P<UserName>[A-Za-z0-9\.\-_]+)$`)
	reUserWithMediaType      = _mustCompile(`/(?P<UserName>[A-Za-z0-9\.\-_]+)(?:/(?P<MediaType>media|videos|photos))?$`)
	reBookmarksWithMediaType = _mustCompile(`/my/collections/bookmarks(?:/(?:all|(?P<ID>[0-9]+))(?:/(?P<MediaType>photos|videos|audios|other|locked))?)?$`)
)

func _mustCompile(rePath string) *regexp.Regexp {
	var re string
	if !strings.HasPrefix(rePath, gof.OFPostDomain) {
		re = `(?i)` + regexp.QuoteMeta(gof.OFPostDomain) + rePath
	} else if rePath == gof.OFPostDomain {
		re = `(?i)` + regexp.QuoteMeta(gof.OFPostDomain)
	} else {
		re = `(?i)` + rePath
	}
	return regexp.MustCompile(re)
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

func ofurlFinds(must, optional []string, url string, res ...*regexp.Regexp) ([]string, bool) {
	if len(res) == 0 {
		return nil, false
	}

	url = common.CorrectOFURL(url, true)

reloop:
	for _, re := range res {
		if !re.MatchString(url) {
			continue
		}
		result := make([]string, len(must)+len(optional))

		if len(must) == 0 && len(optional) == 0 {
			return result, true
		}

		m, ok := common.ReGroup(re, url)
		if ok {
			for i, mustKey := range must {
				v, ok := m[mustKey]
				if ok {
					result[i] = v
				} else {
					continue reloop
				}
			}
			for i, optionalKey := range optional {
				v, ok := m[optionalKey]
				if ok {
					result[len(must)+i] = v
				}
			}
			return result, true
		} else if len(must) == 0 {
			return result, true
		}
	}
	return nil, false
}
