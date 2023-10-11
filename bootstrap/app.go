package bootstrap

import (
	"akil_telegram_bot/mongo"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Application struct {
	Env   *Env
	Mongo mongo.Client
	Bot   *tgbotapi.BotAPI
}

func App() Application {
	app := &Application{}
	app.Env = NewEnv()
	app.Mongo = NewMongoDatabase(app.Env)
	app.Bot = InitTelegram(app.Env)

	return *app
}

func (app *Application) CloseDBConnection() {
	CloseMongoDBConnection(app.Mongo)
}
