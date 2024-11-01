package ofapi

import (
	"strconv"
	"time"
)

func initPublishTimeParam(addParam map[string]string, timePoint time.Time, timeDirection TimeDirection) map[string]string {
	if timePoint.IsZero() {
		timePoint = time.Now()
		timeDirection = TimeDirectionBefore
	}
	initTime := strconv.FormatInt(timePoint.Unix(), 10) + ".000000"

	if addParam == nil {
		addParam = make(map[string]string)
	}
	param := addParam
	if timeDirection == TimeDirectionBefore {
		param["beforePublishTime"] = initTime
	} else {
		param["afterPublishTime"] = initTime
	}
	return param
}

func updatePublishTimeParam(param map[string]string, timeDirection TimeDirection, moreMarker moreMarker) {
	if timeDirection == TimeDirectionBefore {
		param["beforePublishTime"] = moreMarker.TailMarker
	} else {
		if moreMarker.HeadMarker > moreMarker.TailMarker {
			param["afterPublishTime"] = moreMarker.HeadMarker
		} else {
			param["afterPublishTime"] = moreMarker.TailMarker
		}
	}
}

type moreList[T any] struct {
	HasMore  bool `json:"hasMore"`
	List     []T  `json:"list"`
	Counters any  `json:"counters"`
	moreMarker
}

type moreMarker struct {
	HeadMarker string `json:"headMarker"`
	TailMarker string `json:"tailMarker"`
}
