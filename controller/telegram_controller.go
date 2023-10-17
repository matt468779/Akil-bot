package controller

import (
	"akil_telegram_bot/domain"
	"encoding/json"
	"fmt"
	"io"
	"log"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramController struct {
	TelegramUsecase domain.TelegramUsecase
}

func (tc *TelegramController) HandleWebhook(c *gin.Context) {
	var update tgbotapi.Update
	var parseUpdate domain.Update
	defer c.Request.Body.Close()

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(bytes, &parseUpdate)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("%s\n", bytes)
	if update.Message != nil {
		if update.Message.IsCommand() {
			tc.TelegramUsecase.HandleCommand(c, update)
		} else {
			tc.TelegramUsecase.HandlePlainText(c, parseUpdate)
		}
	} else if update.CallbackQuery != nil {
		tc.TelegramUsecase.HandleCallbackQuery(c, update)
	}
}
