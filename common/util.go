package common

import (
	"regexp"
	"strings"

	"github.com/yinyajiang/gof"
)

func ReGroup(re *regexp.Regexp, s string) (map[string]string, bool) {
	matches := re.FindStringSubmatch(s)
	if len(matches) == 0 {
		return nil, false
	}
	names := re.SubexpNames()
	result := make(map[string]string)
	for i, name := range names {
		if i != 0 && name != "" {
			result[name] = matches[i]
		}
	}
	return result, true
}

func CorrectOFURL(url string, removeQuery bool) string {
	if !strings.HasPrefix(url, "https://") {
		url = "https://" + url
	}
	if !strings.Contains(gof.OFPostDomain, "www.") {
		url = strings.Replace(strings.TrimSpace(url), "www.", "", 1)
	}
	if removeQuery {
		if i := strings.Index(url, "?"); i != -1 {
			url = url[:i]
		}
		url = strings.TrimRight(url, "/")
	}
	return url
}
