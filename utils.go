package main

import (
	"fmt"
	"os"
	"unicode/utf8"
)

func Warn(warn string) {
	fmt.Fprintln(os.Stderr, "Warning:", warn)
}

func reverse_str(s string) string {
	reversed := make([]byte, len(s))
	i := 0

	for len(s) > 0 {
		r, size := utf8.DecodeLastRuneInString(s)
		s = s[:len(s)-size]
		i += utf8.EncodeRune(reversed[i:], r)
	}

	return string(reversed)
}

func Assert(b bool) {
	if !b {
		panic("Assert failed")
	}
}

func AssertError(err error) {
	if err != nil {
		panic(err)
	}
}
