package main

import (
	"cdnetwork/internal/util"
	"cdnetwork/pkg/postgresql"
	"fmt"
)

func main() {
	config := util.GetConfig()
	app := postgresql.NewBrandsIncomeInterface(config.Postgresql)

	brands, err := app.GetBrandsIncome()
	if err != nil {
		panic(err)
	}

	for _, brand := range brands {
		fmt.Printf("PlatformCode: %s, CurrencyCode: %s, Date: %s, ActiveUserCount: %d, DailyOrderCount: %s, DailyRevenueUSD: %.2f, CumulativeRevenueUSD: %.2f\n",
			brand.PlatformCode,
			brand.CurrencyCode,
			brand.Date.Format("2006-01-02"), // Formatting date as YYYY-MM-DD
			brand.ActiveUserCount,
			brand.DailyOrderCount,
			brand.DailyRevenueUSD,
			brand.CumulativeRevenueUSD,
		)
	}
}
