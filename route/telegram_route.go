package route

import (
	"akil_telegram_bot/bootstrap"
	"akil_telegram_bot/controller"
	"akil_telegram_bot/domain"
	"akil_telegram_bot/mongo"
	"akil_telegram_bot/repository"
	"akil_telegram_bot/usecase"
	"time"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewTelegramRouter(env *bootstrap.Env, timeout time.Duration, db mongo.Database, group *gin.RouterGroup, bot *tgbotapi.BotAPI) {
	tr := repository.NewTelegramRepository(db, domain.CollectionUpdate)
	tc := &controller.TelegramController{
		TelegramUsecase: usecase.NewTelegramUsecase(tr, timeout, bot, env),
	}
	group.POST("/"+env.BotToken, tc.HandleWebhook)
}
