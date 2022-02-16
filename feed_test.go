package main

import "testing"

func TestGetFeed(t *testing.T) {
	getFeed("http://www3.nhk.or.jp/rss/news/cat0.xml")
}
