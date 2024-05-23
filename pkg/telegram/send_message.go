package telegram

import (
	"strconv"
	"sync"
	"tgs-automation/internal/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	telegramBotInstance *TelegramBot
	once                sync.Once
)

type TelegramBot struct {
	Bot *tgbotapi.BotAPI
}

func NewTelegramBot(token string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{Bot: bot}, nil
}

func (tb *TelegramBot) SendMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := tb.Bot.Send(msg)
	return err
}

func getTelegramBotInstance(token string) (*TelegramBot, error) {
	var err error
	once.Do(func() {
		telegramBotInstance, err = NewTelegramBot(token)
	})
	if err != nil {
		return nil, err
	}
	return telegramBotInstance, nil
}

func SendMessage(msg string) error {
	config := util.GetConfig()
	bot, err := getTelegramBotInstance(config.Telegram.TelegramBotToken)
	if err != nil {
		return err
	}

	bot.Bot.Debug = true

	//log.Info(fmt.Sprintf("Authorized on account %s", bot.Bot.Self.UserName))

	chatid, _ := strconv.Atoi(config.Telegram.TelegramChatId)
	chatID := int64(chatid)

	if err := bot.SendMessage(chatID, msg); err != nil {
		return err
	}
	return nil
}
