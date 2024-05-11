package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
)

type CsvScanner struct {
	input_file   *os.File
	reader       *csv.Reader
	fieldNames   []string
	strip        bool
	nullValues   []string
	rowCount     uint64
	stripby      string
	publicFields []string
}

func (scanner *CsvScanner) Init(input_file string) error {
	if input_file == "-" {
		scanner.input_file = os.Stdin
	} else {
		var err error
		if scanner.input_file, err = os.Open(input_file); err != nil {
			return err
		}
	}

	scanner.stripby = config.ScanCsvStripBy
	scanner.reader = csv.NewReader(scanner.input_file)
	scanner.reader.FieldsPerRecord = -1
	{
		tmp := []rune(config.ScanCsvDelimiter)
		if len(tmp) != 1 {
			scanner.input_file.Close()
			return errors.New("delimiter must be one character")
		}
		scanner.reader.Comma = tmp[0]
	}

	if config.ScanCsvFieldNames != "" {
		scanner.fieldNames = strings.Split(config.ScanCsvFieldNames, ",")
	} else {
		rows, err := scanner.reader.Read()
		if err != nil {
			scanner.input_file.Close()
			return err
		}
		scanner.fieldNames = rows
	}

	if config.ScanCsvAsNull != "" {
		scanner.nullValues = strings.Split(config.ScanCsvAsNull, ",")
	}

	return nil
}

func (scanner *CsvScanner) ReadRow() (map[string]*string, error) {
	for {
		row, err := scanner.reader.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil, nil
			}
			return nil, err
		}
		scanner.rowCount += 1
		if len(row) != len(scanner.fieldNames) {
			Warn(fmt.Sprintf("Row #%d doesn't have as field as expected (%d instead of %d)",
				scanner.rowCount-1, len(row), len(scanner.fieldNames)))
			continue
		}
		ret := make(map[string]*string)

		for i := 0; i < len(row); i++ {
			if scanner.fieldNames[i] == "_" {
				continue
			}

			tmp_val := row[i]
			if scanner.strip {
				tmp_val = strings.Trim(tmp_val, scanner.stripby)
			}

			if slices.Index(scanner.nullValues[:], tmp_val) == -1 {
				ret[scanner.fieldNames[i]] = &tmp_val
			}
		}
		return ret, nil
	}
}

func (scanner *CsvScanner) Fields() []string {
	if scanner.publicFields != nil {
		return scanner.publicFields
	}
	for _, v := range scanner.fieldNames {
		if v != "_" {
			scanner.publicFields = append(scanner.publicFields, v)
		}
	}
	return scanner.publicFields
}

func (scanner *CsvScanner) Close() {
	scanner.input_file.Close()
}
