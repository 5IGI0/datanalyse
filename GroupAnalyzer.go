package main

type GroupAnalyzer interface {
	Init([]FormatterColumn, Formatter) ([]FormatterColumn, []FormatterIndex, error)
	Analyze(*map[string]*string) error
	GetGeneratorInfo() GeneratorInfo
	GetAnalyzerData() any

	// this allows to know which columns matches when analyzers match virtual columns:
	// 1. Username analyzer matches _virtcol_12345
	// 2. VirtualColumnMap["_virtualcol_12345"] -> ["firstname", "lastname"]
	GetVirtualColumnMap() map[string][]string
}

func GetGroupAnalyzer(tag string) GroupAnalyzer {
	switch tag {
	case "realnames":
		return &RealnamesGroupAnalyzer{}
	}

	return nil
}
