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

	userMediaTypeURLs := []testFindKeyURLSt2{
		{url: "https://onlyfans.com/olivoil2/media", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "media"},
		{url: "https://onlyfans.com/olivoil2/media/", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "media"},
		{url: "https://onlyfans.com/olivoil2/media/?test=test", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "media"},
		{url: "https://onlyfans.com/olivoil2/media?test=test", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "media"},
		{url: "https://onlyfans.com/olivoil2/videos", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "videos"},
		{url: "https://onlyfans.com/olivoil2/videos/", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "videos"},
		{url: "https://onlyfans.com/olivoil2/videos/?test=test", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "videos"},
		{url: "https://onlyfans.com/olivoil2/videos?test=test", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "videos"},
		{url: "https://onlyfans.com/olivoil2/photos", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "photos"},
		{url: "https://onlyfans.com/olivoil2/photos/", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "photos"},
		{url: "https://onlyfans.com/olivoil2/photos/?test=test", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "photos"},
		{url: "https://onlyfans.com/olivoil2/photos?test=test", key1: "UserName", value1: "olivoil2", key2: "MediaType", value2: "photos"},
	}

	allBookmarkURLs := []string{
		"https://onlyfans.com/my/collections/bookmarks",
		"https://onlyfans.com/my/collections/bookmarks/",
		"https://onlyfans.com/my/collections/bookmarks/?test=test",
		"https://onlyfans.com/my/collections/bookmarks?test=test",
		"https://onlyfans.com/my/collections/bookmarks/all",
		"https://onlyfans.com/my/collections/bookmarks/all/",
		"https://onlyfans.com/my/collections/bookmarks/all/?test=test",
		"https://onlyfans.com/my/collections/bookmarks/all?test=test",
	}

	allBookmarkByMediaTypeURLs := []testFindKeyURLSt{
		{url: "https://onlyfans.com/my/collections/bookmarks/all/photos", key: "MediaType", value: "photos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/photos/", key: "MediaType", value: "photos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/photos/?test=test", key: "MediaType", value: "photos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/photos?test=test", key: "MediaType", value: "photos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/videos", key: "MediaType", value: "videos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/videos/", key: "MediaType", value: "videos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/videos/?test=test", key: "MediaType", value: "videos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/videos?test=test", key: "MediaType", value: "videos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/audios", key: "MediaType", value: "audios"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/audios/", key: "MediaType", value: "audios"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/audios/?test=test", key: "MediaType", value: "audios"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/audios?test=test", key: "MediaType", value: "audios"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/other", key: "MediaType", value: "other"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/other/", key: "MediaType", value: "other"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/other/?test=test", key: "MediaType", value: "other"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/other?test=test", key: "MediaType", value: "other"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/locked", key: "MediaType", value: "locked"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/locked/", key: "MediaType", value: "locked"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/locked/?test=test", key: "MediaType", value: "locked"},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/locked?test=test", key: "MediaType", value: "locked"},
	}

	singleBookmarkURLs := []testFindKeyURLSt{
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194", key: "ID", value: "1979194"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/", key: "ID", value: "1979194"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/?test=test", key: "ID", value: "1979194"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194?test=test", key: "ID", value: "1979194"},
	}

	singleBookmarkByMediaTypeURLs := []testFindKeyURLSt2{
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/photos", key1: "ID", value1: "1979194", key2: "MediaType", value2: "photos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/photos/", key1: "ID", value1: "1979194", key2: "MediaType", value2: "photos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/photos/?test=test", key1: "ID", value1: "1979194", key2: "MediaType", value2: "photos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/photos?test=test", key1: "ID", value1: "1979194", key2: "MediaType", value2: "photos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/videos", key1: "ID", value1: "1979194", key2: "MediaType", value2: "videos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/videos/", key1: "ID", value1: "1979194", key2: "MediaType", value2: "videos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/videos/?test=test", key1: "ID", value1: "1979194", key2: "MediaType", value2: "videos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/videos?test=test", key1: "ID", value1: "1979194", key2: "MediaType", value2: "videos"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/audios", key1: "ID", value1: "1979194", key2: "MediaType", value2: "audios"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/audios/", key1: "ID", value1: "1979194", key2: "MediaType", value2: "audios"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/audios/?test=test", key1: "ID", value1: "1979194", key2: "MediaType", value2: "audios"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/audios?test=test", key1: "ID", value1: "1979194", key2: "MediaType", value2: "audios"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/other", key1: "ID", value1: "1979194", key2: "MediaType", value2: "other"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/other/", key1: "ID", value1: "1979194", key2: "MediaType", value2: "other"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/other/?test=test", key1: "ID", value1: "1979194", key2: "MediaType", value2: "other"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/other?test=test", key1: "ID", value1: "1979194", key2: "MediaType", value2: "other"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/locked", key1: "ID", value1: "1979194", key2: "MediaType", value2: "locked"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/locked/", key1: "ID", value1: "1979194", key2: "MediaType", value2: "locked"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/locked/?test=test", key1: "ID", value1: "1979194", key2: "MediaType", value2: "locked"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/locked?test=test", key1: "ID", value1: "1979194", key2: "MediaType", value2: "locked"},
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
	testFindKey2URLs(t, reUserByMediaType, userMediaTypeURLs)
	testMatchURLs(t, reAllBookmarks, allBookmarkURLs)
	testFindKeyURLs(t, reAllBookmarksByMediaType, allBookmarkByMediaTypeURLs)
	testFindKeyURLs(t, reSingleBookmark, singleBookmarkURLs)
	testFindKey2URLs(t, reSingleBookmarkByMediaType, singleBookmarkByMediaTypeURLs)

	testNotMatchURLs(t, reSubscriptions, homeURLs)
	testNotMatchURLs(t, reChats, homeURLs)
	testNotMatchURLs(t, reSingleChat, homeURLs)
	testNotMatchURLs(t, reUserList, homeURLs)
	testNotMatchURLs(t, reSinglePost, homeURLs)
	testNotMatchURLs(t, reUser, homeURLs)
	testNotMatchURLs(t, reUserByMediaType, homeURLs)
	testNotMatchURLs(t, reAllBookmarks, homeURLs)
	testNotMatchURLs(t, reSingleBookmark, homeURLs)
	testNotMatchURLs(t, reAllBookmarksByMediaType, homeURLs)
	testNotMatchURLs(t, reSingleBookmarkByMediaType, homeURLs)

	testNotMatchURLs(t, reChats, subscriptionsURLs)
	testNotMatchURLs(t, reSingleChat, subscriptionsURLs)
	testNotMatchURLs(t, reUserList, subscriptionsURLs)
	testNotMatchURLs(t, reSinglePost, subscriptionsURLs)
	testNotMatchURLs(t, reUser, subscriptionsURLs)
	testNotMatchURLs(t, reUserByMediaType, subscriptionsURLs)
	testNotMatchURLs(t, reAllBookmarks, subscriptionsURLs)
	testNotMatchURLs(t, reSingleBookmark, subscriptionsURLs)
	testNotMatchURLs(t, reAllBookmarksByMediaType, subscriptionsURLs)
	testNotMatchURLs(t, reSingleBookmarkByMediaType, subscriptionsURLs)

	testNotMatchURLs(t, reSubscriptions, chartsURLs)
	testNotMatchURLs(t, reSingleChat, chartsURLs)
	testNotMatchURLs(t, reUserList, chartsURLs)
	testNotMatchURLs(t, reSinglePost, chartsURLs)
	testNotMatchURLs(t, reUser, chartsURLs)
	testNotMatchURLs(t, reUserByMediaType, chartsURLs)
	testNotMatchURLs(t, reAllBookmarks, chartsURLs)
	testNotMatchURLs(t, reSingleBookmark, chartsURLs)
	testNotMatchURLs(t, reAllBookmarksByMediaType, chartsURLs)
	testNotMatchURLs(t, reSingleBookmarkByMediaType, chartsURLs)

	testNotMatchKeyURLs(t, reSubscriptions, chatURLs)
	testNotMatchKeyURLs(t, reChats, chatURLs)
	testNotMatchKeyURLs(t, reUserList, chatURLs)
	testNotMatchKeyURLs(t, reSinglePost, chatURLs)
	testNotMatchKeyURLs(t, reUser, chatURLs)
	testNotMatchKeyURLs(t, reUserByMediaType, chatURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, chatURLs)
	testNotMatchKeyURLs(t, reSingleBookmark, chatURLs)
	testNotMatchKeyURLs(t, reAllBookmarksByMediaType, chatURLs)
	testNotMatchKeyURLs(t, reSingleBookmarkByMediaType, chatURLs)

	testNotMatchKeyURLs(t, reSubscriptions, userListURLs)
	testNotMatchKeyURLs(t, reChats, userListURLs)
	testNotMatchKeyURLs(t, reSingleChat, userListURLs)
	testNotMatchKeyURLs(t, reSinglePost, userListURLs)
	testNotMatchKeyURLs(t, reUser, userListURLs)
	testNotMatchKeyURLs(t, reUserByMediaType, userListURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, userListURLs)
	testNotMatchKeyURLs(t, reSingleBookmark, userListURLs)
	testNotMatchKeyURLs(t, reAllBookmarksByMediaType, userListURLs)
	testNotMatchKeyURLs(t, reSingleBookmarkByMediaType, userListURLs)

	testNotMatchKey2URLs(t, reSubscriptions, postURLs)
	testNotMatchKey2URLs(t, reChats, postURLs)
	testNotMatchKey2URLs(t, reSingleChat, postURLs)
	testNotMatchKey2URLs(t, reUserList, postURLs)
	testNotMatchKey2URLs(t, reUser, postURLs)
	testNotMatchKey2URLs(t, reUserByMediaType, postURLs)
	testNotMatchKey2URLs(t, reAllBookmarks, postURLs)
	testNotMatchKey2URLs(t, reSingleBookmark, postURLs)
	testNotMatchKey2URLs(t, reAllBookmarksByMediaType, postURLs)
	testNotMatchKey2URLs(t, reSingleBookmarkByMediaType, postURLs)

	testNotMatchKeyURLs(t, reSubscriptions, userURLs)
	testNotMatchKeyURLs(t, reChats, userURLs)
	testNotMatchKeyURLs(t, reSingleChat, userURLs)
	testNotMatchKeyURLs(t, reUserList, userURLs)
	testNotMatchKeyURLs(t, reSinglePost, userURLs)
	testNotMatchKeyURLs(t, reUserByMediaType, userURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, userURLs)
	testNotMatchKeyURLs(t, reSingleBookmark, userURLs)
	testNotMatchKeyURLs(t, reAllBookmarksByMediaType, userURLs)
	testNotMatchKeyURLs(t, reSingleBookmarkByMediaType, userURLs)

	testNotMatchKey2URLs(t, reSubscriptions, userMediaTypeURLs)
	testNotMatchKey2URLs(t, reChats, userMediaTypeURLs)
	testNotMatchKey2URLs(t, reSingleChat, userMediaTypeURLs)
	testNotMatchKey2URLs(t, reUserList, userMediaTypeURLs)
	testNotMatchKey2URLs(t, reSinglePost, userMediaTypeURLs)
	testNotMatchKey2URLs(t, reUser, userMediaTypeURLs)
	testNotMatchKey2URLs(t, reAllBookmarks, userMediaTypeURLs)
	testNotMatchKey2URLs(t, reSingleBookmark, userMediaTypeURLs)
	testNotMatchKey2URLs(t, reAllBookmarksByMediaType, userMediaTypeURLs)
	testNotMatchKey2URLs(t, reSingleBookmarkByMediaType, userMediaTypeURLs)

	testNotMatchURLs(t, reSubscriptions, allBookmarkURLs)
	testNotMatchURLs(t, reChats, allBookmarkURLs)
	testNotMatchURLs(t, reSingleChat, allBookmarkURLs)
	testNotMatchURLs(t, reUserList, allBookmarkURLs)
	testNotMatchURLs(t, reSinglePost, allBookmarkURLs)
	testNotMatchURLs(t, reUser, allBookmarkURLs)
	testNotMatchURLs(t, reUserByMediaType, allBookmarkURLs)
	testNotMatchURLs(t, reSingleBookmark, allBookmarkURLs)
	testNotMatchURLs(t, reAllBookmarksByMediaType, allBookmarkURLs)
	testNotMatchURLs(t, reSingleBookmarkByMediaType, allBookmarkURLs)

	testNotMatchKeyURLs(t, reSubscriptions, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reChats, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reSingleChat, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reUserList, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reSinglePost, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reUser, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reUserByMediaType, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reAllBookmarksByMediaType, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reSingleBookmarkByMediaType, singleBookmarkURLs)
}

