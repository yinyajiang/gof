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

	userMediaURLs := []testFindKeyURLSt{
		{url: "https://onlyfans.com/olivoil2/media", key: "UserName", value: "olivoil2"},
		{url: "https://onlyfans.com/olivoil2/media/", key: "UserName", value: "olivoil2"},
		{url: "https://onlyfans.com/olivoil2/media/?test=test", key: "UserName", value: "olivoil2"},
		{url: "https://onlyfans.com/olivoil2/media?test=test", key: "UserName", value: "olivoil2"},
	}

	userVideosURLs := []testFindKeyURLSt{
		{url: "https://onlyfans.com/olivoil2/videos", key: "UserName", value: "olivoil2"},
		{url: "https://onlyfans.com/olivoil2/videos/", key: "UserName", value: "olivoil2"},
		{url: "https://onlyfans.com/olivoil2/videos/?test=test", key: "UserName", value: "olivoil2"},
		{url: "https://onlyfans.com/olivoil2/videos?test=test", key: "UserName", value: "olivoil2"},
	}

	userPhotosURLs := []testFindKeyURLSt{
		{url: "https://onlyfans.com/olivoil2/photos", key: "UserName", value: "olivoil2"},
		{url: "https://onlyfans.com/olivoil2/photos/", key: "UserName", value: "olivoil2"},
		{url: "https://onlyfans.com/olivoil2/photos/?test=test", key: "UserName", value: "olivoil2"},
		{url: "https://onlyfans.com/olivoil2/photos?test=test", key: "UserName", value: "olivoil2"},
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

	singleBookmarkURLs := []testFindKeyURLSt{
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194", key: "ID", value: "1979194"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/", key: "ID", value: "1979194"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/?test=test", key: "ID", value: "1979194"},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194?test=test", key: "ID", value: "1979194"},
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
	testFindKeyURLs(t, reUserMedia, userMediaURLs)
	testFindKeyURLs(t, reUserVideos, userVideosURLs)
	testFindKeyURLs(t, reUserPhotos, userPhotosURLs)
	testMatchURLs(t, reAllBookmarks, allBookmarkURLs)
	testFindKeyURLs(t, reSingleBookmark, singleBookmarkURLs)

	testNotMatchURLs(t, reSubscriptions, homeURLs)
	testNotMatchURLs(t, reChats, homeURLs)
	testNotMatchURLs(t, reSingleChat, homeURLs)
	testNotMatchURLs(t, reUserList, homeURLs)
	testNotMatchURLs(t, reSinglePost, homeURLs)
	testNotMatchURLs(t, reUser, homeURLs)
	testNotMatchURLs(t, reUserMedia, homeURLs)
	testNotMatchURLs(t, reAllBookmarks, homeURLs)
	testNotMatchURLs(t, reSingleBookmark, homeURLs)

	testNotMatchURLs(t, reChats, subscriptionsURLs)
	testNotMatchURLs(t, reSingleChat, subscriptionsURLs)
	testNotMatchURLs(t, reUserList, subscriptionsURLs)
	testNotMatchURLs(t, reSinglePost, subscriptionsURLs)
	testNotMatchURLs(t, reUser, subscriptionsURLs)
	testNotMatchURLs(t, reUserMedia, subscriptionsURLs)
	testNotMatchURLs(t, reUserVideos, subscriptionsURLs)
	testNotMatchURLs(t, reUserPhotos, subscriptionsURLs)
	testNotMatchURLs(t, reAllBookmarks, subscriptionsURLs)
	testNotMatchURLs(t, reSingleBookmark, subscriptionsURLs)

	testNotMatchURLs(t, reSubscriptions, chartsURLs)
	testNotMatchURLs(t, reSingleChat, chartsURLs)
	testNotMatchURLs(t, reUserList, chartsURLs)
	testNotMatchURLs(t, reSinglePost, chartsURLs)
	testNotMatchURLs(t, reUser, chartsURLs)
	testNotMatchURLs(t, reUserMedia, chartsURLs)
	testNotMatchURLs(t, reUserVideos, chartsURLs)
	testNotMatchURLs(t, reUserPhotos, chartsURLs)
	testNotMatchURLs(t, reAllBookmarks, chartsURLs)
	testNotMatchURLs(t, reSingleBookmark, chartsURLs)

	testNotMatchKeyURLs(t, reSubscriptions, chatURLs)
	testNotMatchKeyURLs(t, reChats, chatURLs)
	testNotMatchKeyURLs(t, reUserList, chatURLs)
	testNotMatchKeyURLs(t, reSinglePost, chatURLs)
	testNotMatchKeyURLs(t, reUser, chatURLs)
	testNotMatchKeyURLs(t, reUserMedia, chatURLs)
	testNotMatchKeyURLs(t, reUserVideos, chatURLs)
	testNotMatchKeyURLs(t, reUserPhotos, chatURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, chatURLs)
	testNotMatchKeyURLs(t, reSingleBookmark, chatURLs)

	testNotMatchKeyURLs(t, reSubscriptions, userListURLs)
	testNotMatchKeyURLs(t, reChats, userListURLs)
	testNotMatchKeyURLs(t, reSingleChat, userListURLs)
	testNotMatchKeyURLs(t, reSinglePost, userListURLs)
	testNotMatchKeyURLs(t, reUser, userListURLs)
	testNotMatchKeyURLs(t, reUserMedia, userListURLs)
	testNotMatchKeyURLs(t, reUserPhotos, userListURLs)
	testNotMatchKeyURLs(t, reUserVideos, userListURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, userListURLs)
	testNotMatchKeyURLs(t, reSingleBookmark, userListURLs)

	testNotMatchKey2URLs(t, reSubscriptions, postURLs)
	testNotMatchKey2URLs(t, reChats, postURLs)
	testNotMatchKey2URLs(t, reSingleChat, postURLs)
	testNotMatchKey2URLs(t, reUserList, postURLs)
	testNotMatchKey2URLs(t, reUser, postURLs)
	testNotMatchKey2URLs(t, reUserMedia, postURLs)
	testNotMatchKey2URLs(t, reUserPhotos, postURLs)
	testNotMatchKey2URLs(t, reUserVideos, postURLs)
	testNotMatchKey2URLs(t, reAllBookmarks, postURLs)
	testNotMatchKey2URLs(t, reSingleBookmark, postURLs)

	testNotMatchKeyURLs(t, reSubscriptions, userURLs)
	testNotMatchKeyURLs(t, reChats, userURLs)
	testNotMatchKeyURLs(t, reSingleChat, userURLs)
	testNotMatchKeyURLs(t, reUserList, userURLs)
	testNotMatchKeyURLs(t, reSinglePost, userURLs)
	testNotMatchKeyURLs(t, reUserMedia, userURLs)
	testNotMatchKeyURLs(t, reUserPhotos, userURLs)
	testNotMatchKeyURLs(t, reUserVideos, userURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, userURLs)
	testNotMatchKeyURLs(t, reSingleBookmark, userURLs)

	testNotMatchKeyURLs(t, reSubscriptions, userMediaURLs)
	testNotMatchKeyURLs(t, reChats, userMediaURLs)
	testNotMatchKeyURLs(t, reSingleChat, userMediaURLs)
	testNotMatchKeyURLs(t, reUserList, userMediaURLs)
	testNotMatchKeyURLs(t, reSinglePost, userMediaURLs)
	testNotMatchKeyURLs(t, reUser, userMediaURLs)
	testNotMatchKeyURLs(t, reUserPhotos, userMediaURLs)
	testNotMatchKeyURLs(t, reUserVideos, userMediaURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, userMediaURLs)
	testNotMatchKeyURLs(t, reSingleBookmark, userMediaURLs)

	testNotMatchKeyURLs(t, reSubscriptions, userVideosURLs)
	testNotMatchKeyURLs(t, reChats, userVideosURLs)
	testNotMatchKeyURLs(t, reSingleChat, userVideosURLs)
	testNotMatchKeyURLs(t, reUserList, userVideosURLs)
	testNotMatchKeyURLs(t, reSinglePost, userVideosURLs)
	testNotMatchKeyURLs(t, reUser, userVideosURLs)
	testNotMatchKeyURLs(t, reUserMedia, userVideosURLs)
	testNotMatchKeyURLs(t, reUserPhotos, userVideosURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, userVideosURLs)
	testNotMatchKeyURLs(t, reSingleBookmark, userVideosURLs)

	testNotMatchKeyURLs(t, reSubscriptions, userPhotosURLs)
	testNotMatchKeyURLs(t, reChats, userPhotosURLs)
	testNotMatchKeyURLs(t, reSingleChat, userPhotosURLs)
	testNotMatchKeyURLs(t, reUserList, userPhotosURLs)
	testNotMatchKeyURLs(t, reSinglePost, userPhotosURLs)
	testNotMatchKeyURLs(t, reUser, userPhotosURLs)
	testNotMatchKeyURLs(t, reUserMedia, userPhotosURLs)
	testNotMatchKeyURLs(t, reUserVideos, userPhotosURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, userPhotosURLs)
	testNotMatchKeyURLs(t, reSingleBookmark, userPhotosURLs)

	testNotMatchURLs(t, reSubscriptions, allBookmarkURLs)
	testNotMatchURLs(t, reChats, allBookmarkURLs)
	testNotMatchURLs(t, reSingleChat, allBookmarkURLs)
	testNotMatchURLs(t, reUserList, allBookmarkURLs)
	testNotMatchURLs(t, reSinglePost, allBookmarkURLs)
	testNotMatchURLs(t, reUser, allBookmarkURLs)
	testNotMatchURLs(t, reUserMedia, allBookmarkURLs)
	testNotMatchURLs(t, reUserVideos, allBookmarkURLs)
	testNotMatchURLs(t, reUserPhotos, allBookmarkURLs)
	testNotMatchURLs(t, reSingleBookmark, allBookmarkURLs)

	testNotMatchKeyURLs(t, reSubscriptions, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reChats, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reSingleChat, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reUserList, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reSinglePost, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reUser, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reUserMedia, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reUserVideos, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reUserPhotos, singleBookmarkURLs)
	testNotMatchKeyURLs(t, reAllBookmarks, singleBookmarkURLs)
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
