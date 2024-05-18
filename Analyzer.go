package main

type GeneratorInfo struct {
	Name          string `json:"name"`
	VersionString string `json:"ver"`
	VersionId     uint32 `json:"ver_id"`
}

type Analyzer interface {
	Init(FormatterColumn) ([]FormatterColumn, []FormatterIndex, error)
	Analyze(*map[string]*string) error
	GetGeneratorInfo() GeneratorInfo
	GetAnalyzerData() any
}

func GetAnalyzer(tag string) Analyzer {
	switch tag {
	case "email":
		return &EmailAnalyzer{}
	case "username":
		return &UsernameAnalyzer{}
	}
	return nil
}
