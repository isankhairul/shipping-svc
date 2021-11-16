package utils

import (
	"strconv"
	"unicode"
)

func PasswordValidation(s string) (eightOrMore, number, upper, lower, special bool) {
	letters := 0
	for _, c := range s {
		switch {
		case unicode.IsNumber(c):
			number = true
		case unicode.IsUpper(c):
			upper = true
			letters++
		case unicode.IsLower(c):
			lower = true
			letters++
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			special = true
		case unicode.IsLetter(c) || c == ' ':
			letters++
		default:
			return false, false, false, false, false
		}
	}
	eightOrMore = len(s) >= 8
	return
}

func CheckIsNumberic(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
