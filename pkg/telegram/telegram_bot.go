package telegram

import (
	"context"
	"strconv"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type TelegramApi interface {
	SendMessage(chatID int64, text string) error
	SendMessageWithChatIdAndContext(ctx context.Context, chatid string, msg string) error
}
type TelegramBot struct {
	Token string
	Bot   *tgbotapi.BotAPI
}

var (
	telegramBotInstance *TelegramBot
	once                sync.Once
)

func newTelegramBot(token string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{Bot: bot, Token: token}, nil
}

func (tb *TelegramBot) SendMessage(chatid string, text string) error {
	intChatid, err := strconv.Atoi(chatid)
	if err != nil {
		return err
	}
	intChatId := int64(intChatid)
	msg := tgbotapi.NewMessage(intChatId, text)
	_, err = tb.Bot.Send(msg)
	return err
}

func (tb *TelegramBot) SendMessageWithChatIdAndContext(ctx context.Context, chatid string, msg string) error {
	if err := tb.SendMessage(chatid, msg); err != nil {
		return err
	}
	return nil
}

func GetTelegramBotInstance(token string) (*TelegramBot, error) {
	var err error
	once.Do(func() {
		telegramBotInstance, err = newTelegramBot(token)
	})
	if err != nil {
		return nil, err
	}
	return telegramBotInstance, nil
}
