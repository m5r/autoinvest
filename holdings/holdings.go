package holdings

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type holding struct {
	Ticker string
	Weight decimal.Decimal
}
type fund struct {
	updatedAt time.Time
	holdings  map[string]holding
}

var funds = make(map[string]fund) // map fund $TICKER to their holdings

func GetETFHoldings(ticker string) map[string]holding {
	fund, ok := funds[ticker]
	if !ok {
		panic(fmt.Sprintf("ETF \"$%s\" not found", ticker))
	}

	return fund.holdings
}

func GetETFHoldingsHardcoded(ticker string) map[string]holding {
	return map[string]holding{
		"MSFT": {
			Ticker: "MSFT",
			Weight: decimal.NewFromFloat(0.1),
		},
		"AAPL": {
			Ticker: "AAPL",
			Weight: decimal.NewFromFloat(0.25),
		},
		"NET": {
			Ticker: "NET",
			Weight: decimal.NewFromFloat(0.4),
		},
		"SHOP": {
			Ticker: "SHOP",
			Weight: decimal.NewFromFloat(0.25),
		},
	}
}
