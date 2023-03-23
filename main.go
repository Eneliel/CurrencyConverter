package main

import (
	cur_conv "bot_mod/CC"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
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
	str := []string{"Converter", "Show all Cymbols", "Wallet"}
	wallet_keyb := []string{"Добавить", "Убавить", "Назад"}
	keyb := tgbotapi.InlineKeyboardMarkup{}
	for _, i := range str {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(i, i)
		row = append(row, btn)
		keyb.InlineKeyboard = append(keyb.InlineKeyboard, row)
	}

	wal_keyb := tgbotapi.InlineKeyboardMarkup{}
	for _, i := range wallet_keyb {
		var row []tgbotapi.InlineKeyboardButton
		btn := tgbotapi.NewInlineKeyboardButtonData(i, i)
		row = append(row, btn)
		wal_keyb.InlineKeyboard = append(wal_keyb.InlineKeyboard, row)
	}
	//&& update.CallbackQuery == nil
	for update := range updates {
		if update.Message != nil && update.CallbackQuery == nil && !update.Message.IsCommand() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.Text = "Выберите в меню /start, чтобы начать"
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
			case "Wallet":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				id := update.CallbackQuery.Message.Chat.ID
				msg.ReplyMarkup = wal_keyb
				bot.Send(msg)
				for update2 := range updates {
					if update2.CallbackQuery != nil {
						switch update2.CallbackData() {
						case "Добавить":
							msg.ReplyMarkup = nil
							msg.Text = "В какой валюте хотите добавить и сколько\nВся валюта будет конвертированна в рубли"
							bot.Send(msg)
							for update3 := range updates {
								if update3.Message != nil {
									msg = tgbotapi.NewMessage(update3.Message.Chat.ID, update3.Message.Text)
									str := strings.Split(update3.Message.Text, " ")
									if len(str) != 2 {
										msg.Text = "Не правильный ввод, попробуйте ещё раз\nПример: USD 1000"
										bot.Send(msg)
									} else {
										Conv := cur_conv.Conv(str[0], "RUB", str[1])
										msg.Text = "На вашем кошельке:\n" + cur_conv.Wallet(Conv, "sum", id) + " рублей"
										msg.ReplyMarkup = keyb
										bot.Send(msg)
										break
									}
								}
							}
						case "Убавить":
							msg.Text = "В какой валюте хотите убавить и сколько\nВся валюта будет конвертированна в рубли"
							bot.Send(msg)
							for update3 := range updates {
								if update3.Message != nil {
									msg = tgbotapi.NewMessage(update3.Message.Chat.ID, update3.Message.Text)
									str := strings.Split(update3.Message.Text, " ")
									if len(str) != 2 {
										msg.Text = "Не правильный ввод, попробуйте ещё раз\nПример: USD 1000"
										bot.Send(msg)
									} else {
										Conv := cur_conv.Conv(str[0], "RUB", str[1])
										msg.Text = "На вашем кошельке:\n" + cur_conv.Wallet(Conv, "min", id) + " рублей"
										msg.ReplyMarkup = keyb
										bot.Send(msg)
										break
									}
								}
							}
						case "Назад":
							msg.Text = "Выберите:"
							msg.ReplyMarkup = keyb
							bot.Send(msg)

						}
					}
					break
				}
			}
		} else if update.Message.IsCommand() && update.CallbackQuery == nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "start":
				msg.Text = "Выберите:"
				msg.ReplyMarkup = keyb
			}
			bot.Send(msg)
		}
	}

}
