package ofdl

import (
	"regexp"
	"testing"
)

func TestOFURL(t *testing.T) {
	homeURLs := []string{
		"https://onlyfans.com/?test=test",
		"https://onlyfans.com?test=test",
		"https://onlyfans.com/",
		"https://onlyfans.com",
	}

	subscriptionsURLs := []string{
		"https://onlyfans.com/my/collections/user-lists/subscribers/active",
		"https://onlyfans.com/my/collections/user-lists/subscribers/active/",
		"https://onlyfans.com/my/collections/user-lists/subscribers/active/?test=test",
		"https://onlyfans.com/my/collections/user-lists/subscribers/active?test=test",
		"https://onlyfans.com/my/collections/user-lists/subscriptions/active",
		"https://onlyfans.com/my/collections/user-lists/subscriptions/active/",
		"https://onlyfans.com/my/collections/user-lists/subscriptions/active/?test=test",
		"https://onlyfans.com/my/collections/user-lists/subscriptions/active?test=test",
	}

	chartsURLs := []string{
		"https://onlyfans.com/my/chats",
		"https://onlyfans.com/my/chats/",
		"https://onlyfans.com/my/chats/?test=test",
		"https://onlyfans.com/my/chats?test=test",
	}

	chatURLs := []testFindKeyURLSt{
		{url: "https://onlyfans.com/my/chats/chat/342724494", key: "ID", value: "342724494"},
		{url: "https://onlyfans.com/my/chats/chat/342724494/", key: "ID", value: "342724494"},
		{url: "https://onlyfans.com/my/chats/chat/342724494/?test=test", key: "ID", value: "342724494"},
		{url: "https://onlyfans.com/my/chats/chat/342724494?test=test", key: "ID", value: "342724494"},
	}

	userListURLs := []testFindKeyURLSt{
		{url: "https://onlyfans.com/my/collections/user-lists/1141544940", key: "ID", value: "1141544940"},
		{url: "https://onlyfans.com/my/collections/user-lists/1141544940/", key: "ID", value: "1141544940"},
		{url: "https://onlyfans.com/my/collections/user-lists/1141544940/?test=test", key: "ID", value: "1141544940"},
		{url: "https://onlyfans.com/my/collections/user-lists/1141544940?test=test", key: "ID", value: "1141544940"},
	}

	postURLs := []testFindKeyURLSt2{
		{url: "https://onlyfans.com/1353172156/onlyfans", key1: "PostID", value1: "1353172156", key2: "UserName", value2: "onlyfans"},
		{url: "https://onlyfans.com/1353172156/onlyfans/?test=test", key1: "PostID", value1: "1353172156", key2: "UserName", value2: "onlyfans"},
		{url: "https://onlyfans.com/1353172156/onlyfans?test=test", key1: "PostID", value1: "1353172156", key2: "UserName", value2: "onlyfans"},
	}

	userURLs := []testFindKeyURLSt{
		{url: "https://onlyfans.com/kira.asia.ts", key: "UserName", value: "kira.asia.ts"},
		{url: "https://onlyfans.com/kira.asia.ts/", key: "UserName", value: "kira.asia.ts"},
		{url: "https://onlyfans.com/kira.asia.ts/?test=test", key: "UserName", value: "kira.asia.ts"},
		{url: "https://onlyfans.com/kira.asia.ts?test=test", key: "UserName", value: "kira.asia.ts"},
	}

	for _, url := range homeURLs {
		if !isOFHomeURL(url) {
			t.Fail()
		}
	}

	testMatchURLs(t, reSubscriptions, subscriptionsURLs)
	testMatchURLs(t, reChats, chartsURLs)
	testFindKeyURLs(t, reSingleChat, chatURLs)
	testFindKeyURLs(t, reUserList, userListURLs)
	testFindKey2URLs(t, reSinglePost, postURLs)
	testFindKeyURLs(t, reUser, userURLs)

	testNotMatchURLs(t, reSubscriptions, homeURLs)
	testNotMatchURLs(t, reChats, homeURLs)
	testNotMatchURLs(t, reSingleChat, homeURLs)
	testNotMatchURLs(t, reUserList, homeURLs)
	testNotMatchURLs(t, reSinglePost, homeURLs)
	testNotMatchURLs(t, reUser, homeURLs)

	testNotMatchURLs(t, reChats, subscriptionsURLs)
	testNotMatchURLs(t, reSingleChat, subscriptionsURLs)
	testNotMatchURLs(t, reUserList, subscriptionsURLs)
	testNotMatchURLs(t, reSinglePost, subscriptionsURLs)
	testNotMatchURLs(t, reUser, subscriptionsURLs)

	testNotMatchURLs(t, reSubscriptions, chartsURLs)
	testNotMatchURLs(t, reSingleChat, chartsURLs)
	testNotMatchURLs(t, reUserList, chartsURLs)
	testNotMatchURLs(t, reSinglePost, chartsURLs)
	testNotMatchURLs(t, reUser, chartsURLs)

	testNotMatchKeyURLs(t, reSubscriptions, chatURLs)
	testNotMatchKeyURLs(t, reChats, chatURLs)
	testNotMatchKeyURLs(t, reUserList, chatURLs)
	testNotMatchKeyURLs(t, reSinglePost, chatURLs)
	testNotMatchKeyURLs(t, reUser, chatURLs)

	testNotMatchKeyURLs(t, reSubscriptions, userListURLs)
	testNotMatchKeyURLs(t, reChats, userListURLs)
	testNotMatchKeyURLs(t, reSingleChat, userListURLs)
	testNotMatchKeyURLs(t, reSinglePost, userListURLs)
	testNotMatchKeyURLs(t, reUser, userListURLs)

	testNotMatchKey2URLs(t, reSubscriptions, postURLs)
	testNotMatchKey2URLs(t, reChats, postURLs)
	testNotMatchKey2URLs(t, reSingleChat, postURLs)
	testNotMatchKey2URLs(t, reUserList, postURLs)
	testNotMatchKey2URLs(t, reUser, postURLs)

	testNotMatchKeyURLs(t, reSubscriptions, userURLs)
	testNotMatchKeyURLs(t, reChats, userURLs)
	testNotMatchKeyURLs(t, reSingleChat, userURLs)
	testNotMatchKeyURLs(t, reUserList, userURLs)
	testNotMatchKeyURLs(t, reSinglePost, userURLs)
}

