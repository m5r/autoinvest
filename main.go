package main

import (
	"autoinvest/holdings"
	"fmt"
)

func main() {
	//	<init>
	//		fetch portfolio positions
	//		fetch etf composition
	//		filter out unavailable stocks
	//		filter out non-shariah-compliant stocks
	//	</init>
	/* portfolio := &alpaca.Portfolio{}
	portfolio.FetchPositions()

	MIALholdings := holdings.GetETFHoldingsHardcoded("MIAL")
	targetallocation := make(map[string]decimal.Decimal)
	for ddd, holding := range MIALholdings {
		targetallocation[ddd] = holding.Weight
	}

	portfolio.Rebalance(decimal.NewFromInt(1000), targetallocation) */

	holdings.InitHLAL()
	holdings.InitUMMA()
	fmt.Println(holdings.GetETFHoldings("UMMA"))
}
