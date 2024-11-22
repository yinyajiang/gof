package common

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
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

func FileUnmarshal(file string, v any) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func FileMarshal(file string, v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	fileutil.CreateDir(filepath.Dir(file))
	return os.WriteFile(file, data, 0644)
}

func WriteFile(file string, data []byte) error {
	fileutil.CreateDir(filepath.Dir(file))
	return os.WriteFile(file, data, 0644)
}

func ForeachCookie(cookie string, cb func(key, value string) bool) {
	if cookie == "" || cb == nil {
		return
	}
	for _, s := range strings.Split(cookie, ";") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		kv := strings.SplitN(s, "=", 2)
		if len(kv) != 2 {
			continue
		}
		if !cb(strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])) {
			break
		}
	}
}

func FindCookie(cookie string, key string) string {
	var value string
	ForeachCookie(cookie, func(k, v string) bool {
		if strings.EqualFold(k, key) {
			value = v
			return false
		}
		return true
	})
	return value
}

func CleanEmptryString(arr []string) []string {
	return slice.Filter(arr, func(_ int, s string) bool {
		return s != ""
	})
}
