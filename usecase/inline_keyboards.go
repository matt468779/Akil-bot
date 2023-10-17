package usecase

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var mainMenuKeyboard = tgbotapi.NewInlineKeyboardButtonData("Main Menu", "mainmenu")

var mainMenuKeyboardRow = tgbotapi.NewInlineKeyboardRow(
	mainMenuKeyboard,
)

var typeOfEmploymentKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("In Person", "typeOfEmploy,inPerson"),
		tgbotapi.NewInlineKeyboardButtonData("Virtual", "typeOfEmploy,virtual"),
	),
	mainMenuKeyboardRow,
)

var categoriesKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Education", "category,Education"),
		tgbotapi.NewInlineKeyboardButtonData("Social Media", "category,Social Media"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Technology", "category,Technology"),
		tgbotapi.NewInlineKeyboardButtonData("Health", "category,Health"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Economic Development", "category,Economic Development"),
		tgbotapi.NewInlineKeyboardButtonData("Food and Nutrition", "category,Food and Nutrition"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Disaster Response ans Recovery", "category,Disaster Response ans Recovery"),
		tgbotapi.NewInlineKeyboardButtonData("Child and Maternal Health", "category,Child and Maternal Health"),
	),
	mainMenuKeyboardRow,
)
