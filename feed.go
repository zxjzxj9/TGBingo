package main

import (
	"github.com/mmcdole/gofeed"
)

func getFeed(url string) []string {
	// url e.g. "http://www3.nhk.or.jp/rss/news/cat0.xml"
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	ret := make([]string, 0)
	for _, item := range feed.Items {
		ret = append(ret, item.Title+"\n"+item.Links[0])
	}
	return ret
}
