package ofdl

import (
	"regexp"
	"testing"
)

func TestOFURL(t *testing.T) {
	homeURLs := []testURLSt{
		{url: "https://onlyfans.com/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com", must: []string{}, mustValue: []string{}},
	}

	subscriptionsURLs := []testURLSt{
		{url: "https://onlyfans.com/my/collections/user-lists/subscribers", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscribers/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscribers/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscribers?test=test", must: []string{}, mustValue: []string{}},

		{url: "https://onlyfans.com/my/collections/user-lists/subscribers/expired", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscribers/expired/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscribers/expired/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscribers/expired?test=test", must: []string{}, mustValue: []string{}},

		{url: "https://onlyfans.com/my/collections/user-lists/restricted", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/restricted/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/restricted/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/restricted?test=test", must: []string{}, mustValue: []string{}},

		{url: "https://onlyfans.com/my/collections/user-lists/blocked", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/blocked/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/blocked/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/blocked?test=test", must: []string{}, mustValue: []string{}},

		{url: "https://onlyfans.com/my/collections/user-lists/subscribers/active", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscribers/active/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscribers/active/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscribers/active?test=test", must: []string{}, mustValue: []string{}},

		{url: "https://onlyfans.com/my/collections/user-lists/subscriptions/active", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscriptions/active/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscriptions/active/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/subscriptions/active?test=test", must: []string{}, mustValue: []string{}},
	}

	chatURLs := []testURLSt{
		{url: "https://onlyfans.com/my/chats", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/chats/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/chats/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/chats?test=test", must: []string{}, mustValue: []string{}},

		{url: "https://onlyfans.com/my/chats/chat/342724494", must: []string{"ID"}, mustValue: []string{"342724494"}},
		{url: "https://onlyfans.com/my/chats/chat/342724494/", must: []string{"ID"}, mustValue: []string{"342724494"}},
		{url: "https://onlyfans.com/my/chats/chat/342724494/?test=test", must: []string{"ID"}, mustValue: []string{"342724494"}},
		{url: "https://onlyfans.com/my/chats/chat/342724494?test=test", must: []string{"ID"}, mustValue: []string{"342724494"}},
	}

	userListURLs := []testURLSt{
		{url: "https://onlyfans.com/my/collections/user-lists", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/user-lists?test=test", must: []string{}, mustValue: []string{}},

		{url: "https://onlyfans.com/my/collections/user-lists/1141544940", must: []string{"ID"}, mustValue: []string{"1141544940"}},
		{url: "https://onlyfans.com/my/collections/user-lists/1141544940/", must: []string{"ID"}, mustValue: []string{"1141544940"}},
		{url: "https://onlyfans.com/my/collections/user-lists/1141544940/?test=test", must: []string{"ID"}, mustValue: []string{"1141544940"}},
		{url: "https://onlyfans.com/my/collections/user-lists/1141544940?test=test", must: []string{"ID"}, mustValue: []string{"1141544940"}},
	}

	postURLs := []testURLSt{
		{url: "https://onlyfans.com/1353172156/onlyfans", must: []string{"PostID", "UserName"}, mustValue: []string{"1353172156", "onlyfans"}},
		{url: "https://onlyfans.com/1353172156/onlyfans/?test=test", must: []string{"PostID", "UserName"}, mustValue: []string{"1353172156", "onlyfans"}},
		{url: "https://onlyfans.com/1353172156/onlyfans?test=test", must: []string{"PostID", "UserName"}, mustValue: []string{"1353172156", "onlyfans"}},
	}

	userURLs := []testURLSt{
		{url: "https://onlyfans.com/kira.asia.ts", must: []string{"UserName"}, mustValue: []string{"kira.asia.ts"}},
		{url: "https://onlyfans.com/kira.asia.ts/", must: []string{"UserName"}, mustValue: []string{"kira.asia.ts"}},
		{url: "https://onlyfans.com/kira.asia.ts/?test=test", must: []string{"UserName"}, mustValue: []string{"kira.asia.ts"}},
		{url: "https://onlyfans.com/kira.asia.ts?test=test", must: []string{"UserName"}, mustValue: []string{"kira.asia.ts"}},

		{url: "https://onlyfans.com/olivoil2/media", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "media"}},
		{url: "https://onlyfans.com/olivoil2/media/", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "media"}},
		{url: "https://onlyfans.com/olivoil2/media/?test=test", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "media"}},
		{url: "https://onlyfans.com/olivoil2/media?test=test", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "media"}},
		{url: "https://onlyfans.com/olivoil2/videos", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "videos"}},
		{url: "https://onlyfans.com/olivoil2/videos/", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "videos"}},
		{url: "https://onlyfans.com/olivoil2/videos/?test=test", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "videos"}},
		{url: "https://onlyfans.com/olivoil2/videos?test=test", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "videos"}},
		{url: "https://onlyfans.com/olivoil2/photos", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "photos"}},
		{url: "https://onlyfans.com/olivoil2/photos/", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "photos"}},
		{url: "https://onlyfans.com/olivoil2/photos/?test=test", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "photos"}},
		{url: "https://onlyfans.com/olivoil2/photos?test=test", must: []string{"UserName", "MediaType"}, mustValue: []string{"olivoil2", "photos"}},
	}

	bookmarkURLs := []testURLSt{
		{url: "https://onlyfans.com/my/collections/bookmarks", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/bookmarks/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/bookmarks/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/bookmarks?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/?test=test", must: []string{}, mustValue: []string{}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all?test=test", must: []string{}, mustValue: []string{}},

		{url: "https://onlyfans.com/my/collections/bookmarks/all/photos", must: []string{"MediaType"}, mustValue: []string{"photos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/photos/", must: []string{"MediaType"}, mustValue: []string{"photos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/photos/?test=test", must: []string{"MediaType"}, mustValue: []string{"photos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/photos?test=test", must: []string{"MediaType"}, mustValue: []string{"photos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/videos", must: []string{"MediaType"}, mustValue: []string{"videos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/videos/", must: []string{"MediaType"}, mustValue: []string{"videos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/videos/?test=test", must: []string{"MediaType"}, mustValue: []string{"videos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/videos?test=test", must: []string{"MediaType"}, mustValue: []string{"videos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/audios", must: []string{"MediaType"}, mustValue: []string{"audios"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/audios/", must: []string{"MediaType"}, mustValue: []string{"audios"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/audios/?test=test", must: []string{"MediaType"}, mustValue: []string{"audios"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/audios?test=test", must: []string{"MediaType"}, mustValue: []string{"audios"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/other", must: []string{"MediaType"}, mustValue: []string{"other"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/other/", must: []string{"MediaType"}, mustValue: []string{"other"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/other/?test=test", must: []string{"MediaType"}, mustValue: []string{"other"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/other?test=test", must: []string{"MediaType"}, mustValue: []string{"other"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/locked", must: []string{"MediaType"}, mustValue: []string{"locked"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/locked/", must: []string{"MediaType"}, mustValue: []string{"locked"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/locked/?test=test", must: []string{"MediaType"}, mustValue: []string{"locked"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/all/locked?test=test", must: []string{"MediaType"}, mustValue: []string{"locked"}},

		{url: "https://onlyfans.com/my/collections/bookmarks/1979194", must: []string{"ID"}, mustValue: []string{"1979194"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/", must: []string{"ID"}, mustValue: []string{"1979194"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/?test=test", must: []string{"ID"}, mustValue: []string{"1979194"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194?test=test", must: []string{"ID"}, mustValue: []string{"1979194"}},

		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/photos", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "photos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/photos/", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "photos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/photos/?test=test", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "photos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/photos?test=test", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "photos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/videos", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "videos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/videos/", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "videos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/videos/?test=test", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "videos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/videos?test=test", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "videos"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/audios", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "audios"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/audios/", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "audios"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/audios/?test=test", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "audios"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/audios?test=test", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "audios"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/other", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "other"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/other/", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "other"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/other/?test=test", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "other"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/other?test=test", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "other"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/locked", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "locked"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/locked/", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "locked"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/locked/?test=test", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "locked"}},
		{url: "https://onlyfans.com/my/collections/bookmarks/1979194/locked?test=test", must: []string{"ID", "MediaType"}, mustValue: []string{"1979194", "locked"}},
	}

	testShouldMatchURLs(t, reHome, homeURLs)
	testShouldMatchURLs(t, reSubscriptions, subscriptionsURLs)
	testShouldMatchURLs(t, reChat, chatURLs)
	testShouldMatchURLs(t, reUserList, userListURLs)
	testShouldMatchURLs(t, reSinglePost, postURLs)
	testShouldMatchURLs(t, reUserWithMediaType, userURLs)
	testShouldMatchURLs(t, reBookmarksWithMediaType, bookmarkURLs)

	testShouldNotMatchURLs(t, []*regexp.Regexp{
		/*reHome,*/ reSubscriptions, reChat,
		reUserList, reSinglePost, reUserWithMediaType,
		reBookmarksWithMediaType}, homeURLs)

	testShouldNotMatchURLs(t, []*regexp.Regexp{
		reHome /*reSubscriptions */, reChat,
		reUserList, reSinglePost, reUserWithMediaType,
		reBookmarksWithMediaType}, subscriptionsURLs)

	testShouldNotMatchURLs(t, []*regexp.Regexp{
		reHome, reSubscriptions, /*reChat */
		reUserList, reSinglePost, reUserWithMediaType,
		reBookmarksWithMediaType}, chatURLs)

	testShouldNotMatchURLs(t, []*regexp.Regexp{
		reHome, reSubscriptions, reChat,
		/*reUserList,*/ reSinglePost, reUserWithMediaType,
		reBookmarksWithMediaType}, userListURLs)

	testShouldNotMatchURLs(t, []*regexp.Regexp{
		reHome, reSubscriptions, reChat,
		reUserList /*reSinglePost,*/, reUserWithMediaType,
		reBookmarksWithMediaType}, postURLs)

	testShouldNotMatchURLs(t, []*regexp.Regexp{
		reHome, reSubscriptions, reChat,
		reUserList, reSinglePost, /*reUserWithMediaType,*/
		reBookmarksWithMediaType}, userURLs)

	testShouldNotMatchURLs(t, []*regexp.Regexp{
		reHome, reSubscriptions, reChat,
		reUserList, reSinglePost, reUserWithMediaType,
		/*reBookmarksWithMediaType,*/}, bookmarkURLs)
}

type testURLSt struct {
	url           string
	must          []string
	mustValue     []string
	optional      []string
	optionalValue []string
}

func testShouldMatchURLs(t *testing.T, re *regexp.Regexp, tests []testURLSt) {
	for _, ts := range tests {
		if len(ts.must) == 0 && len(ts.optional) == 0 {
			if !ofurlMatchs(ts.url, re) {
				t.Logf("url should match: %s, re: %s", ts.url, re.String())
				t.Fail()
			}
			continue
		}

		if (len(ts.must) != 0 || len(ts.mustValue) != 0) && len(ts.mustValue) != len(ts.must) {
			t.Logf("mustValue length must be equal to must length")
			t.Fail()
		}

		if (len(ts.optional) != 0 || len(ts.optionalValue) != 0) && len(ts.optionalValue) != len(ts.optional) {
			t.Logf("optionalValue length must be equal to optional length")
			t.Fail()
		}

		if founds, ok := ofurlFinds(ts.must, ts.optional, ts.url, re); !ok {
			t.Logf("url should match: %s, re: %s", ts.url, re.String())
			t.Fail()
		} else {
			for i := range ts.must {
				if founds[i] != ts.mustValue[i] {
					t.Logf("url should match: %s, re: %s", ts.url, re.String())
					t.Fail()
				}
			}
			for i := range ts.optional {
				if founds[len(ts.must)+i] != ts.optionalValue[i] {
					t.Logf("url should match: %s, re: %s", ts.url, re.String())
					t.Fail()
				}
			}
		}
	}
}

func testShouldNotMatchURLs(t *testing.T, res []*regexp.Regexp, tests []testURLSt) {
	for _, ts := range tests {
		for _, re := range res {
			if ok := ofurlMatchs(ts.url, re); ok {
				t.Logf("url should not match: %s, re: %s", ts.url, re.String())
				t.Fail()
			}
		}
	}
}