func testMatchURLs(t *testing.T, re *regexp.Regexp, urls []string) {
	for _, url := range urls {
		if !ofurlMatch(re, url) {
			t.Logf("url should match: %s, re: %s", url, re.String())
			t.Fail()
		}
	}
}

func testNotMatchURLs(t *testing.T, re *regexp.Regexp, urls []string) {
	for _, url := range urls {
		if ofurlMatch(re, url) {
			t.Logf("url should not match: %s, re: %s", url, re.String())
			t.Fail()
		}
	}
}

type testFindKeyURLSt struct {
	url   string
	key   string
	value string
}

func testFindKeyURLs(t *testing.T, re *regexp.Regexp, tests []testFindKeyURLSt) {
	for _, ts := range tests {
		if value, ok := ofurlFind(re, ts.url, ts.key); !ok || value != ts.value {
			t.Logf("url should match: %s, re: %s", ts.url, re.String())
			t.Fail()
		}
	}
}

func testNotMatchKeyURLs(t *testing.T, re *regexp.Regexp, tests []testFindKeyURLSt) {
	for _, ts := range tests {
		if _, ok := ofurlFind(re, ts.url, ts.key); ok {
			t.Logf("url should not match: %s, re: %s", ts.url, re.String())
			t.Fail()
		}
	}
}

type testFindKeyURLSt2 struct {
	url    string
	key1   string
	value1 string
	key2   string
	value2 string
}

func testFindKey2URLs(t *testing.T, re *regexp.Regexp, tests []testFindKeyURLSt2) {
	for _, ts := range tests {
		if value1, value2, ok := ofurlFind2(re, ts.url, ts.key1, ts.key2); !ok || value1 != ts.value1 || value2 != ts.value2 {
			t.Logf("url should not match: %s, re: %s", ts.url, re.String())
			t.Fail()
		}
	}
}

func testNotMatchKey2URLs(t *testing.T, re *regexp.Regexp, tests []testFindKeyURLSt2) {
	for _, ts := range tests {
		if _, _, ok := ofurlFind2(re, ts.url, ts.key1, ts.key2); ok {
			t.Logf("url should not match: %s, re: %s", ts.url, re.String())
			t.Fail()
		}
	}
}
