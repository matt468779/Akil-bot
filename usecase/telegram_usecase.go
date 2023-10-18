package usecase

import (
	"akil_telegram_bot/bootstrap"
	"akil_telegram_bot/domain"
	"akil_telegram_bot/gpt"
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/go-querystring/query"
	"github.com/sashabaranov/go-openai"
)

type telegramUsecase struct {
	telegramRepository domain.TelegramRepository
	bot                *tgbotapi.BotAPI
	env                *bootstrap.Env
}

var err error

func NewTelegramUsecase(telegramRepository domain.TelegramRepository, timeout time.Duration, Bot *tgbotapi.BotAPI, env *bootstrap.Env) domain.TelegramUsecase {
	return &telegramUsecase{
		telegramRepository: telegramRepository,
		bot:                Bot,
		env:                env,
	}
}

func (tu *telegramUsecase) HandlePlainText(c context.Context, update domain.Update) error {
	var err error
	messages := tu.telegramRepository.GetMessages(c, &update)
	messages = append(messages, update)
	parsedMessages := parseMessages(messages)
	reply := gpt.GetResponse(parsedMessages)
	msg := tgbotapi.NewMessage(int64(update.Message.Chat.Id), reply)
	msg.ParseMode = tgbotapi.ModeHTML

	var replyUpdate domain.Update
	replyUpdate.UpdateId = update.UpdateId
	replyUpdate.Message.Chat.Id = update.Message.Chat.Id
	replyUpdate.Message.Date = time.Now().Unix()
	replyUpdate.Message.From.IsBot = true
	replyUpdate.Message.Text = reply

	if _, err = tu.bot.Send(msg); err != nil {
		panic(err)
	}

	err = tu.telegramRepository.SaveMessage(c, &update)
	if err != nil {
		panic(err)
	}
	err = tu.telegramRepository.SaveMessage(c, &replyUpdate)
	if err != nil {
		panic(err)
	}
	log.Printf("From: %+v Text: %+v\n", update.Message.From, update.Message.Text)

	return nil
}

func (tu *telegramUsecase) HandleCommand(c context.Context, update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	command := update.Message.Command()

	if command == "opportunities" {
		return startOpportunities(tu.bot, chatId, update)
	}

	return nil
}

func (tu *telegramUsecase) HandleCallbackQuery(c context.Context, update tgbotapi.Update) error {
	data := strings.Split(update.CallbackQuery.Data, ",")
	selected := data[1]

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, selected)
	if _, err := tu.bot.Request(callback); err != nil {
		panic(err)
	}

	// editMessage := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{}))
	// deleteMessage := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	// if _, err := tu.bot.Send(deleteMessage); err != nil {
	// 	println(err)
	// }

	chatId := update.CallbackQuery.Message.Chat.ID
	// callbackId := update.CallbackQuery.ID

	keyboard := update.CallbackQuery.Message.ReplyMarkup.InlineKeyboard
	prevData := *keyboard[len(keyboard)-1][0].CallbackData

	cummulatedData := prevData + "," + selected

	currState := data[0]

	if currState == "typeOfEmploy" {
		chooseCategories(cummulatedData, chatId, tu.bot)
	} else if currState == "category" {
		searchOpportunities(cummulatedData, chatId, tu.bot, tu.env.BackendURL, tu.env.FrontendURL)
	}

	return nil
}

func searchOpportunities(data string, chatId int64, bot *tgbotapi.BotAPI, backendURL string, frontendURL string) {
	searchURL := backendURL + "opportunities/search?"

	opportunityFilter := gpt.OpportunityFilter{}
	arguments := strings.Split(data, ",")
	opportunityFilter.OpportunityType = arguments[1]
	opportunityFilter.Categories = []string{arguments[2]}

	// json.Unmarshal([]byte(arguments), &opportunityFilter)

	queryString, err := query.Values(opportunityFilter)
	if err != nil {
		panic(err)
	}
	println(arguments)
	println(searchURL + queryString.Encode())
	res, err := http.Get(searchURL + queryString.Encode())
	if err != nil {
		panic(err)
	}

	result, _ := io.ReadAll(res.Body)
	var opportunities Opportunity
	json.Unmarshal(result, &opportunities)
	msg := tgbotapi.NewMessage(chatId, FormatOpportunitiesMarkdown(opportunities.Data, frontendURL))
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	if _, err = bot.Send(msg); err != nil {
		panic(err)
	}
}

func chooseCategories(data string, chatId int64, bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(chatId, "Please choose category")
	categoriesKeyboard.InlineKeyboard[len(categoriesKeyboard.InlineKeyboard)-1][0].CallbackData = &data
	msg.ReplyMarkup = categoriesKeyboard

	if _, err = bot.Send(msg); err != nil {
		panic(err)
	}
}

func (tu *telegramUsecase) PostOpportunityToChannel(username string, message string, opportunityId int) {
	msg := tgbotapi.NewMessageToChannel(username, message)
	goToBotLink := fmt.Sprintf("https://t.me/%s?start=%d", "akilconnectbot", opportunityId)
	applyButton := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Apply", goToBotLink))
	applyMarkup := tgbotapi.NewInlineKeyboardMarkup(applyButton)
	msg.ReplyMarkup = applyMarkup

	if _, err = tu.bot.Send(msg); err != nil {
		panic(err)
	}
}

func parseMessages(updates []domain.Update) []openai.ChatCompletionMessage {
	var messages []openai.ChatCompletionMessage
	for _, update := range updates {
		var message = openai.ChatCompletionMessage{}
		message.Content = update.Message.Text
		if update.Message.From.IsBot {
			message.Role = openai.ChatMessageRoleSystem
		} else {
			message.Role = openai.ChatMessageRoleUser
		}

		messages = append(messages, message)
	}

	return messages
}

func startOpportunities(bot *tgbotapi.BotAPI, chatId int64, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(chatId, "Choose Type of Employment")
	msg.ReplyMarkup = typeOfEmploymentKeyboard

	if _, err = bot.Send(msg); err != nil {
		panic(err)
	}

	return err
}
