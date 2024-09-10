package main

import "strings"

type ReverseIndexEmulator struct {
	Indexes        []FormatterIndex
	HasGeneratedAs bool
}

type ReverseIndexEmulatorIdxInfo struct {
	LinkedColumn string `json:"linked_column"`
}

func (e *ReverseIndexEmulator) Init(indexes *[]FormatterIndex, f Formatter) []FormatterColumn {
	var ret []FormatterColumn

	e.HasGeneratedAs = (f.GetFeatures() & FMT_FEATURE_GENERATED_AS) != 0

	for i, index := range *indexes {
		if index.Reversed {
			ret = append(ret, FormatterColumn{
				Name:              "__emidx_" + index.ColumnName,
				Type:              FMT_TYPE_STR,
				Tags:              []string{"nullable"},
				IsInvisible:       true,
				AlwaysGeneratedAs: CheckNull(index.ColumnName, newSqlExpr(index.ColumnName).Reverse()).String(),
				Generator:         e,
				GeneratorData: ReverseIndexEmulatorIdxInfo{
					LinkedColumn: index.ColumnName}})
			e.Indexes = append(e.Indexes, index)
			(*indexes)[i].ColumnName = "__emidx_" + index.ColumnName
			(*indexes)[i].Reversed = false
		}
	}

	return ret
}

func (e *ReverseIndexEmulator) Apply(rows *map[string]*string) {
	if e.HasGeneratedAs {
		return
	}
	for _, index := range e.Indexes {
		if v, e := (*rows)[index.ColumnName]; e && v != nil {
			tmp := reverse_str(*v)
			(*rows)["__emidx_"+index.ColumnName] = &tmp
		}
	}
}

func EscapeString(input string) string {
	var output strings.Builder
	output.WriteByte('\'')
	for _, c := range input {
		switch c {
		case '\x00':
			output.WriteString("\\0")
		case '\'':
			output.WriteString("\\'")
		case '"':
			output.WriteString("\\\"")
		case '\x08':
			output.WriteString("\\b")
		case '\n':
			output.WriteString("\\n")
		case '\r':
			output.WriteString("\\r")
		case '\t':
			output.WriteString("\\t")
		case '\x1a': // the doc says "ASCII 26" but idk if it is decimal, hexadecimal or octal
			output.WriteString("\\Z")
		case '\\':
			output.WriteString("\\\\")
		default:
			output.WriteRune(c)
		}
	}

	output.WriteByte('\'')
	return output.String()
}

func (e *ReverseIndexEmulator) GetGeneratorInfo() GeneratorInfo {
	return GeneratorInfo{
		Name:          "reverse_index_emulator",
		VersionString: "1.0.0",
		VersionId:     0x010000,
	}
}
