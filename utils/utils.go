package utils

import "strings"

func TrimString(text string) string {
	text = strings.ReplaceAll(text, "  ", "")
	text = strings.ReplaceAll(text, "\r", "")
	return strings.ReplaceAll(text, "\n", "")
}

var FaToEn = strings.NewReplacer(
	"۰", "0",
	"۱", "1",
	"۲", "2",
	"۳", "3",
	"۴", "4",
	"۵", "5",
	"۶", "6",
	"۷", "7",
	"۸", "8",
	"۹", "9",
)
