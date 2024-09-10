package main

type PhoneAnalyzer struct {
	ColumnName     string
	HasGeneratedAs bool
	Data           struct {
		Sanitized string `json:"sanitized"`
	}
}

type PhoneAnalyzerMetaColumn struct {
	LinkedColumn string `json:"linked_col"`
	Version      uint32 `json:"ver"`
}

func (u *PhoneAnalyzer) Init(col FormatterColumn, f Formatter) ([]FormatterColumn, []FormatterIndex, error) {
	u.ColumnName = col.Name
	u.HasGeneratedAs = (f.GetFeatures() & FMT_FEATURE_GENERATED_AS) != 0

	u.Data.Sanitized = "__" + col.Name + "__phone_sanitized"

	return []FormatterColumn{
			{Name: u.Data.Sanitized,
				Type:              FMT_TYPE_STR,
				Tags:              []string{"nullable"},
				MaxLen:            col.MaxLen,
				MinLen:            0,
				IsInvisible:       true,
				AlwaysGeneratedAs: CheckNull(col.Name, newSqlExpr(col.Name).OnlyNum()).String(),
				Generator:         u,
				GeneratorData: PhoneAnalyzerMetaColumn{
					LinkedColumn: col.Name,
					Version:      1}}},
		[]FormatterIndex{
			{ColumnName: u.Data.Sanitized,
				IndexName: "__" + col.Name + "__phone_sanitized",
				Reversed:  false}},
		nil
}

func (u *PhoneAnalyzer) Analyze(row *map[string]*string) error {
	if u.HasGeneratedAs {
		return nil
	}
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
