package main

type Formatter interface {
	Init(output_file string, columns []FormatterColumn, indexes []FormatterIndex) error
	WriteRow(map[string]*string) error
	Close()
}

type FormatterColumn struct {
	Name       string
	Type       int8
	Tags       []string
	MaxLen     int
	MinLen     int
	IsLenFixed bool
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
