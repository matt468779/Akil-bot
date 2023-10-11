package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	Bot       *tgbotapi.BotAPI
	BotToken  = os.Getenv("BOT_TOKEN")
	BaseURL   = os.Getenv("BASE_URL")
	ChannelID = -1001742144512
)

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://t.me/akilconnectbot?start=12345"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
)

func InitTelegram() {
	var err error

	Bot, err = tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Println(err)
		return
	}

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
	url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s/%s", Bot.Token, BaseURL, Bot.Token)
	_, err = http.Get(url) // Set webhook
	if err != nil {
		log.Println("webhook not set: ", err)
	}
}

func WebhookHandler(c *gin.Context) error {
	defer c.Request.Body.Close()

	bytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return nil
	}
	if update.Message.IsCommand() {
		if update.Message.Command() == "/setprofile" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is a regular message with some hidden text: [hidden](invisible)")
			msg.ParseMode = tgbotapi.ModeMarkdownV2

			Bot.Send(msg)
		}
	} else if update.Message != nil {
		// reply := gpt.GetResponse(update.Message.Text)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is a regular message with some hidden text: [hidden](invisible)")
		msg.ParseMode = tgbotapi.ModeMarkdownV2

		// PostOpportunityToChannel("@akilconnect", "This is just a test opportunity", 12345)
		// EditOpportunityOnChannel(int64(ChannelID), 20, "This edited test message", false)
		if _, err = Bot.Send(msg); err != nil {
			panic(err)
		}
		log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)

	} else if update.CallbackQuery != nil {
		editMessage := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{}))
		if _, err := Bot.Send(editMessage); err != nil {
			panic(err)
		}
		msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		if _, err := Bot.Send(msg); err != nil {
			panic(err)
		}
	}

	return nil
}

func PostOpportunityToChannel(username string, message string, opportunityId int) {
	msg := tgbotapi.NewMessageToChannel(username, message)
	goToBotLink := fmt.Sprintf("https://t.me/%s?start=%d", "akilconnectbot", opportunityId)
	applyButton := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Apply", goToBotLink))
	applyMarkup := tgbotapi.NewInlineKeyboardMarkup(applyButton)
	msg.ReplyMarkup = applyMarkup
	var res tgbotapi.Message
	var err error
	if res, err = Bot.Send(msg); err != nil {
		panic(err)
	}
	println("??????????????????????????")
	println(res.MessageID)
}

func EditOpportunityOnChannel(chatId int64, messageId int, message string, is_active bool) {
	var msg tgbotapi.EditMessageTextConfig
	if is_active {
		msg = tgbotapi.NewEditMessageText(chatId, messageId, message)
	} else {
		emptyMarkup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{})
		msg = tgbotapi.NewEditMessageTextAndMarkup(chatId, messageId, message, emptyMarkup)
	}

	if _, err := Bot.Send(msg); err != nil {
		panic(err)
	}
}
