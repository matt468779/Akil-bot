package usecase

import (
	"akil_telegram_bot/domain"
	"akil_telegram_bot/gpt"
	"context"
	"log"
	"time"

	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

type telegramUsecase struct {
	telegramRepository domain.TelegramRepository
	bot                *tgbotapi.BotAPI
}

func NewTelegramUsecase(telegramRepository domain.TelegramRepository, timeout time.Duration, Bot *tgbotapi.BotAPI) domain.TelegramUsecase {
	return &telegramUsecase{
		telegramRepository: telegramRepository,
		bot:                Bot,
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
	msg := tgbotapi.NewMessage(int64(update.Message.Chat.ID), "This is a command")
	msg.ParseMode = tgbotapi.ModeMarkdownV2

	_, err := tu.bot.Send(msg)
	return err
}

func (tu *telegramUsecase) HandleCallbackQuery(c context.Context, update tgbotapi.Update) error {
	editMessage := tgbotapi.NewEditMessageReplyMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{}))
	if _, err := tu.bot.Send(editMessage); err != nil {
		panic(err)
	}
	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	if _, err := tu.bot.Send(msg); err != nil {
		panic(err)
	}
	return nil
}

func (tu *telegramUsecase) PostOpportunityToChannel(username string, message string, opportunityId int) {
	msg := tgbotapi.NewMessageToChannel(username, message)
	goToBotLink := fmt.Sprintf("https://t.me/%s?start=%d", "akilconnectbot", opportunityId)
	applyButton := tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonURL("Apply", goToBotLink))
	applyMarkup := tgbotapi.NewInlineKeyboardMarkup(applyButton)
	msg.ReplyMarkup = applyMarkup
	var err error
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
