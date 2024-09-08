package main

import (
	"flag"
	"fmt"
	"os"
)

type ColumnInfo struct {
	Type   string
	Tags   []string
	MaxLen int
	MinLen int
}

func (c *ColumnInfo) Init() {
	c.MaxLen = -1
	c.MinLen = -1
}

type Config struct {
	// scanner config
	Scanner   string
	InputFile string

	// CSV Scanner config
	ScanCsvDelimiter   string
	ScanCsvQuote       string
	ScanCsvAsNull      []string
	ScanCsvNoNull      bool
	ScanCsvStrip       bool
	ScanCsvColumnNames string
	ScanCsvStripBy     string

	// Formatters config
	Formater           string
	DatasetName        string
	DatasetDescription string
	DatasetMeta        map[string]string

	// SQL Formatters config
	FmtSQLTable        string
	FmtSQLMaxQuerySize int64

	// Analyzer config
	ColumnInfos map[string]*ColumnInfo
}

var config = Config{}

func ParseConfig() {
	config.ColumnInfos = make(map[string]*ColumnInfo)
	config.DatasetMeta = make(map[string]string)

	// scanner config
	flag.StringVar(&config.Scanner, "scanner", "csv", "What scanner to use (default: csv)")

	// CSV config
	flag.StringVar(&config.ScanCsvDelimiter, "csv-delimiter", ",", "CSV delimiter (commas by default)")
	// TODO: currently disabled since the golang csv reader doesn't support it (i guess i'm going to recode it)
	//flag.StringVar(&config.ScanCsvSeparator, "csv-quote", "\"", "CSV quotes, empty to disable it (double quotes bu default)")
	flag.BoolVar(&config.ScanCsvStrip, "csv-strip", true, "Strip columns (default: true)")
	flag.BoolVar(&config.ScanCsvNoNull, "csv-no-null", false, "Disable nullation of csv columns")
	flag.Var(MultipleVars{&config.ScanCsvAsNull}, "csv-null", "Values that should be parsed as null (default: empty values, overwritten by this or --csv-no-null)")
	flag.StringVar(&config.ScanCsvColumnNames, "csv-column-names", "", "Name of columns in the CSV file, columns `_` are ignored (by default the scanner reads the first row)")
	flag.StringVar(&config.ScanCsvStripBy, "csv-strip-by", " \n\f\r\t\v", "Character used to strip values (default: ASCII spaces (isspace + C_LOCAL))")

	// Formatters config
	flag.StringVar(&config.Formater, "formatter", "MySQL", "What formatter to use (default: MySQL)")
	flag.StringVar(&config.DatasetName, "dataset-name", "", "Specify dataset's name")
	flag.StringVar(&config.DatasetDescription, "dataset-description", "", "Specify dataset's description")
	flag.Var(DatasetMetaVar{&config.DatasetMeta}, "dataset-meta", "Specify dataset meta (format: <key>:<value>)")

	// SQL Formatters config
	flag.StringVar(&config.FmtSQLTable, "sql-table", "", "Output table name")
	flag.Int64Var(&config.FmtSQLMaxQuerySize, "sql-max-query-size", 1048576, "Max SQL query size")

	// Analyzer config
	flag.Var(&ColumnTypeVar{&config.ColumnInfos}, "column-type", "Specify column's type, format: <column name>:<type>")
	flag.Var(&ColumnLenVar{&config.ColumnInfos}, "column-len", "Specify column's len, format: <column name>:<max len>[:<min len>]")
	flag.Var(&ColumnTagsVar{&config.ColumnInfos}, "column-tags", "Specify column's tags, format: <column name>:<tags>:...")

	flag.Parse()

	if len(flag.Args()) == 0 {
		fmt.Println("Usage:", os.Args[0], "[FLAGS...] <subcommand>")
		fmt.Println("subcommands:")
		fmt.Println("  - hint:    analyze columns and try to guess columns' types (takes input file)")
		fmt.Println("  - analyze: analyze columns and use `formater` to store it (takes input and output file)")
		os.Exit(0)
	}
}
