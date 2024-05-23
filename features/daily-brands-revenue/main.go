package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

	log.LogInfo("Message:" + message)
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
	log.LogInfo(fmt.Sprintf("Brands: %v", brands))
	configFilePath := "/etc/config/currency.json" // Path to the mounted ConfigMap file
	curMap, err := loadCurrencyConfig(configFilePath)
	log.LogInfo(fmt.Sprintf("Currency HKD Map: %v", curMap["HKD"]))

	if err != nil {
		log.LogFatal(fmt.Sprintf("Error loading config file: %v", err))
	}

	message := "日期: " + fmt.Sprintf("%v", brands[0].Date.Format("2006-01-02")) + "<br>"

	tempCurrencyCode := ""
	for _, brand := range brands {
		log.LogInfo(fmt.Sprintf("Brand CurrencyCode: %v", brand.CurrencyCode))
		log.LogInfo(fmt.Sprintf("Brand CurrencyCode mapped: %v", curMap[brand.CurrencyCode]))

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
	return "", nil
}

// loadCurrencyConfig loads a currency configuration from a JSON file
func loadCurrencyConfig(filePath string) (map[string]string, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Read the file content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	// Initialize a map to hold the currency configuration
	currencyMap := make(map[string]string)

	// Parse the JSON content into the map
	err = json.Unmarshal(byteValue, &currencyMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}

	return currencyMap, nil
}
