package main

type UsernameAnalyzer struct {
	ColumnName  string
	PrimaryType string
	Data        struct {
		Sanitized string `json:"sanitized"`
		Bidirect  string `json:"bidirect"`
	}
}

func (u *UsernameAnalyzer) Init(col FormatterColumn, f Formatter) ([]FormatterColumn, []FormatterIndex, error) {
	u.ColumnName = col.Name

	u.Data.Sanitized = "__" + col.Name + "__username_sanitized"
	u.Data.Bidirect = "__" + col.Name + "__username_bidirect"

	if u.PrimaryType == "" {
		u.PrimaryType = "username"
	}

	return []FormatterColumn{
			{Name: u.Data.Sanitized,
				ForceString: true,
				Tags:        []string{"nullable"},
				IsInvisible: true,
				Generator:   u,
				GeneratorData: &GeneratorData{
					LinkedColumn: col.Name,
					Format:       "sanitized",
					PrimaryType:  u.PrimaryType,
					Version:      1}},
			{Name: u.Data.Bidirect,
				ForceString: true,
				Tags:        []string{"nullable"},
				IsInvisible: true,
				Generator:   u,
				GeneratorData: &GeneratorData{
					LinkedColumn: col.Name,
					Format:       "bidirect_sanitized",
					PrimaryType:  u.PrimaryType,
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
				Reversed:  false}},
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

func (u *UsernameAnalyzer) GetAnalyzerType() int {
	return ANALYZER_USERNAME
}
