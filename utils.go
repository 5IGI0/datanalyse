package main

import (
	"fmt"
	"os"
)

func Warn(warn string) {
	fmt.Fprintln(os.Stderr, "Warning:", warn)
}
