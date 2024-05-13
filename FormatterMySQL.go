package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

type MySQLFormatter struct {
	OutputFile  *os.File
	Writer      *bufio.Writer
	Columns     []FormatterColumn
	Indexes     []FormatterIndex
	reverse_idx ReverseIndexEmulator
	CachedQuery strings.Builder
}

func (f *MySQLFormatter) Init(output_file string, columns []FormatterColumn, indexes []FormatterIndex) error {
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

	columns = append(columns, f.reverse_idx.Init(&indexes)...)

	f.Columns = columns
	f.Indexes = indexes

	f.Writer.WriteString("CREATE TABLE `")
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
		if slices.Index(col.Tags, "nonnull") != -1 {
			f.Writer.WriteString(" NOT NULL")
		}
	}
	f.Writer.WriteString(");\n")
	return nil
}

func (f *MySQLFormatter) get_string_type(column FormatterColumn) string {
	if column.IsLenFixed {
		return fmt.Sprint("CHAR(", column.IsLenFixed, ")")
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
	for i, col := range f.Columns {
		if i != 0 {
			f.CachedQuery.WriteString("`,`")

		}
		f.CachedQuery.WriteString(col.Name)
	}
	f.CachedQuery.WriteString("`) VALUES\n")
}

func (f *MySQLFormatter) _encodeRow(row map[string]*string) string {
	var output string = "("

	for i, col := range f.Columns {
		tmp := row[col.Name]

		if i != 0 {
			output += ","
		}

		if tmp == nil {
			output += "NULL"
		} else {
			output += EscapeString(*tmp)
		}
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
