package main

import (
	"strings"

	"golang.org/x/net/idna"
)

type EmailAnalyzer struct {
	ColumnName string
	Data       struct {
		SanitizedLogin string `json:"reverse_login"`
		BidirectLogin  string `json:"bidirect_login"`
		ReverseDomain  string `json:"ReverseDomain"`
	}
}

func (a *EmailAnalyzer) Init(Column FormatterColumn, formatter Formatter) ([]FormatterColumn, []FormatterIndex, error) {
	a.ColumnName = Column.Name

	a.Data.SanitizedLogin = "__" + Column.Name + "__email_sanitized_login"
	a.Data.BidirectLogin = "__" + Column.Name + "__email_bidirect_login"
	a.Data.ReverseDomain = "__" + Column.Name + "__email_reverse_domain"

	return []FormatterColumn{
			{
				Name:        a.Data.ReverseDomain,
				ForceString: true,
				Tags:        []string{"nullable"},
				IsInvisible: true,
				Generator:   a,
				GeneratorData: &GeneratorData{
					LinkedColumn: a.ColumnName,
					Format:       "reverse_domain",
					PrimaryType:  "email_domain",
					Tags:         []string{"domain"},
					Version:      1,
				}},
			{
				Name:        a.Data.SanitizedLogin,
				ForceString: true,
				Tags:        []string{"nullable"},
				IsInvisible: true,
				Generator:   a,
				GeneratorData: &GeneratorData{
					LinkedColumn: a.ColumnName,
					Format:       "sanitized",
					PrimaryType:  "email_login",
					Tags:         []string{"username"},
					Version:      1,
				}},
			{
				Name:        a.Data.BidirectLogin,
				ForceString: true,
				Tags:        []string{"nullable"},
				IsInvisible: true,
				Generator:   a,
				GeneratorData: &GeneratorData{
					LinkedColumn: a.ColumnName,
					Format:       "bidirect_sanitized",
					PrimaryType:  "email_login",
					Tags:         []string{"username"},
					Version:      1,
				}}},
		[]FormatterIndex{
			{ColumnName: Column.Name,
				IndexName: "__" + Column.Name + "__email",
				Reversed:  false},
			{ColumnName: a.Data.SanitizedLogin,
				IndexName: a.Data.SanitizedLogin,
				Reversed:  false},
			{ColumnName: a.Data.SanitizedLogin,
				IndexName: a.Data.SanitizedLogin + "__reverse",
				Reversed:  true},
			{ColumnName: a.Data.BidirectLogin,
				IndexName: a.Data.BidirectLogin,
				Reversed:  false},
			{ColumnName: a.Data.ReverseDomain,
				IndexName: a.Data.ReverseDomain,
				Reversed:  false}},
		nil
}

func (a *EmailAnalyzer) Analyze(row *map[string]*string) error {
	v, e := (*row)[a.ColumnName]

	if !e || v == nil || strings.IndexByte(*v, '@') == -1 {
		return nil
	}

	vv := *v

	var login string
	if strings.IndexByte(vv, '+') != -1 && strings.IndexByte(vv, '+') < strings.IndexByte(vv, '@') {
		login = vv[:strings.IndexByte(vv, '+')]
	} else {
		login = vv[:strings.IndexByte(vv, '@')]
	}
	login = OnlyAlphaNum(login)

	{
		domain, err := idna.ToASCII(vv[strings.IndexByte(vv, '@')+1:])
		AssertError(err)
		tmp := reverse_str(strings.ToLower(domain))
		(*row)[a.Data.ReverseDomain] = &tmp
	}

	{
		tmp := login
		(*row)[a.Data.SanitizedLogin] = &tmp
	}

	{
		tmp := BidirectionalizeTextA(login)
		(*row)[a.Data.BidirectLogin] = &tmp
	}

	return nil
}

func (a *EmailAnalyzer) GetGeneratorInfo() GeneratorInfo {
	return GeneratorInfo{
		Name:          "email_analyzer",
		VersionString: "1.0",
		VersionId:     0x010000,
	}
}

func (a *EmailAnalyzer) GetAnalyzerData() any {
	return a.Data
}

func (a *EmailAnalyzer) GetAnalyzerType() int {
	return ANALYZER_PHONE
}
