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

	c, i, analyzers := analyzer_setup_analyzers(scanner)
	formatter := MySQLFormatter{}
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

func analyzer_setlen(col *FormatterColumn) {
	if v, e := config.ColumnLens[col.Name]; e && v != nil {
		col.IsLenFixed = v.Max == v.Min
		col.MaxLen = v.Max
	}
}

func analyzer_setup_analyzers(scanner InputScanner) ([]FormatterColumn, []FormatterIndex, []Analyzer) {
	var output_columns []FormatterColumn
	var output_indexes []FormatterIndex
	var output_analyzers []Analyzer

	fields := scanner.Fields()

	for _, f := range fields {
		types, e := config.ColumnTypes[f]

		if !e || len(types) == 0 {
			Warn(fmt.Sprintf("no type was specified for `%s`, interpreted as nullable string.", f))
			output_columns = append(output_columns, FormatterColumn{
				Name: f, Type: FMT_TYPE_STR, Tags: []string{"nullable"}})

			analyzer_setlen(&output_columns[len(output_columns)-1])
		} else {
			fmt_type := StrType2FmtType(types[0])

			if fmt_type == FMT_TYPE_UNKNOWN {
				Warn(fmt.Sprintf("invalid type `%s`, `%s` interpreted as nullable string.", types[0], f))
				output_columns = append(output_columns, FormatterColumn{
					Name: f, Type: FMT_TYPE_STR, Tags: []string{"nullable"}})
				analyzer_setlen(&output_columns[len(output_columns)-1])
			} else {
				output_columns = append(output_columns, FormatterColumn{
					Name: f, Type: fmt_type, Tags: types[1:]})

				analyzer_setlen(&output_columns[len(output_columns)-1])

				for _, tag := range types[1:] {
					//log.Println(tag)
					if analyzer := GetAnalyzer(tag); analyzer != nil {
						c, i, e := analyzer.Init(output_columns[len(output_columns)-1])
						if e != nil {
							panic(e)
						}

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
