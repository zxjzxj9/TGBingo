package main

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
)

// Scraper for books
var (
	header http.Header
	wg sync.WaitGroup
)

func init() {
	header := make(http.Header)
	header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	header.Add("accept-language", "zh-CN,zh;q=0.9,en;q=0.8")
	header.Add("sec-ch-ua", "\"Chromium\";v=\"88\", \"Google Chrome\";v=\"88\", \";Not A Brand\";v=\"99\"")
	header.Add("sec-ch-ua-mobile", "?0")
	header.Add("sec-fetch-dest", "document")
	header.Add("sec-fetch-mode", "navigate")
	header.Add("sec-fetch-site", "none")
	header.Add("sec-fetch-user", "?1")
	header.Add("upgrade-insecure-requests", "1")
}

func crawl(pageId int, c chan <- float32) (string, error) {
	// defer wg.Done()
	proxyUrl, err := url.Parse("http://127.0.0.1:8118")

	if err != nil {
		fmt.Println("Failed to create proxy ", err)
		return "", err
	}
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	// client := &http.Client{}
	fmt.Println(fmt.Sprintf("https://cn.uukanshu.cc/book/%d/", pageId))
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://cn.uukanshu.cc/book/%d/", pageId), nil)
	req.Header = header
	if err != nil {
		fmt.Println("Fail to create request ", err)
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Fail to retrieve the book page ", err)
		return "", err
	}

	// data, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(resp.StatusCode, string(data))
	indexPage, err := htmlquery.Parse(resp.Body)
	if err != nil {
		fmt.Println("not a valid XPath expression.")
		return "", err
	}

	fname := fmt.Sprintf("cache/%d.txt", pageId)
	fout, err := os.OpenFile(fname, os.O_WRONLY | os.O_CREATE, 0755)
	defer fout.Close()

	// fmt.Println(htmlquery.InnerText(indexPage))
	nodes, err := htmlquery.QueryAll(indexPage, "//*[@id='list-chapterAll']/div/dd")
	if err != nil {
		fmt.Println("Nodes query failed")
		return "", err
	}

	for idx, node := range nodes {
		p, err := htmlquery.Query(node, "a[1]/@href")
		if err != nil {
			fmt.Println("Nodes query failed")
			return "", err
		}

		// fmt.Println(idx)
		c <- float32(idx+1)/float32(len(nodes))

		url := htmlquery.InnerText(p)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		req.Header = header
		if err != nil {
			fmt.Println("Fail to create request ", err)
			return "", err
		}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Fail to retrieve the book page ", err)
			return "", err
		}

		dataPage, err := htmlquery.Parse(resp.Body)
		if err != nil {
			fmt.Println("not a valid XPath expression.")
			return "", err
		}

		head, err := htmlquery.Query(dataPage, "/html/body/div[2]/div[1]/div/h1/text()")
		if err != nil {
			fmt.Println("Nodes query failed")
			return "", err
		}
		// fmt.Println(htmlquery.InnerText(head))
		fout.WriteString(htmlquery.InnerText(head) + "\n")

		contentNodes, err := htmlquery.QueryAll(dataPage, "/html/body/div[2]/div[1]/div/p[2]/text()")
		if err != nil {
			fmt.Println("Nodes query failed")
			return "", err
		}

		for _, content := range contentNodes {
			// fmt.Println(htmlquery.InnerText(content))
			fout.WriteString(strings.Trim(htmlquery.InnerText(content), " \n"))
		}
		fout.WriteString("\n")
	}

	return fname, nil
}
