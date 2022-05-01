package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
	"io/ioutil"
	"math/rand"
	"net/http"
)

func randImage(query string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.google.com/search?q=%s&tbm=isch", html.EscapeString(query)))
	if err != nil {
		data, err := ioutil.ReadAll(resp.Body)
		fmt.Println(data)
		fmt.Println(err)
		return "", err
	}

	indexPage, err := htmlquery.Parse(resp.Body)
	nodes, err := htmlquery.QueryAll(indexPage, "//a/div/img")

	indexRandom := rand.Intn(len(nodes))
	node := nodes[indexRandom]
	var src string
	for _, attr := range node.Attr {
		if attr.Key == "src" {
			src = attr.Val
		}
		fmt.Println(src)
	}

	return src, nil
}
