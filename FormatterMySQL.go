package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strings"
)

type MySQLFormatter struct {
	OutputFile     *os.File
	Writer         *bufio.Writer
	Columns        []FormatterColumn
	Indexes        []FormatterIndex
	GroupAnalyzers []GroupAnalyzer
	reverse_idx    ReverseIndexEmulator
	CachedQuery    strings.Builder
	InternalId     uint32
}

func (f *MySQLFormatter) GetFeatures() int {
	return FMT_FEATURE_GENERATED_AS
}

func (f *MySQLFormatter) Init(output_file string, columns_map map[string]FormatterColumn, indexes []FormatterIndex, group_analyzers []GroupAnalyzer) error {
	if output_file == "-" {
		f.OutputFile = os.Stdout
	} else {
		var err error
		f.OutputFile, err = os.OpenFile(output_file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
		if err != nil {
			panic(err)
		}
	}

	f.Writer = bufio.NewWriter(f.OutputFile)
	f.GroupAnalyzers = group_analyzers

	columns := []FormatterColumn{}
	for _, col := range columns_map {
		columns = append(columns, col)
	}
	columns = append(columns, f.reverse_idx.Init(&indexes, f)...)

	f.Columns = append([]FormatterColumn{{Name: "__internal_id", Type: FMT_TYPE_UINT32, Tags: []string{"nonnull"}}}, columns...)
	f.Indexes = indexes

	f.Writer.WriteString("SET SESSION sql_mode='';\nDROP TABLE IF EXISTS `")
	f.Writer.WriteString(config.FmtSQLTable)
	f.Writer.WriteString("`;\nCREATE TABLE `")
	f.Writer.WriteString(config.FmtSQLTable)
	f.Writer.WriteString("` (\n")

	for i, col := range f.Columns {
		if i != 0 {
			f.Writer.WriteString(",\n")
		}
		f.Writer.WriteString("\t`")
		f.Writer.WriteString(col.Name)
		f.Writer.WriteString("` ")
		f.Writer.WriteString(map[int8]string{
			FMT_TYPE_STR:    f.get_string_type(col),
			FMT_TYPE_INT8:   "TINYINT",
			FMT_TYPE_UINT8:  "TINYINT UNSIGNED",
			FMT_TYPE_INT16:  "SMALLINT",
			FMT_TYPE_UINT16: "SMALLINT UNSIGNED",
			FMT_TYPE_INT32:  "INT",
			FMT_TYPE_UINT32: "INT UNSIGNED",
			FMT_TYPE_INT64:  "BIG",
			FMT_TYPE_UINT64: "BIG UNSIGNED",
		}[col.Type])
		f.Writer.WriteString(" ")

		/* always generated as */

		if col.AlwaysGeneratedAs == "" {
			/* cannot use NOT NULL on generated values (mariadb) */
			if slices.Index(col.Tags, "nonnull") != -1 {
				f.Writer.WriteString(" NOT NULL")
			}
		} else {
			f.Writer.WriteString(" GENERATED ALWAYS AS (")
			f.Writer.WriteString(col.AlwaysGeneratedAs)

			/* if it is invisible, then we don't store it */
			if col.IsInvisible {
				f.Writer.WriteString(") VIRTUAL")
			} else {
				f.Writer.WriteString(") STORED")
			}
		}

		if col.IsInvisible {
			f.Writer.WriteString(" INVISIBLE")
		}

		f.Writer.WriteString(" COMMENT ")
		f.Writer.WriteString(EscapeString(f.generate_column_comment(&col)))
	}
	f.Writer.WriteString(",\n\tPRIMARY KEY(`__internal_id`)) COMMENT ")
	f.Writer.WriteString(EscapeString(f.generate_table_comment()) + ";\n")
	return nil
}

func (f *MySQLFormatter) get_string_type(column FormatterColumn) string {
	if column.IsLenFixed {
		return fmt.Sprint("CHAR(", column.MaxLen, ")")
	} else if column.MaxLen != 0 && column.MaxLen < 1024 {
		return fmt.Sprint("VARCHAR(", column.MaxLen, ")")
	}
	return "TEXT"
}

func (f *MySQLFormatter) _startInsertQuery() {
	f.CachedQuery.Reset()

	f.CachedQuery.WriteString("INSERT INTO `")
	f.CachedQuery.WriteString(config.FmtSQLTable)
	f.CachedQuery.WriteString("`(`")
	i := 0
	for _, col := range f.Columns {
		if col.AlwaysGeneratedAs != "" {
			continue
		}
		if i != 0 {
			f.CachedQuery.WriteString("`,`")
		}
		f.CachedQuery.WriteString(col.Name)
		i++
	}
	f.CachedQuery.WriteString("`) VALUES\n")
}

func (f *MySQLFormatter) _encodeRow(row map[string]*string) string {
	var output string = "("

	{
		tmp := fmt.Sprint(f.InternalId)
		row["__internal_id"] = &tmp
		f.InternalId++
	}

	i := 0
	for _, col := range f.Columns {
		if col.AlwaysGeneratedAs != "" {
			continue
		}
		tmp := row[col.Name]

		if i != 0 {
			output += ","
		}

		if tmp == nil {
			output += "NULL"
		} else {
			output += EscapeString(*tmp)
		}

		i++
	}

	output += ")\n"
	return output
}

func (f *MySQLFormatter) WriteRow(row map[string]*string) error {
	f.reverse_idx.Apply(&row)
	if f.CachedQuery.Len() == 0 {
		f._startInsertQuery()
		f.CachedQuery.WriteString(f._encodeRow(row))
		return nil
	}

	encoded := f._encodeRow(row)

	if len(encoded)+f.CachedQuery.Len() > int(config.FmtSQLMaxQuerySize)-3 {
		f.CachedQuery.WriteString(";\n")
		f.Writer.WriteString(f.CachedQuery.String())
		f._startInsertQuery()
	} else {
		f.CachedQuery.WriteByte(',')
	}

	f.CachedQuery.WriteString(encoded)

	return nil
}

func (f *MySQLFormatter) Close() {
	if f.CachedQuery.Len() != 0 {
		f.CachedQuery.WriteString(";\n")
		f.Writer.WriteString(f.CachedQuery.String())
	}
	for _, index := range f.Indexes {
		fmt.Fprint(f.Writer, "CREATE INDEX `", index.IndexName, "` ON `", config.FmtSQLTable, "`(`", index.ColumnName, "`);\n")
	}
	f.Writer.Flush()
	f.OutputFile.Close()
}

func (f *MySQLFormatter) generate_column_comment(column *FormatterColumn) string {
	type AnalyzerManifest struct {
		Analyzer GeneratorInfo `json:"analyzer"`
		Data     any           `json:"data"`
	}

	var tmp struct {
		Generator     *GeneratorInfo `json:"generator"`
		GeneratorData any            `json:"generator_data"`
		IsInvisible   bool           `json:"is_invisible"`
		Tags          []string       `json:"tags"`
		Analyzers     map[string]AnalyzerManifest
	}
	tmp.Analyzers = make(map[string]AnalyzerManifest)

	if column.Generator != nil {
		inf := column.Generator.GetGeneratorInfo()
		tmp.Generator = &inf
	}
	tmp.GeneratorData = column.GeneratorData
	tmp.IsInvisible = column.IsInvisible
	tmp.Tags = column.Tags

	for _, analyzer := range column.Analyzers {
		info := analyzer.GetGeneratorInfo()
		tmp.Analyzers[info.Name] = AnalyzerManifest{Analyzer: info, Data: analyzer.GetAnalyzerData()}
	}

	b, _ := json.Marshal(tmp)
	return string(b)
}

func (f *MySQLFormatter) generate_table_comment() string {
	type AnalyzerManifest struct {
		Analyzer GeneratorInfo `json:"analyzer"`
		Data     any           `json:"data"`
	}

	var tmp struct {
		Name          string              `json:"name"`
		Description   string              `json:"description"`
		Version       uint32              `json:"version"`
		Meta          map[string]string   `json:"meta"`
		GroupAnalyzer []AnalyzerManifest  `json:"group_analyzers"`
		VirtColMap    map[string][]string `json:"virtcol_map"`
	}

	tmp.GroupAnalyzer = []AnalyzerManifest{}
	tmp.VirtColMap = make(map[string][]string)

	for _, a := range f.GroupAnalyzers {
		tmp.GroupAnalyzer = append(tmp.GroupAnalyzer, AnalyzerManifest{
			Analyzer: a.GetGeneratorInfo(),
			Data:     a.GetAnalyzerData(),
		})

		m := a.GetVirtualColumnMap()

		for v, c := range m {
			tmp.VirtColMap[v] = c
		}
	}

	if config.DatasetName != "" {
		tmp.Name = config.DatasetName
	} else {
		tmp.Name = config.FmtSQLTable
	}
	tmp.Description = config.DatasetDescription
	tmp.Version = 1
	tmp.Meta = config.DatasetMeta

	b, _ := json.Marshal(tmp)
	return string(b)
}
