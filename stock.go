package main

import (
	"fmt"
	"github.com/markcheno/go-quote"
	"io/ioutil"
	"log"
	"os"
)

func GetQuote(symbol string, date string) (string, error) {
	spy, err := quote.NewQuoteFromYahoo(symbol, date, date, quote.Daily, true)
	if err != nil {
		return "", err
	}
	return spy.CSV(), nil
}

func GetQuotes(symbols []string, date string) (string, error) {
	file, err := ioutil.TempFile(".", "tmp")
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
