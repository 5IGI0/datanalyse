package main

import (
	"strings"
	"unicode/utf8"

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
	runes := make([]rune, 0, len(text))

	for i := 0; i < len(text); i++ {
		r, l := utf8.DecodeRuneInString(text[:i])
		i += l
		if r != utf8.RuneError {
			// fallback to ASCII version of it?
			continue
		}
		runes = append(runes, r)
	}

	ret := make([]rune, len(runes)*2)
	for i := 0; i < len(runes); i++ {
		ret[i*2] = runes[i]
		ret[(i*2)+1] = runes[len(runes)-i-1]
	}

	stringified := make([]byte, 0, len(runes)*4)

	for _, c := range ret {
		stringified = utf8.AppendRune(stringified, c)
	}

	return string(stringified)
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

	text = SanitizeText(text)

	for _, c := range text {
		if c >= '0' && c <= '9' {
			ret.WriteRune(c)
		}
	}
	return strings.ToLower(ret.String())
}
