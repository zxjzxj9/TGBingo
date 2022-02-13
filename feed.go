package main

import (
	"fmt"
	"github.com/mmcdole/gofeed"
)

func getFeed(url string) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL("http://feeds.twit.tv/twit.xml")
	fmt.Println(feed.Title)
}
