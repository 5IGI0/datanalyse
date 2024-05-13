package main

import (
	"flag"
	"fmt"
)

func main() {
	ParseConfig()

	if flag.Arg(0) == "hint" {
		hinter()
	} else if flag.Arg(0) == "analyze" {
		analyzer()
	} else {
		fmt.Println("unknown subcommand")
	}
}
