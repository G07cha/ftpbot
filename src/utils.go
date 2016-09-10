package main

import "strings"

// RoundNumber does proper rounding for decimals(in case if .1 and more will add one)
func RoundNumber(num float32) int {
	if num > float32(int(num)) {
		return int(num) + 1
	}
	return int(num)
}

// GetRemainingText removes command from given message and return remaining part of the message
// For example: "/getme foo bar" -> "foo bar"
func GetRemainingText(text string) string {
	return text[strings.Index(text, " ")+1:]
}
