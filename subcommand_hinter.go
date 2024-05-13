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
		fmt.Print("  --column-type=", k, ":", h.GetType(), ":", strings.Join(h.get_tags(), ":"))

		if h.GetType() == "str" {
			fmt.Print(" \\\n  --column-len=", k, ":", h.GetMaxLen(), ":", h.GetMinLen())
		}
	}
	fmt.Print("\n")
}
