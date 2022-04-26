package main

import "testing"

func TestGetStockQuotes(t *testing.T) {
	quotes, err := GetQuote("AAPL")
	if err != nil {
		t.Errorf("Error getting stock quotes: %s", err)
	}
	if len(quotes) == 0 {
		t.Errorf("No quotes returned")
	}
}
