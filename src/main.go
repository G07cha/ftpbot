package main

import (
	"log"
	"time"

	"github.com/mkideal/cli"
	"github.com/tucnak/telebot"
)

type argT struct {
	Token string `cli:"*token" usage:"enter telegram bot's token"`
}

var bot *telebot.Bot
var router *Router

func main() {
	router = GetRouter()
	cli.Run(new(argT), func(ctx *cli.Context) error {
		var err error
		argv := ctx.Argv().(*argT)

		bot, err = telebot.NewBot(argv.Token)
		if err != nil {
			return err
		}

		return nil
	})

	bot.Messages = make(chan telebot.Message)
	bot.Callbacks = make(chan telebot.Callback)

	go messages()
	go callbacks()

	log.Println("Bot started")
	bot.Start(1 * time.Second)
}

func messages() {
	for message := range bot.Messages {
		router.handle(bot, message)
	}
}

func callbacks() {
	for callback := range bot.Callbacks {
		callback.Message.Text = callback.Data
		callback.Message.Sender = callback.Sender

		// Mark callback query as readed
		bot.AnswerCallbackQuery(&callback, &telebot.CallbackResponse{
			CallbackID: callback.ID,
		})

		router.handle(bot, callback.Message)
	}
}
