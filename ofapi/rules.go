package ofapi

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/duke-git/lancet/v2/slice"
	"github.com/yinyajiang/gof/common"
)

func loadDynamicRules(rulesURL ...string) (rules, error) {
	const fixURL = "https://raw.githubusercontent.com/deviint/onlyfans-dynamic-rules/main/dynamicRules.json"

	if len(rulesURL) == 0 {
		rulesURL = []string{fixURL}
	} else {
		if !slice.Contain(rulesURL, fixURL) {
			rulesURL = append(rulesURL, fixURL)
		}
	}

	ruleList := make([]*rules, len(rulesURL))
	wg := sync.WaitGroup{}
	for i, url := range rulesURL {
		wg.Add(1)
		go func(i int, url string) {
			defer wg.Done()
			var rules rules
			err := common.HttpGetUnmarshalJson(url, &rules)
			if err != nil {
				fmt.Printf("get rules from %s failed, err: %v\n", url, err)
			} else {
				ruleList[i] = &rules
			}
		}(i, url)
	}
	wg.Wait()

	var latestRules *rules
	var latestRevisionTime int64 = -1
	for _, rules := range ruleList {
		if rules == nil || !isValidRules(*rules) {
			continue
		}
		if rules.Revision == "" {
			rules.Revision = time.Now().Format("202310311103") + "-000000"
		}
		revision := strings.Split(rules.Revision, "-")[0]
		revisionTime, e := strconv.ParseInt(revision, 10, 64)
		if e == nil && revisionTime > latestRevisionTime {
			latestRevisionTime = revisionTime
			latestRules = rules
		}
	}
	if latestRules == nil {
		return rules{}, fmt.Errorf("no valid rules found")
	}
	return *latestRules, nil
}

func isValidRules(rules rules) bool {
	return rules.AppToken != "" &&
		rules.ChecksumConstant != 0 &&
		len(rules.ChecksumIndexes) > 0 &&
		rules.Prefix != "" &&
		rules.StaticParam != "" &&
		rules.Suffix != ""

}
