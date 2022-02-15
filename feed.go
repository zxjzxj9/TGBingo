package main

import (
	"github.com/mmcdole/gofeed"
	"strings"
)

func getFeed(url string) string {
	// url e.g. "http://www3.nhk.or.jp/rss/news/cat0.xml"
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(url)
	var sb strings.Builder
	for _, item := range feed.Items {
		sb.WriteString(item.Title + "\n")
		sb.WriteString(item.Links[0] + "\n")
	}
	return sb.String()
}
