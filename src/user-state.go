package main

import (
	"os"
	"os/user"
	"path"

	"github.com/tucnak/telebot"
)

// UserState is storing current state for each user that uses bot
type UserState struct {
	user           *telebot.User
	currentPath    string
	selectedAction userAction
	selectedFile   string
}

type userAction uint8

// UserActions used for storing available action and preventing random integer insertion
var UserActions = struct {
	NONE, SELECT, RENAME userAction
}{0, 1, 2}

// UsersList global list of users with their current states
var UsersList []UserState

func init() {
	Router.routes["/cancel"] = ResetAction
}

// SelectFile appends provided filename to current directory
func (u *UserState) SelectFile(filename string) error {
	fullpath := path.Join(u.currentPath, filename)
	_, err := os.Stat(fullpath) // Check if item exists
	if err != nil {
		return err
	}

	u.selectedAction = UserActions.SELECT
	u.selectedFile = fullpath

	return nil
}

// NewUser create new user with default parameters
func NewUser(u *telebot.User) UserState {
	usr, err := user.Current() // Home directory
	if err != nil {
		usr = &user.User{HomeDir: "/"}
	}

	return UserState{
		user:        u,
		currentPath: usr.HomeDir,
	}
}

// GetCurrentState retrives state for user or creates new one
func GetCurrentState(u *telebot.User) *UserState {
	for index, state := range UsersList {
		if state.user.ID == u.ID {
			return &UsersList[index]
		}
	}

	// Create new user if no user found
	newState := NewUser(u)
	UsersList = append(UsersList, newState)

	return &UsersList[len(UsersList)-1]
}

// ResetAction set current action for user to NONE
func ResetAction(bot *telebot.Bot, msg telebot.Message) error {
	state := GetCurrentState(&msg.Sender)

	state.selectedFile = ""
	state.selectedAction = UserActions.NONE

	return bot.SendMessage(msg.Chat, "Got it, aborting!", nil)
}
