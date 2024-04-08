package main

import (
	"cdnetwork/internal/log"
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

	message := ""

	for _, brand := range brands {
		message += "\n" + brand.PlatformCode + "\n" +
			"當日營收：" + brand.DailyRevenueUSD + "\n" +
			"當日訂單數量：" + brand.DailyOrderCount + "\n" +
			"當日活躍人數：" + fmt.Sprintf("%d", brand.ActiveUserCount) + "\n" +
			"當月營收：" + brand.CumulativeRevenueUSD + "\n\n"

	}
	log.LogInfo(message)
}
