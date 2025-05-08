package main

import "strconv"

func IsNumber(number string) bool {
	if _, err := strconv.Atoi(number); err == nil {
		return true
	}
	return false
}
