package telegram

import (
	"cdnetwork/internal/log"
	"cdnetwork/internal/util"
	"fmt"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

func SendMessage(msg string) {

	config := util.GetConfig()

	bot, err := NewTelegramBot(config.Telegram.TelegramBotToken)
	if err != nil {
		log.LogFatal(err.Error())
	}

	bot.Bot.Debug = true

	log.LogError(fmt.Sprintf("Authorized on account %s", bot.Bot.Self.UserName))

	chatid, _ := strconv.Atoi(config.Telegram.TelegramChatId)
	chatID := int64(chatid)

	if err := bot.SendMessage(chatID, msg); err != nil {
		log.LogFatal(err.Error())
	}
}
