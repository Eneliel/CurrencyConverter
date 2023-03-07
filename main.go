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

	for update := range updates {
		if update.Message != nil { // If we got a message
			str := strings.Split(update.Message.Text, " ")
			in_str := strings.ToLower(str[0])
			switch in_str {
			case "conv":
				from := strings.ToUpper(str[1])
				to := strings.ToUpper(str[2])
				amount := str[3]
				out := cur_conv.Conv(from, to, amount)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, out))
			case "symb":
				m := cur_conv.Sym()
				var out string
				for k, v := range m {
					out += k + " : " + v + "\n"
				}
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, out))
			default:
			}

		}
	}
}
