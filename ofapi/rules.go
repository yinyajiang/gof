package ofapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/fileutil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof/common"
)

func loadRules(cacheDir string, rulesURL []string) (rules, error) {
	var allRules []rules

	urlRules, urlErr := _loadURLRules(rulesURL)
	if urlErr == nil {
		allRules = append(allRules, urlRules)
	}
	cachedRules, cachedErr := _loadCachedRules(cacheDir)
	if cachedErr == nil {
		allRules = append(allRules, cachedRules)
	}
	if cachedErr != nil && urlErr != nil {
		return rules{}, urlErr
	}

	latest := _selectLatestRules(allRules)
	if cachedRules.Revision != latest.Revision {
		_cacheRules(cacheDir, latest)
	}
	return latest, nil
}

func _isValidRules(rules *rules) bool {
	if rules == nil {
		return false
	}
	if rules.AppToken_Old != "" && rules.AppToken == "" {
		rules.AppToken = rules.AppToken_Old
	}
	return rules.AppToken != "" &&
		rules.ChecksumConstant != 0 &&
		len(rules.ChecksumIndexes) > 0 &&
		rules.Prefix != "" &&
		rules.StaticParam != "" &&
		rules.Suffix != ""
}

func _cacheRules(cacheDir string, rules rules) {
	data, err := json.Marshal(rules)
	if err != nil {
		fmt.Printf("marshal rules failed, err: %v\n", err)
		return
	}
	os.WriteFile(filepath.Join(cacheDir, "rules"), data, 0644)
}

func _loadCachedRules(cacheDir string) (rules, error) {
	data, err := os.ReadFile(filepath.Join(cacheDir, "rules"))
	if err != nil {
		return rules{}, err
	}
	var rules rules
	err = json.Unmarshal(data, &rules)
	return rules, err
}

func _loadURLRules(rulesURI []string) (rules, error) {
	const fixURL = "https://git.ofdl.tools/sim0n00ps/dynamic-rules/raw/branch/main/rules.json"

	if len(rulesURI) == 0 {
		rulesURI = []string{fixURL}
	} else {
		if !slice.Contain(rulesURI, fixURL) {
			rulesURI = append(rulesURI, fixURL)
		}
	}

	pruleList := make([]*rules, len(rulesURI))
	wg := sync.WaitGroup{}
	for i, url := range rulesURI {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			var rules rules

			var e error
			if !strings.HasPrefix(url, "http") && fileutil.IsExist(url) {
				e = common.FileUnmarshal(url, &rules)
			} else {
				e = common.HttpGetUnmarshal(url, &rules)
			}
			if e != nil {
				fmt.Printf("get rules from %s failed, err: %v\n", url, e)
			} else {
				pruleList[i] = &rules
			}
		}(i, url)
	}
	wg.Wait()

	ruleList := []rules{}
	for _, item := range pruleList {
		if item != nil && _isValidRules(item) {
			ruleList = append(ruleList, *item)
		}
	}
	if len(ruleList) == 0 {
		return rules{}, errors.New("no url valid rules")
	}
	return _selectLatestRules(ruleList), nil
}

func _selectLatestRules(rulesList []rules) rules {
	if len(rulesList) == 0 {
		return rules{}
	}

	var latestRules rules
	var latestRevisionTime int64 = -1
	for _, item := range rulesList {
		if !_isValidRules(&item) {
			continue
		}
		if item.Revision == "" {
			item.Revision = time.Now().Format("202310311103") + "-000000"
		}
		revision := strings.Split(item.Revision, "-")[0]
		revisionTime, e := strconv.ParseInt(revision, 10, 64)
		if e == nil && revisionTime > latestRevisionTime {
			latestRevisionTime = revisionTime
			latestRules = item
		}
	}
	return latestRules
}
