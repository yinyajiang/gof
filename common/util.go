package common

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/duke-git/lancet/v2/fileutil"
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
