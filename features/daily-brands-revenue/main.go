package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"tgs-automation/internal/log"
	"tgs-automation/internal/util"
	"tgs-automation/pkg/letstalk"
	"tgs-automation/pkg/postgresql"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go Run(ctx)

	<-ctx.Done() // Wait for signal
	log.LogInfo("Shutting down...")
}

func Run(ctx context.Context) {
	config := util.GetConfig()

	message, err := getMessageFromBrandsRevenue(config.Postgresql)
	if err != nil {
		log.LogError("getMessageFromBrandsRevenue Error;" + err.Error())
	}

	err = sendMessageToLetsTalk(config.LetsTalk, message)

	if err != nil {
		log.LogError("sendMessageToLetsTalk Error;" + err.Error())
	}
}

func sendMessageToLetsTalk(config util.LetsTalkConfig, message string) error {
	token, err := letstalk.GetToken(config)
	if err != nil {
		log.LogError("Get Token Error:" + token)
		return err
	}
	var rooms []letstalk.Room
	rooms, err = letstalk.GetRooms(token)
	if err != nil {
		log.LogError("Get Room Error:" + err.Error())
		return err
	}

	for _, room := range rooms {
		log.LogInfo("Room:" + room.Title + " Token:" + room.Token)
	}

	letstalkChatGroupTitles := []string{"PG Daily Report"}
	roomTokens := []string{}
	for _, room := range rooms {
		for _, key := range letstalkChatGroupTitles {
			if room.Title == key {
				roomTokens = append(roomTokens, room.Token)
			}
		}
	}

	err = letstalk.SendMessage(token, roomTokens, message)
	if err != nil {
		log.LogError("SendMessage Error:" + err.Error())
		return err
	}
	return nil
}

func getMessageFromBrandsRevenue(config util.PostgresqlConfig) (string, error) {
	app := postgresql.NewDailyBrandsRevenueInterface(config)

	brands, err := app.GetDailyBrandsRevenue()
	if err != nil {
		log.LogError("GetDailyBrandsRevenue Error:" + err.Error())
		return "", err

	}

	curMap := map[string]string{
		"PHP":      "PHP - 菲律賓幣",
		"HKD":      "HKD - 港幣",
		"VND_1000": "VND_1000 - 越南盾",
	}

	message := "日期: " + fmt.Sprintf("%v", brands[0].Date.Format("2006-01-02")) + "<br>"

	tempCurrencyCode := ""
	for _, brand := range brands {
		if tempCurrencyCode != brand.CurrencyCode {
			message += "<br><b>[" + curMap[brand.CurrencyCode] + "]</b>"
			tempCurrencyCode = brand.CurrencyCode
		}
		message += "<br>[" + brand.PlatformCode + "]<br>" +
			"當日訂單數量：" + brand.DailyOrderCount + "<br>" +
			"當日活躍人數：" + fmt.Sprintf("%d", brand.ActiveUserCount) + "<br>" +
			"當日營收：$ " + brand.DailyRevenueUSD + "<br>" +
			"當月營收：$ " + brand.CumulativeRevenueUSD + "<br>"
	}
	log.LogInfo(message)
	return message, nil
}