func testMatchURLs(t *testing.T, re *regexp.Regexp, urls []string) {
	for _, url := range urls {
		if !ofurlMatchs(url, re) {
			t.Logf("url should match: %s, re: %s", url, re.String())
			t.Fail()
		}
	}
}

func testNotMatchURLs(t *testing.T, re *regexp.Regexp, urls []string) {
	for _, url := range urls {
		if ofurlMatchs(url, re) {
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
		if value, ok := ofurlFinds(ts.url, ts.key, re); !ok || value != ts.value {
			t.Logf("url should match: %s, re: %s", ts.url, re.String())
			t.Fail()
		}
	}
}

func testNotMatchKeyURLs(t *testing.T, re *regexp.Regexp, tests []testFindKeyURLSt) {
	for _, ts := range tests {
		if _, ok := ofurlFinds(ts.url, ts.key, re); ok {
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
		if value1, value2, ok := ofurlFinds2(ts.url, ts.key1, ts.key2, re); !ok || value1 != ts.value1 || value2 != ts.value2 {
			t.Logf("url should not match: %s, re: %s", ts.url, re.String())
			t.Fail()
		}
	}
}

func testNotMatchKey2URLs(t *testing.T, re *regexp.Regexp, tests []testFindKeyURLSt2) {
	for _, ts := range tests {
		if _, _, ok := ofurlFinds2(ts.url, ts.key1, ts.key2, re); ok {
			t.Logf("url should not match: %s, re: %s", ts.url, re.String())
			t.Fail()
		}
	}
}
