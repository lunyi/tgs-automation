package main

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"cdnetwork/pkg/letstalk"
	"cdnetwork/pkg/postgresql"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	config := util.GetConfig()

	message := getMessageFromBrandsRevenue(config.Postgresql)
	err := sendMessageToLetsTalk(config.LetsTalk, message)

	if err != nil {
		log.LogError("sendMessageToLetsTalk Error;" + err.Error())
	}

	sig := <-signals
	log.LogInfo(fmt.Sprintf("Received signal: %v, initiating shutdown", sig))
	os.Exit(0)
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

func getMessageFromBrandsRevenue(config util.PostgresqlConfig) string {
	app := postgresql.NewDailyBrandsRevenueInterface(config)

	brands, err := app.GetDailyBrandsRevenue()
	if err != nil {
		log.LogError("GetDailyBrandsRevenue Error:" + err.Error())
		panic(err)
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
			"當日營收：" + brand.DailyRevenueUSD + "<br>" +
			"當日訂單數量：" + brand.DailyOrderCount + "<br>" +
			"當日活躍人數：" + fmt.Sprintf("%d", brand.ActiveUserCount) + "<br>" +
			"當月營收：" + brand.CumulativeRevenueUSD + "<br>"
	}
	log.LogInfo(message)
	return message
}
