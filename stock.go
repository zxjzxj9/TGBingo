package main

import (
	"github.com/markcheno/go-quote"
	"time"
)

func GetQuote(symbol string) (string, error) {
	date := time.Now().Format("2006-01-02")
	spy, err := quote.NewQuoteFromYahoo(symbol, date, date, quote.Daily, true)
	if err != nil {
		return "", err
	}
	// fmt.Print(spy.CSV())
	// talib.Rsi(spy.Close, 2)
	return spy.CSV(), nil
}
