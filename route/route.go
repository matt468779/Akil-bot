package route

import (
	"akil_telegram_bot/bootstrap"
	"akil_telegram_bot/mongo"
	"time"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Setup(env *bootstrap.Env, timeout time.Duration, db mongo.Database, gin *gin.Engine, bot *tgbotapi.BotAPI) {
	publicRouter := gin.Group("")

	NewTelegramRouter(env, timeout, db, publicRouter, bot)
}
