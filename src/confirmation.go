package main

import (
	"fmt"
	"strings"

	"github.com/tucnak/telebot"
)

// Confirm action
func Confirm(action string) func(*telebot.Bot, telebot.Message) error {
	return func(bot *telebot.Bot, msg telebot.Message) error {
		item := msg.Text[strings.Index(msg.Text, " ")+1:]
		message := fmt.Sprintf("Are you sure that you want to %s %s?", action, item)

		markup := telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.KeyboardButton{
				[]telebot.KeyboardButton{
					telebot.KeyboardButton{
						Text: "Yes",
						Data: fmt.Sprintf("/%s %s", action, item),
					},
					telebot.KeyboardButton{
						Text: "No",
						Data: "/cancel",
					},
				},
			},
		}

		return bot.SendMessage(msg.Chat, message, &telebot.SendOptions{
			ReplyMarkup: markup,
		})
	}
}
