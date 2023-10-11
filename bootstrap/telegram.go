package bootstrap

import (
	"fmt"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Bot      *tgbotapi.BotAPI
	BotToken string
	BaseURL  string
	err      error
)

func InitTelegram(env *Env) *tgbotapi.BotAPI {
	BotToken = env.BotToken
	BaseURL = env.BaseURL

	Bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s/%s", Bot.Token, BaseURL, Bot.Token)
	_, err = http.Get(url) // Set webhook
	if err != nil {
		panic(err)
	}

	return Bot
}
