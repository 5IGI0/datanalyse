package main

import "fmt"

type InputScanner interface {
	Init(input_file string) error
	ReadRow() (map[string]*string, error)
	// NOTE: they may not be present
	Fields() []string
	Close()
}

func InitScanner(input_file string) (InputScanner, error) {
	if config.Scanner == "csv" {
		var scanner CsvScanner
		if err := scanner.Init(input_file); err != nil {
			return nil, err
		}
		return &scanner, nil
	} else {
		return nil, fmt.Errorf("unknown scanner `%s`", config.Scanner)
	}
}
