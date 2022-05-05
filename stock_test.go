package main

import "testing"

func TestGetStockQuote(t *testing.T) {
	quotes, err := GetQuote("AAPL")
	if err != nil {
		t.Errorf("Error getting stock quotes: %s", err)
	}
	if len(quotes) == 0 {
		t.Errorf("No quotes returned")
	}
}

func TestGetStockQuotes(t *testing.T) {
	quotes, err := GetQuotes([]string{"AAPL", "GOOG"})
	if err != nil {
		t.Errorf("Error getting stock quotes: %s", err)
	}
	if len(quotes) == 0 {
		t.Errorf("No quotes returned")
	}
}
