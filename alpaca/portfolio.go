package alpaca

import (
	"fmt"

	"github.com/alpacahq/alpaca-trade-api-go/v2/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v2/marketdata"
	"github.com/shopspring/decimal"
)

type Portfolio struct {
	stocks map[string]*Stock
}

type Stock struct {
	symbol string
	id     string
	price  decimal.Decimal
	shares decimal.Decimal
}

var client = alpaca.NewClient(alpaca.ClientOpts{
	ApiKey:    "PKESJMK8ZJ35XDVJ1N3B",
	ApiSecret: "k8KIlcoaayeMIzW83rfLIDXYB84zxCkzijDQQrBm",
	BaseURL:   "https://paper-api.alpaca.markets",
})
var marketDataClient = marketdata.NewClient(marketdata.ClientOpts{
	ApiKey:    "PKESJMK8ZJ35XDVJ1N3B",
	ApiSecret: "k8KIlcoaayeMIzW83rfLIDXYB84zxCkzijDQQrBm",
})

func (p *Portfolio) FetchPositions() {
	// Retrieve the list of positions for your account
	positions, err := client.ListPositions()
	if err != nil {
		panic(fmt.Errorf("error retrieving positions: %v", err))
	}

	// Update the portfolio with the current positions from Alpaca
	p.stocks = make(map[string]*Stock)
	for _, position := range positions {
		p.stocks[position.Symbol] = &Stock{
			symbol: position.Symbol,
			id:     position.AssetID,
			price:  *position.CurrentPrice,
			shares: position.Qty,
		}
	}
}

func (p *Portfolio) Rebalance(amount decimal.Decimal, targetAllocation map[string]decimal.Decimal) {
	totalValue := decimal.New(0, 0)
	allocations := make(map[string]decimal.Decimal)

	// Calculate the total value of the portfolio and the current allocation of each stock
	for _, stock := range p.stocks {
		value := stock.price.Mul(stock.shares)
		totalValue = totalValue.Add(value)
		allocations[stock.symbol] = value.Div(totalValue)
	}

	// TODO: what to do when `portfolio` holds something not present in `targetAllocation`

	// Calculate the number of shares needed for each stock to reach the target allocation
	for symbol, target := range targetAllocation {
		current := allocations[symbol]
		difference := target.Sub(current)
		if difference.Equal(decimal.New(0, 0)) {
			continue
		}

		// Find the stock in the portfolio
		stock, ok := p.stocks[symbol]
		if !ok {
			// Stock is not in the portfolio, so create a new stock object
			asset, err := client.GetAsset(symbol)
			if err != nil {
				panic(fmt.Errorf("error retrieving asset \"%s\": %v", symbol, err))
			}

			quotes, err := marketDataClient.GetLatestQuote(symbol)
			if err != nil {
				panic(fmt.Errorf("error retrieving asset's quotes \"%s\": %v", symbol, err))
			}

			stock = &Stock{
				symbol: symbol,
				id:     asset.ID,
				price:  decimal.NewFromFloat(quotes.AskPrice),
				shares: decimal.New(0, 0),
			}
			p.stocks[symbol] = stock
		}

		// Calculate the number of shares to buy or sell
		fmt.Println(stock.price)
		shares := amount.Mul(difference).Div(stock.price)
		fmt.Println(shares)
		if shares.LessThanOrEqual(decimal.New(0, 0)) {
			continue
		}

		// Place the trade using the Alpaca API (or other mechanism)
		order, err := client.PlaceOrder(alpaca.PlaceOrderRequest{
			AssetKey:    &symbol,
			Qty:         &shares,
			Type:        "market",
			Side:        "buy",
			TimeInForce: "day",
		})
		fmt.Println(order)
		if err != nil {
			panic(fmt.Errorf("error placing buy order for %q: %v", symbol, err))
		}

		// Update the number of shares in the portfolio
		stock.shares = stock.shares.Add(shares)
	}
}
