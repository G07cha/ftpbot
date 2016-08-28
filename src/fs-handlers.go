package main

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/tucnak/telebot"
)

// LS provide interface between telebot and native function to list files in current directory
func LS(bot *telebot.Bot, msg telebot.Message) error {
	page := 0
	path := GetCurrentState(&msg.Sender).currentPath
	args := strings.Split(msg.Text, " ")
	if len(args) > 1 {
		page, _ = strconv.Atoi(args[1])
	}

	markup, err := lsToMarkup(path, page)
	if err != nil {
		return err
	}

	return bot.SendMessage(msg.Chat, "List of items in "+path, &telebot.SendOptions{
		ReplyMarkup: markup,
	})
}

// ShowActions determinates and show possible actions that could be done with file from passed query
func ShowActions(bot *telebot.Bot, msg telebot.Message) error {
	filename := msg.Text[strings.Index(msg.Text, " ")+1:]
	currentPath := GetCurrentState(&msg.Sender).currentPath
	file, err := os.Open(path.Join(currentPath, filename))
	if err != nil {
		return err
	}

	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	folderActions := []telebot.KeyboardButton{
		telebot.KeyboardButton{
			Text: "Open",
			Data: "/cd " + fileInfo.Name(),
		},
	}
	fileActions := []telebot.KeyboardButton{
		telebot.KeyboardButton{
			Text: "Download",
			Data: "/download " + fileInfo.Name(),
		},
	}

	// Setup markup with base actions
	replyMarkup := &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.KeyboardButton{
			[]telebot.KeyboardButton{
				telebot.KeyboardButton{
					Text: "Delete",
					Data: "/confirm delete " + fileInfo.Name(),
				},
			},
			[]telebot.KeyboardButton{
				telebot.KeyboardButton{
					Text: "Copy",
					Data: "/cp " + fileInfo.Name(),
				},
				telebot.KeyboardButton{
					Text: "Move",
					Data: "/mv " + fileInfo.Name(),
				},
			},
			[]telebot.KeyboardButton{
				telebot.KeyboardButton{
					Text: "Rename",
					Data: "/rename " + fileInfo.Name(),
				},
			},
		},
	}

	var selectedActions []telebot.KeyboardButton

	if fileInfo.IsDir() == true {
		selectedActions = folderActions
	} else {
		selectedActions = fileActions
	}

	replyMarkup.InlineKeyboard = append(replyMarkup.InlineKeyboard, selectedActions)

	return bot.SendMessage(msg.Chat, "Choose action for "+fileInfo.Name(), &telebot.SendOptions{
		ReplyMarkup: *replyMarkup,
	})
}

// ChangeDirectory updates current working directory for users and return file listing in new directory
func ChangeDirectory(bot *telebot.Bot, msg telebot.Message) error {
	state := GetCurrentState(&msg.Sender)

	newPath := path.Join(state.currentPath, strings.Split(msg.Text, " ")[1])
	state.currentPath = newPath

	markup, err := lsToMarkup(newPath, 0)
	if err != nil {
		return err
	}

	return bot.SendMessage(msg.Chat, "You are now in "+newPath, &telebot.SendOptions{
		ReplyMarkup: markup,
	})
}

// Download used for downloading files from fs
func Download(bot *telebot.Bot, msg telebot.Message) error {
	filename := msg.Text[strings.Index(msg.Text, " ")+1:]
	fileExt := filename[strings.LastIndex(filename, ".")+1:]

	file, err := telebot.NewFile(path.Join(GetCurrentState(&msg.Sender).currentPath, filename))
	if err != nil {
		return err
	}

	switch {
	case fileExt == "png" || fileExt == "jpg":
		bot.SendPhoto(msg.Sender, &telebot.Photo{File: file}, nil)
	case fileExt == "mp3":
		bot.SendAudio(msg.Sender, &telebot.Audio{File: file}, nil)
	case fileExt == "mp4":
		bot.SendVideo(msg.Sender, &telebot.Video{Audio: telebot.Audio{File: file}}, nil)
	}
	return bot.SendDocument(msg.Sender, &telebot.Document{File: file}, nil)
}

func lsToMarkup(path string, page int) (telebot.ReplyMarkup, error) {
	const (
		maxItemsPerPage int = 30
		maxItemsPerRow  int = 3
	)

	reservedButtons := 1 // Back button by default
	files, err := ioutil.ReadDir(path)
	files = files[page*maxItemsPerPage:] // Apply paging
	if err != nil {
		return telebot.ReplyMarkup{}, err
	}
	if len(files) > maxItemsPerPage {
		reservedButtons++ // Reserve slot for next page button
		files = files[0:maxItemsPerPage]
	}

	markup := telebot.ReplyMarkup{
		InlineKeyboard: make([][]telebot.KeyboardButton, RoundNumber(float32(len(files))/float32(maxItemsPerRow))+reservedButtons),
	}
	// Add "Back" button at start
	markup.InlineKeyboard[0] = []telebot.KeyboardButton{
		telebot.KeyboardButton{
			Text: "..",
			Data: "/cd ..",
		},
	}

	for i := 1; i <= len(files); i += maxItemsPerRow {
		var row []telebot.KeyboardButton
		if len(files)-i > maxItemsPerRow {
			row = make([]telebot.KeyboardButton, maxItemsPerRow)
		} else {
			row = make([]telebot.KeyboardButton, len(files)-i)
		}

		for index, file := range files[i : i+len(row)] {
			row[index] = telebot.KeyboardButton{
				Text: file.Name(),
				Data: "/actions " + file.Name(),
			}
		}
		markup.InlineKeyboard[i/maxItemsPerRow+1] = row
	}

	if reservedButtons > 1 {
		markup.InlineKeyboard[len(markup.InlineKeyboard)-1] = []telebot.KeyboardButton{
			telebot.KeyboardButton{
				Text: "Next page",
				Data: "/ls " + strconv.Itoa(page+1),
			},
		}
	}

	return markup, nil
}
