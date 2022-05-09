package main

import (
	"fmt"
	"testing"
)

func TestGetStockQuote(t *testing.T) {
	quotes, err := GetQuote("AAPL", "2022-05-06")
	if err != nil {
		t.Errorf("Error getting stock quotes: %s", err)
	}
	if len(quotes) == 0 {
		t.Errorf("No quotes returned")
	}
}

func TestGetStockQuotes(t *testing.T) {
	quotes, err := GetQuotes([]string{"AAPL", "GOOG"}, "2022-05-01")
	if err != nil {
		t.Errorf("Error getting stock quotes: %s", err)
	}
	if len(quotes) == 0 {
		t.Errorf("No quotes returned")
	}
	fmt.Println(quotes)
}
