package main

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/tucnak/telebot"
)

func init() {
	Router.routes["/ls"] = ls
	Router.routes["/actions"] = showActions
	Router.routes["/cd"] = cd
	Router.routes["/download"] = download
	Router.routes["/rm"] = Confirm("delete")
	Router.routes["/delete"] = rm
	Router.routes["/rename"] = rename
	Router.routes["/select"] = selectItem
	Router.routes["/move"] = moveItem
	Router.routes["/copy"] = copyItem
	Router.routes["text"] = handleText
}

func handleText(bot *telebot.Bot, msg telebot.Message) error {
	state := GetCurrentState(&msg.Sender)

	if state.selectedAction == UserActions.RENAME {
		state.selectedAction = UserActions.NONE
		dir, _ := path.Split(state.selectedFile)
		newName := path.Join(dir, msg.Text)

		err := os.Rename(state.selectedFile, newName)
		if err != nil {
			return bot.SendMessage(msg.Chat, "Failed to rename "+state.selectedFile, nil)
		}

		return bot.SendMessage(msg.Chat, "File "+state.selectedFile+" renamed successfully", nil)
	}

	return nil
}

func ls(bot *telebot.Bot, msg telebot.Message) error {
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

func selectItem(bot *telebot.Bot, msg telebot.Message) error {
	filename := GetRemainingText(msg.Text)
	err := GetCurrentState(&msg.Sender).SelectFile(filename)

	if err != nil {
		return err
	}

	return bot.SendMessage(msg.Chat, "File selected, navigate to desired folder and click Copy or Move button", &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.KeyboardButton{
				[]telebot.KeyboardButton{
					telebot.KeyboardButton{
						Text: "Cancel",
						Data: "/cancel",
					},
				},
			},
		},
	})
}

func showActions(bot *telebot.Bot, msg telebot.Message) error {
	filename := GetRemainingText(msg.Text)
	state := GetCurrentState(&msg.Sender)
	file, err := os.Open(path.Join(state.currentPath, filename))
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

	if state.selectedAction == UserActions.SELECT {
		_, filename := path.Split(state.selectedFile)
		folderActions = append(folderActions, telebot.KeyboardButton{
			Text: "Copy " + filename,
			Data: "/copy " + fileInfo.Name(),
		}, telebot.KeyboardButton{
			Text: "Move " + filename,
			Data: "/move " + fileInfo.Name(),
		})
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
					Data: "/rm " + fileInfo.Name(),
				},
			},
			[]telebot.KeyboardButton{
				telebot.KeyboardButton{
					Text: "Select",
					Data: "/select " + fileInfo.Name(),
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

func moveItem(bot *telebot.Bot, msg telebot.Message) error {
	state := GetCurrentState(&msg.Sender)
	foldername := GetRemainingText(msg.Text)
	_, itemname := path.Split(state.selectedFile)
	newPath := path.Join(state.currentPath, foldername, itemname)

	err := os.Rename(state.selectedFile, newPath)
	if err != nil {
		return err
	}

	state.selectedAction = UserActions.NONE
	state.selectedFile = ""

	return bot.SendMessage(msg.Chat, itemname+" moved successfully", nil)
}

func copyItem(bot *telebot.Bot, msg telebot.Message) error {
	state := GetCurrentState(&msg.Sender)
	foldername := GetRemainingText(msg.Text)
	_, itemname := path.Split(state.selectedFile)
	newPath := path.Join(state.currentPath, foldername, itemname)

	info, err := os.Stat(path.Join(state.selectedFile))
	if err != nil {
		return err
	}

	if info.IsDir() == true {
		err = CopyDir(state.selectedFile, newPath)
	} else {
		err = CopyFile(state.selectedFile, newPath)
	}

	if err != nil {
		return err
	}

	state.selectedAction = UserActions.NONE
	state.selectedFile = ""

	return bot.SendMessage(msg.Chat, itemname+" copied successfully", nil)
}

func cd(bot *telebot.Bot, msg telebot.Message) error {
	state := GetCurrentState(&msg.Sender)

	newPath := path.Join(state.currentPath, GetRemainingText(msg.Text))
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
func download(bot *telebot.Bot, msg telebot.Message) error {
	filename := GetRemainingText(msg.Text)
	fileExt := path.Ext(filename)

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

func rename(bot *telebot.Bot, msg telebot.Message) error {
	state := GetCurrentState(&msg.Sender)
	state.selectedAction = UserActions.RENAME
	state.selectedFile = path.Join(state.currentPath, GetRemainingText(msg.Text))

	return bot.SendMessage(msg.Chat, "Please send me a new name", &telebot.SendOptions{
		ReplyMarkup: telebot.ReplyMarkup{
			InlineKeyboard: [][]telebot.KeyboardButton{
				[]telebot.KeyboardButton{
					telebot.KeyboardButton{
						Text: "Cancel",
						Data: "/cancel",
					},
				},
			},
		},
	})
}

// Remove file or folder from filesystem
func rm(bot *telebot.Bot, msg telebot.Message) error {
	filename := GetRemainingText(msg.Text)
	fullpath := path.Join(GetCurrentState(&msg.Sender).currentPath, filename)

	err := os.RemoveAll(fullpath)
	if err != nil {
		return bot.SendMessage(msg.Chat, "Failed to remove "+filename, nil)
	}

	return bot.SendMessage(msg.Chat, filename+" removed successfully", nil)
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
