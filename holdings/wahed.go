package holdings

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/shopspring/decimal"
)

type rawCsvEntry struct {
	Date              string `csv:"Date"`
	Account           string `csv:"Account"`
	StockTicker       string `csv:"StockTicker"`
	CUSIP             string `csv:"CUSIP"`
	SecurityName      string `csv:"SecurityName"`
	Shares            string `csv:"Shares"`
	Price             string `csv:"Price"`
	MarketValue       string `csv:"MarketValue"`
	Weightings        string `csv:"Weightings"`
	NetAssets         string `csv:"NetAssets"`
	SharesOutstanding string `csv:"SharesOutstanding"`
	CreationUnits     string `csv:"CreationUnits"`
	MoneyMarketFlag   string `csv:"MoneyMarketFlag"`
}

func InitHLAL() {
	const ticker = "HLAL"
	const url = "https://funds.wahedinvest.com/etf-holdings.csv"
	fetchHoldings(ticker, url)
}

func InitUMMA() {
	const ticker = "UMMA"
	const url = "https://global-uploads.webflow.com/6258aa32b493a205485f0800/63a1c11d50a64e2d0303333c_Holdings_hlal.csv"
	fetchHoldings(ticker, url)
}

func fetchHoldings(ticker string, url string) {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()
	holdings := mapEntries(response.Body)
	funds[ticker] = fund{
		updatedAt: time.Now(),
		holdings:  holdings,
	}
}

func mapEntries(body io.ReadCloser) map[string]holding {
	entries := []*rawCsvEntry{}
	err := gocsv.Unmarshal(body, &entries)
	if err != nil {
		panic(err)
	}

	var holdings = make(map[string]holding)
	for _, entry := range entries {
		if entry.MoneyMarketFlag == "Y" {
			continue
		}

		weight, err := decimal.NewFromString(strings.Replace(entry.Weightings, "%", "", -1))
		if err != nil {
			panic(err)
		}

		holdings[entry.StockTicker] = holding{
			Ticker: entry.StockTicker,
			Weight: weight,
		}
	}

	return holdings
}
