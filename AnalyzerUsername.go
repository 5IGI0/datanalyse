package main

type UsernameAnalyzer struct {
	ColumnName string
	Data       struct {
		Sanitized string `json:"sanitized"`
		Bidirect  string `json:"bidirect"`
	}
}

type UsernameAnalyzerMetaColumnInfo struct {
	LinkedColumn string `json:"linked_col"`
	ColumnType   string `json:"coltyp"`
	Version      uint32 `json:"ver"`
}

func (u *UsernameAnalyzer) Init(col FormatterColumn) ([]FormatterColumn, []FormatterIndex, error) {
	u.ColumnName = col.Name

	u.Data.Sanitized = "__" + col.Name + "__username_sanitized"
	u.Data.Bidirect = "__" + col.Name + "__username_bidirect"

	return []FormatterColumn{
			{Name: u.Data.Sanitized,
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
			{Name: u.Data.Bidirect,
				Type:        FMT_TYPE_STR,
				Tags:        []string{"nullable"},
				MaxLen:      col.MaxLen * 2,
				MinLen:      0,
				IsInvisible: true,
				Generator:   u,
				GeneratorData: UsernameAnalyzerMetaColumnInfo{
					LinkedColumn: col.Name,
					ColumnType:   "bidirect",
					Version:      1}}},
		[]FormatterIndex{
			{ColumnName: u.Data.Sanitized,
				IndexName: "__" + col.Name + "__username_sanitized",
				Reversed:  false},
			{ColumnName: u.Data.Sanitized,
				IndexName: "__" + col.Name + "__username_reverse_sanitized",
				Reversed:  true},
			{ColumnName: u.Data.Bidirect,
				IndexName: "__" + col.Name + "__username_bidirect",
				Reversed:  true}},
		nil
}

func (u *UsernameAnalyzer) Analyze(row *map[string]*string) error {
	if v, e := (*row)[u.ColumnName]; e && v != nil {
		sanitized := OnlyAlphaNum(*v)

		(*row)[u.Data.Sanitized] = &sanitized
		bidirect := BidirectionalizeTextA(sanitized)
		(*row)[u.Data.Bidirect] = &bidirect
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

func (u *UsernameAnalyzer) GetAnalyzerData() any {
	return u.Data
}
