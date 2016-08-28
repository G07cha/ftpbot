package main

import (
	"os/user"

	"github.com/tucnak/telebot"
)

// UserState is storing current state for each user that uses bot
type UserState struct {
	user           *telebot.User
	currentPath    string
	selectedAction userAction
}

type userAction uint8

// UserActions used for storing available action and preventing random integer insertion
var UserActions = struct {
	COPY, MOVE, RENAME userAction
}{1, 2, 3}

// UsersList global list of users with their current states
var UsersList []UserState

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

	newState := NewUser(u)
	UsersList = append(UsersList, newState)

	return &UsersList[len(UsersList)-1]
}