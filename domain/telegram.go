package domain

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// A Telegram Chat indicates the conversation to which the message belongs.
type Chat struct {
	Id int `json:"id"`
}

// Message is a Telegram object that can be found in an update.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
	Date int64  `json:"date"`
	From From   `json:"from"`
}

type From struct {
	IsBot bool `json:"is_bot"`
}

// Update is a Telegram object that the handler receives every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

const (
	CollectionUpdate = "update"
)

type TelegramRepository interface {
	SaveMessage(c context.Context, update *Update) error
	GetMessages(c context.Context, update *Update) []Update
}

type TelegramUsecase interface {
	HandlePlainText(c context.Context, update Update) error
	HandleCommand(c context.Context, update tgbotapi.Update) error
	HandleCallbackQuery(c context.Context, update tgbotapi.Update) error
}
