package main

import (
	"flag"
	"fmt"
)

func main() {
	ParseConfig()

	if flag.Arg(0) == "hint" {
		hinter()
	} else {
		fmt.Println("unknown subcommand")
	}
}
