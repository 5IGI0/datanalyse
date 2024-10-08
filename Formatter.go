package main

type Formatter interface {
	Init(output_file string, columns map[string]FormatterColumn, indexes []FormatterIndex, group_analyzers []GroupAnalyzer) error
	GetFeatures() int
	WriteRow(map[string]*string) error
	Close()
}

const (
	FMT_FEATURE_GENERATED_AS = 1
)

type GeneratorData struct {
	Format       string   `json:"format"`
	PrimaryType  string   `json:"primary_type"`
	Tags         []string `json:"tags"`
	Version      int      `json:"ver"`
	LinkedColumn string   `json:"linked_col"`
	CustomData   any      `json:"custom_data"`
}

type FormatterColumn struct {
	Name        string
	ForceString bool
	Tags        []string
	IsInvisible bool
	// expression used to generate the targeted value
	AlwaysGeneratedAs string
	Generator         interface {
		GetGeneratorInfo() GeneratorInfo
	}
	// data of the analyzer that generated this column
	GeneratorData *GeneratorData
	// analyzers that will pass on it to generate new columns
	Analyzers []Analyzer
}

const (
	FMT_TYPE_UNKNOWN = -1
	FMT_TYPE_STR     = 0
	FMT_TYPE_INT8    = 1
	FMT_TYPE_UINT8   = 2
	FMT_TYPE_INT16   = 3
	FMT_TYPE_UINT16  = 4
	FMT_TYPE_INT32   = 5
	FMT_TYPE_UINT32  = 6
	FMT_TYPE_INT64   = 7
	FMT_TYPE_UINT64  = 8
	FMT_TYPE_ENUM    = 9
)

type FormatterIndex struct {
	ColumnName string
	IndexName  string
	Reversed   bool
}

func StrType2FmtType(typ string) int8 {
	switch typ {
	case "str":
		return FMT_TYPE_STR
	case "int8":
		return FMT_TYPE_INT8
	case "int16":
		return FMT_TYPE_INT16
	case "int32":
		return FMT_TYPE_INT32
	case "int64":
		return FMT_TYPE_INT64
	case "uint16":
		return FMT_TYPE_UINT16
	case "uint32":
		return FMT_TYPE_UINT32
	}

	return FMT_TYPE_UNKNOWN
}
