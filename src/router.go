package main

import (
	"log"
	"regexp"
	"strings"

	"github.com/tucnak/telebot"
)

// Router used for flexible working with multiple commands and handlers
type Router struct {
	routes map[*regexp.Regexp]func(*telebot.Bot, telebot.Message) error
	handle func(*telebot.Bot, telebot.Message)
}

// GetRouter initializes application's router
func GetRouter() *Router {
	router := Router{
		handle: func(bot *telebot.Bot, message telebot.Message) {
			for route, handler := range router.routes {
				if route.MatchString(message.Text) == true {
					err := handler(bot, message)
					if err != nil {
						log.Println(message.Text + " HANDLER ERROR:")
						log.Println(err)
					}
				}
			}

		},
		routes: make(map[*regexp.Regexp]func(*telebot.Bot, telebot.Message) error),
	}

	routes := map[string]func(*telebot.Bot, telebot.Message) error{
		"^/ls*":       LS,
		"^/actions*":  ShowActions,
		"^/cd*":       ChangeDirectory,
		"^/download*": Download,
		"^/rm*":       Confirm("delete"),
		"^/delete*":   Remove,
		"^/cancel*":   ResetAction,
	}

	for route, handler := range routes {
		expression, err := regexp.Compile(route)
		if err != nil {
			log.Println(route + " is invalid regexp, skipping")
			log.Println(err)
		} else {
			router.routes[expression] = handler
		}
	}

	return &router
}

func getCommand(text string) string {
	return strings.Split(text, " ")[0]
}
