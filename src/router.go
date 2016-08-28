package main

import (
	"log"
	"strings"

	"github.com/tucnak/telebot"
)

// Router used for flexible working with multiple commands and handlers
type Router struct {
	routes map[string]func(*telebot.Bot, telebot.Message) error
	handle func(*telebot.Bot, telebot.Message)
}

// GetRouter initializes application's router
func GetRouter() *Router {
	router := Router{
		handle: func(bot *telebot.Bot, message telebot.Message) {
			command := getCommand(message.Text)

			if handler, ok := router.routes[command]; ok {
				err := handler(bot, message)
				if err != nil {
					log.Println(message.Text + " HANDLER ERROR:")
					log.Println(err)
				}
			}
		},
		routes: make(map[string]func(*telebot.Bot, telebot.Message) error),
	}

	router.routes["/ls"] = LS
	router.routes["/actions"] = ShowActions
	router.routes["/cd"] = ChangeDirectory
	router.routes["/download"] = Download

	return &router
}

func getCommand(text string) string {
	return strings.Split(text, " ")[0]
}
