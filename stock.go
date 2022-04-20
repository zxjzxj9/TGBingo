package main

// from https://gist.github.com/mickelsonm/ecbfe59979e71f075995
import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	API_URL        = "http://query.yahooapis.com/v1/public/yql?q="
	queryString    = "select * from yahoo.finance.quotes where symbol in ('%s')"
	configSettings = "&format=%s&env=http://datatables.org/alltables.env"

	JSON_FORMAT = 0
	XML_FORMAT  = 1
)

type Quotes []Quote
type Quote struct {
	Name               string `json:"Name" xml:"Name"`
	Symbol             string `json:"Symbol" xml:"symbol,attr"`
	StockExchange      string `json:"StockExchange" xml:"StockExchange"`
	Currency           string `json:"Currency" xml:"Currency"`
	Volume             string `json:"Volume" xml:"Volume"`
	AverageDailyVolume string `json:"AverageDailyVolume" xml:"AverageDailyVolume"`

	Ask         string `json:"Ask" xml:"Ask"`
	AskRealtime string `json:"AskRealtime" xml:"AskRealtime"`
	BidRealtime string `json:"BidRealtime" xml:"BidRealtime"`

	Open          string `json:"Open" xml:"Open"`
	PreviousClose string `json:"PreviousClose" xml:"PreviousClose"`
	LastTradeDate string `json:"LastTradeDate" xml:"LastTradeDate"`

	DaysLow           string `json:"DaysLow" xml:"DaysLow"`
	DaysHigh          string `json:"DaysHigh" xml:"DaysHigh"`
	DaysRange         string `json:"DaysRange" xml:"DaysRange"`
	DaysRangeRealtime string `json:"DaysRangeRealtime" xml:"DaysRangeRealtime"`

	YearLow   string `json:"YearLow" xml:"YearLow"`
	YearHigh  string `json:"YearHigh" xml:"YearHigh"`
	YearRange string `json:"YearRange" xml:"YearRange"`

	//TODO: there's a lot more fields to add if they are needed
}

type YahooQuoteResponse struct {
	Query struct {
		Results struct {
			Quotes Quotes `json:"quote" xml:"quote"`
		}
	}
}

type YahooSingleQuoteResponse struct {
	Query struct {
		Results struct {
			Quote Quote `json:"quote" xml:"quote"`
		}
	}
}

func GetStockQuotes(symbols []string, dataFormat ...int) (quotes Quotes, err error) {
	var buf []byte
	var resp *http.Response
	var symbolsString string
	var urlString string
	var currentDataFormat int

	//handle the stock symbols
	if len(symbols) == 0 {
		err = fmt.Errorf("Must have at least one symbol.")
		return
	} else if len(symbols) > 1 {
		for _, s := range symbols {
			symbolsString += s + ","
		}
	} else {
		symbolsString = symbols[0]
	}

	//handle the data format, defaults to json
	urlString = API_URL + url.QueryEscape(fmt.Sprintf(queryString, symbolsString))
	if len(dataFormat) > 0 && dataFormat[0] == XML_FORMAT {
		currentDataFormat = XML_FORMAT
		urlString += fmt.Sprintf(configSettings, "xml")
	} else {
		currentDataFormat = JSON_FORMAT
		urlString += fmt.Sprintf(configSettings, "json")
	}

	//fmt.Println(urlString)
	if resp, err = http.Get(urlString); err != nil {
		return
	}
	defer resp.Body.Close()

	if buf, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	//unmarshal the json/xml response
	var yqr YahooQuoteResponse
	var yqrs YahooSingleQuoteResponse
	if currentDataFormat == XML_FORMAT {
		if len(symbols) > 1 {
			if err = xml.Unmarshal(buf, &yqr); err != nil {
				return
			}
		} else {
			if err = xml.Unmarshal(buf, &yqrs); err != nil {
				return
			}
		}
	} else {
		if len(symbols) > 1 {
			if err = json.Unmarshal(buf, &yqr); err != nil {
				return
			}
		} else {
			if err = json.Unmarshal(buf, &yqrs); err != nil {
				return
			}
		}
	}

	//handle the single response or multiple responses
	if len(symbols) > 1 {
		quotes = yqr.Query.Results.Quotes
	} else {
		quotes = Quotes{yqrs.Query.Results.Quote}
	}

	return
}
