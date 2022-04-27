package main

import (
	"fmt"
	"github.com/markcheno/go-quote"
	"time"
)

func GetQuote(symbol string) (string, error) {
	date := time.Now().Format("2006-01-02")
	spy, _ := quote.NewQuoteFromYahoo(symbol, date, date, quote.Daily, true)
	fmt.Print(spy.CSV())
	// talib.Rsi(spy.Close, 2)
	return "", nil
}
