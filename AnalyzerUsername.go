package main

type UsernameAnalyzer struct {
	ColumnName string
}

type UsernameAnalyzerMetaColumnInfo struct {
	LinkedColumn string `json:"linked_col"`
	ColumnType   string `json:"coltyp"`
	Version      uint32 `json:"ver"`
}

func (u *UsernameAnalyzer) Init(col FormatterColumn) ([]FormatterColumn, []FormatterIndex, error) {
	u.ColumnName = col.Name

	return []FormatterColumn{
			{Name: "__" + col.Name + "__username_sanitized",
				Type:        FMT_TYPE_STR,
				Tags:        []string{"nullable"},
				MaxLen:      col.MaxLen,
				MinLen:      0,
				IsInvisible: true,
				Generator:   u,
				GeneratorData: UsernameAnalyzerMetaColumnInfo{
					LinkedColumn: col.Name,
					ColumnType:   "sanitized",
					Version:      1}},
			{Name: "__" + col.Name + "__username_bidirect",
				Type:        FMT_TYPE_STR,
				Tags:        []string{"nullable"},
				MaxLen:      col.MaxLen,
				MinLen:      0,
				IsInvisible: true,
				Generator:   u,
				GeneratorData: UsernameAnalyzerMetaColumnInfo{
					LinkedColumn: col.Name,
					ColumnType:   "bidirect",
					Version:      1}}},
		[]FormatterIndex{
			{ColumnName: "__" + col.Name + "__username_sanitized",
				IndexName: "__" + col.Name + "__username_sanitized",
				Reversed:  false},
			{ColumnName: "__" + col.Name + "__username_sanitized",
				IndexName: "__" + col.Name + "__username_reverse_sanitized",
				Reversed:  true},
			{ColumnName: "__" + col.Name + "__username_bidirect",
				IndexName: "__" + col.Name + "__username_bidirect",
				Reversed:  true}},
		nil
}

func (u *UsernameAnalyzer) Analyze(row *map[string]*string) error {
	if v, e := (*row)[u.ColumnName]; e && v != nil {
		sanitized := OnlyAlphaNum(*v)

		(*row)["__"+u.ColumnName+"__username_sanitized"] = &sanitized
		bidirect := BidirectionalizeTextA(sanitized)
		(*row)["__"+u.ColumnName+"__username_bidirect"] = &bidirect
	}

	return nil
}

func (u *UsernameAnalyzer) GetGeneratorInfo() GeneratorInfo {
	return GeneratorInfo{
		Name:          "username_analyzer",
		VersionString: "1.0.0",
		VersionId:     0x010000,
	}
}
