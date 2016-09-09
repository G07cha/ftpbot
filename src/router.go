package main

import (
	"log"
	"strings"

	"github.com/tucnak/telebot"
)

type router struct {
	routes map[string]func(*telebot.Bot, telebot.Message) error
	handle func(*telebot.Bot, telebot.Message)
}

// Router used for flexible working with multiple commands and handlers
var Router = router{
	routes: make(map[string]func(*telebot.Bot, telebot.Message) error),
}

func init() {
	Router.handle = handler
}

func handler(bot *telebot.Bot, message telebot.Message) {
	command := getCommand(message.Text)

	if message.Document.FileName != "" && Router.routes["file"] != nil {
		Router.routes["file"](bot, message)
	}

	if length := len(command); length == 0 && Router.routes["text"] != nil {
		Router.routes["text"](bot, message)
	}

	if handler, ok := Router.routes[command]; ok {
		err := handler(bot, message)
		if err != nil {
			log.Println(message.Text + " HANDLER ERROR:")
			log.Println(err)
		}
	}
}

func getCommand(text string) string {
	if text[0] == '/' {
		spaceSymbol := strings.Index(text, " ")
		if spaceSymbol > -1 {
			return text[:spaceSymbol]
		}
		return text
	}
	return ""
}
