package main

import (
	"strings"

	gotextsanitizer "github.com/5IGI0/go-text-sanitizer"
)

func SanitizeText(input string) string {
	san, err := gotextsanitizer.Unidecode(input)
	AssertError(err)
	return san
}

func BidirectionalizeTextA(text string) string {
	ret := make([]byte, len(text)*2)

	for i := 0; i < len(text); i++ {
		ret[i*2] = text[i]
		ret[(i*2)+1] = text[len(text)-i-1]
	}

	return string(ret)
}

func BidirectionalizeTextW(text string) string {
	ret := make([]rune, len(text)*2)
	runes := []rune(text)

	for i := 0; i < len(runes); i++ {
		ret[i*2] = runes[i]
		ret[(i*2)+1] = runes[len(runes)-i-1]
	}

	return string(ret)
}

func OnlyAlphaNum(text string) string {
	var ret strings.Builder

	text = SanitizeText(text)

	for _, c := range text {
		if (c >= 'a' && c <= 'z') ||
			(c >= 'A' && c <= 'Z') ||
			(c >= '0' && c <= '9') {
			ret.WriteRune(c)
		}
	}
	return strings.ToLower(ret.String())
}

func OnlyNum(text string) string {
	var ret strings.Builder

	for _, c := range text {
		if c >= '0' && c <= '9' {
			ret.WriteRune(c)
		}
	}
	return ret.String()
}
