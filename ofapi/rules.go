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

	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof/common"
)

func LoadRules(cacheDir string, rulesURL []string, cachePriority ...bool) (rules, error) {
	if len(cachePriority) > 0 && cachePriority[0] {
		cachedRules, e := loadCachedRules(cacheDir)
		if e == nil {
			return cachedRules, nil
		}
	}

	var allRules []rules

	urlRules, urlErr := loadURLRules(rulesURL)
	if urlErr == nil {
		allRules = append(allRules, urlRules)
	}
	cachedRules, cachedErr := loadCachedRules(cacheDir)
	if cachedErr == nil {
		allRules = append(allRules, cachedRules)
	}
	if cachedErr != nil && urlErr != nil {
		return rules{}, urlErr
	}

	latest := selectLatestRules(allRules)
	if cachedRules.Revision != latest.Revision {
		cacheRules(cacheDir, latest)
	}
	return latest, nil
}

func isValidRules(rules rules) bool {
	return rules.AppToken != "" &&
		rules.ChecksumConstant != 0 &&
		len(rules.ChecksumIndexes) > 0 &&
		rules.Prefix != "" &&
		rules.StaticParam != "" &&
		rules.Suffix != ""
}

func cacheRules(cacheDir string, rules rules) {
	data, err := json.Marshal(rules)
	if err != nil {
		fmt.Printf("marshal rules failed, err: %v\n", err)
		return
	}
	os.WriteFile(filepath.Join(cacheDir, "rules"), data, 0644)
}

func loadCachedRules(cacheDir string) (rules, error) {
	data, err := os.ReadFile(filepath.Join(cacheDir, "rules"))
	if err != nil {
		return rules{}, err
	}
	var rules rules
	err = json.Unmarshal(data, &rules)
	return rules, err
}

func loadURLRules(rulesURL []string) (rules, error) {
	const fixURL = "https://raw.githubusercontent.com/deviint/onlyfans-dynamic-rules/main/dynamicRules.json"

	if len(rulesURL) == 0 {
		rulesURL = []string{fixURL}
	} else {
		if !slice.Contain(rulesURL, fixURL) {
			rulesURL = append(rulesURL, fixURL)
		}
	}

	pruleList := make([]*rules, len(rulesURL))
	wg := sync.WaitGroup{}
	for i, url := range rulesURL {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			var rules rules
			e := common.HttpGetUnmarshal(url, &rules)
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
		if item != nil && isValidRules(*item) {
			ruleList = append(ruleList, *item)
		}
	}
	if len(ruleList) == 0 {
		return rules{}, errors.New("no url valid rules")
	}
	return selectLatestRules(ruleList), nil
}

func selectLatestRules(rulesList []rules) rules {
	if len(rulesList) == 0 {
		return rules{}
	}

	var latestRules rules
	var latestRevisionTime int64 = -1
	for _, item := range rulesList {
		if !isValidRules(item) {
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
