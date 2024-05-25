package main

// used for columns that don't need specific analysis
type DummyAnalyzer struct {
	Data      any
	Generator GeneratorInfo
}

func (*DummyAnalyzer) Init(Column FormatterColumn) ([]FormatterColumn, []FormatterIndex, error) {
	return nil, []FormatterIndex{
		{
			ColumnName: Column.Name,
			IndexName:  "__" + Column.Name + "_idx",
			Reversed:   false},
	}, nil
}

func (*DummyAnalyzer) Analyze(*map[string]*string) error {
	return nil
}

func (d *DummyAnalyzer) GetGeneratorInfo() GeneratorInfo {
	return d.Generator
}

func (d *DummyAnalyzer) GetAnalyzerData() any {
	return d.Data
}
