package model

import "time"

type purchased struct {
	ResponseType        string     `json:"responseType"`
	Text                string     `json:"text"`
	GiphyID             any        `json:"giphyId"`
	LockedText          any        `json:"lockedText"`
	IsFree              bool       `json:"isFree"`
	Price               any        `json:"price"`
	IsMediaReady        bool       `json:"isMediaReady"`
	MediaCount          int        `json:"mediaCount"`
	Media               []Media    `json:"media"`
	Previews            []any      `json:"previews"`
	Preview             []any      `json:"preview"`
	IsTip               bool       `json:"isTip"`
	IsReportedByMe      bool       `json:"isReportedByMe"`
	IsCouplePeopleMedia bool       `json:"isCouplePeopleMedia"`
	QueueID             any        `json:"queueId"`
	FromUser            UserIDView `json:"fromUser"`
	Author              UserIDView `json:"author"`
	IsFromQueue         bool       `json:"isFromQueue"`
	CanUnsendQueue      bool       `json:"canUnsendQueue"`
	UnsendSecondsQueue  any        `json:"unsendSecondsQueue"`
	ID                  int64      `json:"id"`
	IsOpened            bool       `json:"isOpened"`
	IsNew               bool       `json:"isNew"`
	CreatedAt           time.Time  `json:"createdAt"`
	PostedAt            time.Time  `json:"postedAt"`
	ChangedAt           time.Time  `json:"changedAt"`
	CancelSeconds       int        `json:"cancelSeconds"`
	IsLiked             bool       `json:"isLiked"`
	CanPurchase         bool       `json:"canPurchase"`
	CanReport           bool       `json:"canReport"`
	IsCanceled          bool       `json:"isCanceled"`
	IsArchived          bool       `json:"isArchived"`
}
