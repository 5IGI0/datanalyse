package main

type GeneratorInfo struct {
	Name          string `json:"name"`
	VersionString string `json:"ver"`
	VersionId     uint32 `json:"ver_id"`
}

const (
	ANALYZER_EMAIL       = 1
	ANALYZER_USERNAME    = 2
	ANALYZER_PHONE       = 3
	ANALYZER_FACEBOOK_ID = 4
)

type Analyzer interface {
	Init(FormatterColumn, Formatter) ([]FormatterColumn, []FormatterIndex, error)
	Analyze(*map[string]*string) error
	GetGeneratorInfo() GeneratorInfo
	GetAnalyzerData() any
	GetAnalyzerType() int
}

func GetAnalyzer(tag string) Analyzer {
	switch tag {
	case "email":
		return &EmailAnalyzer{}
	case "username":
		return &UsernameAnalyzer{}
	case "realnames":
		return &UsernameAnalyzer{PrimaryType: "realnames"}
	case "phone":
		return &PhoneAnalyzer{}
	case "facebook_id":
		return &DummyAnalyzer{Data: nil, Generator: GeneratorInfo{
			Name:          "fbid_analyzer",
			VersionString: "1.0",
			VersionId:     0x010000,
		}, Type: ANALYZER_FACEBOOK_ID}
	}
	return nil
}
