package main

import (
	"flag"
	"fmt"
	"os"
)

func analyzer() {
	scanner, err := InitScanner(flag.Arg(1))

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to initialise scanner:", err.Error())
		os.Exit(1)
	}
	defer scanner.Close()

	formatter := MySQLFormatter{}
	c, i, analyzers := analyzer_setup_analyzers(scanner, &formatter)
	output := flag.Arg(2)
	if output == "" {
		output = "-"
	}
	formatter.Init(flag.Arg(2), c, i)
	defer formatter.Close()

	for {
		data, err := scanner.ReadRow()
		if err != nil {
			panic(err)
		}
		if data == nil {
			break
		}

		for _, analyzer := range analyzers {
			if err := analyzer.Analyze(&data); err != nil {
				panic(err)
			}
		}

		formatter.WriteRow(data)
	}

}

func analyzer_setup_analyzers(scanner InputScanner, formatter Formatter) ([]FormatterColumn, []FormatterIndex, []Analyzer) {
	var output_columns []FormatterColumn
	var output_indexes []FormatterIndex
	var output_analyzers []Analyzer

	fields := scanner.Fields()

	for _, f := range fields {
		colinf, e := config.ColumnInfos[f]

		if !e || len(colinf.Type) == 0 {
			Warn(fmt.Sprintf("no type was specified for `%s`, interpreted as nullable string.", f))
			output_columns = append(output_columns, FormatterColumn{
				Name: f, Type: FMT_TYPE_STR, Tags: []string{"nullable"},
				MaxLen: colinf.MaxLen, MinLen: colinf.MinLen, IsLenFixed: colinf.MaxLen == colinf.MinLen})
		} else {
			fmt_type := StrType2FmtType(colinf.Type)

			if fmt_type == FMT_TYPE_UNKNOWN {
				Warn(fmt.Sprintf("invalid type `%s`, `%s` interpreted as nullable string.", colinf.Type, f))
				output_columns = append(output_columns, FormatterColumn{
					Name: f, Type: FMT_TYPE_STR, Tags: []string{"nullable"},
					MaxLen: colinf.MaxLen, MinLen: colinf.MinLen, IsLenFixed: colinf.MaxLen == colinf.MinLen})
			} else {
				output_columns = append(output_columns, FormatterColumn{
					Name: f, Type: fmt_type, Tags: colinf.Tags,
					MaxLen: colinf.MaxLen, MinLen: colinf.MinLen, IsLenFixed: colinf.MaxLen == colinf.MinLen})

				for _, tag := range colinf.Tags {
					//log.Println(tag)
					if analyzer := GetAnalyzer(tag); analyzer != nil {
						c, i, e := analyzer.Init(output_columns[len(output_columns)-1], formatter)
						if e != nil {
							panic(e)
						}

						output_columns[len(output_columns)-1].Analyzers = append(output_columns[len(output_columns)-1].Analyzers, analyzer)
						output_columns = append(output_columns, c...)
						output_indexes = append(output_indexes, i...)
						output_analyzers = append(output_analyzers, analyzer)
					}
				}
			}
		}

	}

	return output_columns, output_indexes, output_analyzers
}
