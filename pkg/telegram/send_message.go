package telegram

import (
	"context"
	"tgs-automation/internal/util"
)

func getTelegramBotInstance(token string) (*TelegramBot, error) {
	var err error
	once.Do(func() {
		telegramBotInstance, err = newTelegramBot(token)
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
	if err := bot.SendMessage(config.Telegram.TelegramChatId, msg); err != nil {
		return err
	}
	return nil
}

func SendMessageWithChatId(msg string, chatid string) error {
	config := util.GetConfig()
	bot, err := getTelegramBotInstance(config.Telegram.TelegramBotToken)
	if err != nil {
		return err
	}

	bot.Bot.Debug = true
	if err := bot.SendMessage(chatid, msg); err != nil {
		return err
	}
	return nil
}

func SendMessageWithChatIdAndContext(ctx context.Context, msg string, chatid string) error {
	config := util.GetConfig()
	bot, err := getTelegramBotInstance(config.Telegram.TelegramBotToken)
	if err != nil {
		return err
	}

	bot.Bot.Debug = true
	if err := bot.SendMessage(chatid, msg); err != nil {
		return err
	}
	return nil
}
