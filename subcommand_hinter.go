package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func hinter() {
	scanner, err := InitScanner(flag.Arg(1))

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to initialise scanner:", err.Error())
		os.Exit(1)
	}
	defer scanner.Close()

	hinters := make(map[string]*Hinter)
	for _, k := range scanner.Fields() {
		hinters[k] = &Hinter{}
		hinters[k].Init()
	}

	for {
		data, err := scanner.ReadRow()
		if err != nil {
			panic(err)
		}
		if data == nil {
			break
		}

		for k, h := range hinters {
			h.Analyze(data[k])
		}
	}

	fmt.Println("recommended flags:")
	a := false
	for k, h := range hinters {
		if a {
			fmt.Print(" \\\n")
		} else {
			a = true
		}
		fmt.Print("  --column-type=", k, ":", h.GetType())
		fmt.Print(" \\\n  --column-tags=", k, ":", strings.Join(h.GetTags(), ":"))

		if h.GetType() == "str" {
			fmt.Print(" \\\n  --column-len=", k, ":", h.GetMaxLen(), ":", h.GetMinLen())
		}
	}

	// stuff that need to be "manually checked"
	a = false
	for k, h := range hinters {
		tags := h.GuessTags(k)

		if len(tags) != 0 {
			if !a {
				a = true
				fmt.Print(" \\\n ")
			}
			fmt.Print(" \\\n  --column-tags=", k, ":", strings.Join(tags, ":"))
		}
	}
	fmt.Print("\n")
}
