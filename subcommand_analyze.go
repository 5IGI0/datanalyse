package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
)

func analyzer() {
	scanner, err := InitScanner(flag.Arg(1))

	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to initialise scanner:", err.Error())
		os.Exit(1)
	}
	defer scanner.Close()

	formatter := MySQLFormatter{}
	c, i, analyzers, group_analyzers := analyzer_setup_analyzers(scanner, &formatter)
	output := flag.Arg(2)
	if output == "" {
		output = "-"
	}
	formatter.Init(flag.Arg(2), c, i, group_analyzers)
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

		for _, analyzer := range group_analyzers {
			if err := analyzer.Analyze(&data); err != nil {
				panic(err)
			}
		}

		formatter.WriteRow(data)
	}

}

func analyzer_setup_analyzers(scanner InputScanner, formatter Formatter) (map[string]FormatterColumn, []FormatterIndex, []Analyzer, []GroupAnalyzer) {
	var output_columns = make(map[string]FormatterColumn)
	var output_indexes []FormatterIndex
	var output_analyzers []Analyzer
	var output_group_analyzers []GroupAnalyzer

	fields := scanner.Fields()

	for _, f := range fields {
		colinf, e := config.ColumnInfos[f]

		if !e {
			colinf.Init()
		}

		column := FormatterColumn{
			Name: f, Tags: colinf.Tags,
			MaxLen: colinf.MaxLen, MinLen: colinf.MinLen,
			IsLenFixed: colinf.MaxLen == colinf.MinLen}

		if len(colinf.Type) == 0 {
			colinf.Type = "str"
			Warn(fmt.Sprintf("no type was specified for `%s`, interpreted as string.", f))
		}

		column.Type = StrType2FmtType(colinf.Type)
		if column.Type == FMT_TYPE_UNKNOWN {
			column.Type = FMT_TYPE_STR
			Warn(fmt.Sprintf("invalid type `%s`, `%s` interpreted as string.", colinf.Type, f))
		}

		if column.Tags == nil ||
			(!slices.Contains(column.Tags, "nonnull") && !slices.Contains(column.Tags, "nullable")) {
			Warn(fmt.Sprintf("`%s` has not been marked as nonnull or nullable, defaulting to nullable.", colinf.Type))
			column.Tags = append(column.Tags, "nullable")
		}

		if slices.Contains(column.Tags, "nonnull") && slices.Contains(column.Tags, "nullable") {
			Warn(fmt.Sprintf("`%s` has been marked as nonnull _AND_ nullable, defaulting to nullable.", colinf.Type))
			i := slices.Index(column.Tags, "nonnull")
			column.Tags = slices.Delete(column.Tags, i, i+1)
		}

		for _, tag := range colinf.Tags {
			if analyzer := GetAnalyzer(tag); analyzer != nil {
				/* check if it is already present (some tags can trigger the same analyzer) */
				a_type := analyzer.GetAnalyzerType()
				for _, a := range column.Analyzers {
					if a.GetAnalyzerType() == a_type {
						continue
					}
				}

				/* add analyzer */
				c, i, e := analyzer.Init(column, formatter)
				AssertError(e)

				column.Analyzers = append(column.Analyzers, analyzer)
				for _, col := range c {
					output_columns[col.Name] = col
				}
				output_indexes = append(output_indexes, i...)
				output_analyzers = append(output_analyzers, analyzer)
			}
		}

		output_columns[column.Name] = column
	}

	/* initialize group analyzer */
	for _, gi := range config.GroupInfos {
		ga := GetGroupAnalyzer(gi.Kind)
		if ga == nil {
			Warn(fmt.Sprintf("Unknown `%s` group kind, ignoring it.", gi.Kind))
			continue
		}

		cols := []FormatterColumn{}

		s := true
		for _, colname := range gi.Fields {
			col := output_columns[colname]
			if col.Name == "" {
				Warn(fmt.Sprintf("Unknown column `%s` for `%s` group, abandoning it.", colname, gi.Kind))
				s = false
				break
			}
			cols = append(cols, col)
		}
		if !s {
			continue
		}

		columns, i, e := ga.Init(cols, formatter)
		AssertError(e)

		for _, col := range columns {
			output_columns[col.Name] = col
		}
		output_indexes = append(output_indexes, i...)
		output_group_analyzers = append(output_group_analyzers, ga)
	}

	return output_columns, output_indexes, output_analyzers, output_group_analyzers
}
