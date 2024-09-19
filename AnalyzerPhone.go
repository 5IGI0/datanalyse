package main

type PhoneAnalyzer struct {
	ColumnName string
	Data       struct {
		Sanitized string `json:"sanitized"`
	}
}

func (u *PhoneAnalyzer) Init(col FormatterColumn, f Formatter) ([]FormatterColumn, []FormatterIndex, error) {
	u.ColumnName = col.Name

	u.Data.Sanitized = "__" + col.Name + "__phone_sanitized"

	return []FormatterColumn{
			{Name: u.Data.Sanitized,
				ForceString:       true,
				Tags:              []string{"nullable"},
				IsInvisible:       true,
				AlwaysGeneratedAs: CheckNull(col.Name, newSqlExpr(col.Name).OnlyNum()).String(),
				Generator:         u,
				GeneratorData: &GeneratorData{
					LinkedColumn: col.Name,
					Format:       "numeric",
					PrimaryType:  "phone_number",
					Version:      1}}},
		[]FormatterIndex{
			{ColumnName: u.Data.Sanitized,
				IndexName: "__" + col.Name + "__phone_sanitized",
				Reversed:  false}},
		nil
}

func (u *PhoneAnalyzer) Analyze(row *map[string]*string) error {
	if v, e := (*row)[u.ColumnName]; e && v != nil {
		sanitized := OnlyNum(*v)

		(*row)[u.Data.Sanitized] = &sanitized
	}

	return nil
}

func (u *PhoneAnalyzer) GetGeneratorInfo() GeneratorInfo {
	return GeneratorInfo{
		Name:          "phone_analyzer",
		VersionString: "1.0.0",
		VersionId:     0x010000,
	}
}

func (u *PhoneAnalyzer) GetAnalyzerData() any {
	return u.Data
}

func (u *PhoneAnalyzer) GetAnalyzerType() int {
	return ANALYZER_PHONE
}
