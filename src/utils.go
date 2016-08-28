package main

// RoundNumber does proper rounding for decimals(in case if .1 and more will add one)
func RoundNumber(num float32) int {
	if num > float32(int(num)) {
		return int(num) + 1
	}
	return int(num)
}
