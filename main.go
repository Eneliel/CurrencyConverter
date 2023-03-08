package main

import (
	cur_conv "bot_mod/CC"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("6288334300:AAEbaCuAKl1UdB2aJ3xN2snU4uw-JnjBG-o")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	str := []string{"Converter", "Show all Cymbols"}
	keyb := tgbotapi.InlineKeyboardMarkup{}
	for _, i := range str {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(i, i)
		row = append(row, btn)
		keyb.InlineKeyboard = append(keyb.InlineKeyboard, row)
	}

	for update := range updates {
		if update.Message != nil && update.CallbackQuery == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			switch update.Message.Text {
			case "start":
				msg.Text = "Выберите:"
				msg.ReplyMarkup = keyb
			default:
				msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Отправте 'start' чтобы начать!")
			}
			if _, err := bot.Send(msg); err != nil {
				panic(err)
			}
		} else if update.CallbackQuery != nil {
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			if _, err := bot.Request(callback); err != nil {
				panic(err)
			}
		}
		if update.CallbackQuery != nil {
			switch update.CallbackData() {
			case "Converter":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				msg.Text = "Напишите какую валюту, куда конвертировать и колличество.\nПример:\nUSD RUB 1000\n(Мы конвертируем 1000 долларов в рубли)"
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
				for update := range updates {
					if update.Message != nil {
						msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
						str := strings.Split(update.Message.Text, " ")
						if len(str) != 3 {
							msg.Text = "Не правильный ввод, попробуйте снова"
							if _, err := bot.Send(msg); err != nil {
								panic(err)
							}
						} else {
							msg.Text = cur_conv.Conv(str[0], str[1], str[2])
							if _, err := bot.Send(msg); err != nil {
								panic(err)
							}
							break
						}
					}
				}
				msg.Text = "Выберите:"
				msg.ReplyMarkup = keyb
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			case "Show all Cymbols":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				msg.Text = cur_conv.Sym()
				msg.ReplyMarkup = keyb
				if _, err := bot.Send(msg); err != nil {
					panic(err)
				}
			}
		}
	}

}
