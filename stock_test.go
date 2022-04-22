package main

import "testing"

func TestGetStockQuotes(t *testing.T) {
	quotes, err := GetStockQuotes([]string{"AAPL"}, XML_FORMAT)
	if err != nil {
		t.Errorf("Error getting stock quotes: %s", err)
	}
	if len(quotes) == 0 {
		t.Errorf("No quotes returned")
	}
}
