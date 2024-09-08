package main

type GeneratorInfo struct {
	Name          string `json:"name"`
	VersionString string `json:"ver"`
	VersionId     uint32 `json:"ver_id"`
}

type Analyzer interface {
	Init(FormatterColumn, Formatter) ([]FormatterColumn, []FormatterIndex, error)
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
	case "facebook_id":
		return &DummyAnalyzer{Data: nil, Generator: GeneratorInfo{
			Name:          "fbid_analyzer",
			VersionString: "1.0",
			VersionId:     0x010000,
		}}
	}
	return nil
}
