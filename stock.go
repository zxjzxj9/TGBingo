package main

import (
	"fmt"
	"github.com/markcheno/go-quote"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func GetQuote(symbol string) (string, error) {
	date := time.Now().Format("2006-01-02")
	spy, err := quote.NewQuoteFromYahoo(symbol, date, date, quote.Daily, true)
	if err != nil {
		return "", err
	}
	return spy.CSV(), nil
}

func GetQuotes(symbols []string) (string, error) {
	date := time.Now().Format("2006-01-02")
	file, err := ioutil.TempFile("dir", "prefix")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(file.Name())
	for _, symbol := range symbols {
		file.WriteString(fmt.Sprintf("%s\n", symbol))
	}
	fmt.Println(file.Name())
	spy, err := quote.NewQuotesFromYahoo(file.Name(), date, date, quote.Daily, true)
	if err != nil {
		return "", err
	}
	return spy.CSV(), nil
}
