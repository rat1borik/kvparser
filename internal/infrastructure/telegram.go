package infrastructure

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Bot interface {
	SendMessage(text string) error
}

type telegramBot struct {
	bot           *tgbotapi.BotAPI
	chatIdDefault int64
}

func NewTgBot(token string, chatId int64) (Bot, error) {
	b, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	return &telegramBot{bot: b, chatIdDefault: chatId}, nil
}

func (t *telegramBot) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(t.chatIdDefault, text)
	msg.ParseMode = "HTML"
	_, err := t.bot.Send(msg)
	return err
}
