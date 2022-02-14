package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
)

func getFeed(url string) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("http://www3.nhk.or.jp/rss/news/cat0.xml")
	fmt.Println(feed.Title)
	fmt.Println(feed.Links)
	for _, item := range feed.Items {
		fmt.Println(item.Title)
	}
}
