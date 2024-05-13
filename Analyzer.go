package main

type Analyzer interface {
	Init(string) ([]FormatterColumn, []FormatterIndex, error)
	Analyze(map[string]*string) error
}

func GetAnalyzer(tag string) Analyzer {
	switch tag {
	case "email":
		return &EmailAnalyzer{}
	}
	return nil
}
