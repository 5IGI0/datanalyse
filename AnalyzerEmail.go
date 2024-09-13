package main

import (
	"strings"

	"golang.org/x/net/idna"
)

type EmailAnalyzer struct {
	ColumnName     string
	HasGeneratedAs bool
	Data           struct {
		ReverseLogin  string `json:"reverse_login"`
		BidirectLogin string `json:"bidirect_login"`
		ReverseDomain string `json:"ReverseDomain"`
	}
}

type EmailAnalyzerMetaColumnInfo struct {
	LinkedColumn string `json:"linked_col"`
	ColumnType   string `json:"coltyp"`
	Version      uint32 `json:"ver"`
}

func (a *EmailAnalyzer) Init(Column FormatterColumn, formatter Formatter) ([]FormatterColumn, []FormatterIndex, error) {
	a.ColumnName = Column.Name
	a.HasGeneratedAs = (formatter.GetFeatures() & FMT_FEATURE_GENERATED_AS) != 0

	a.Data.ReverseLogin = "__" + Column.Name + "__email_reverse_login"
	a.Data.BidirectLogin = "__" + Column.Name + "__email_bidirect_login"
	a.Data.ReverseDomain = "__" + Column.Name + "__email_reverse_domain"

	return []FormatterColumn{
			{
				Name:        a.Data.ReverseDomain,
				Type:        FMT_TYPE_STR,
				Tags:        []string{"nullable"},
				MaxLen:      Column.MaxLen,
				IsInvisible: true,
				Generator:   a,
				GeneratorData: EmailAnalyzerMetaColumnInfo{
					LinkedColumn: a.ColumnName,
					ColumnType:   "reverse_domain",
					Version:      1,
				}},
			{
				Name:        a.Data.ReverseLogin,
				Type:        FMT_TYPE_STR,
				Tags:        []string{"nullable"},
				MaxLen:      Column.MaxLen,
				IsInvisible: true,
				AlwaysGeneratedAs: CheckNull(a.ColumnName, newSqlExpr(a.ColumnName).
					SplitBefore("@").
					SplitBefore("+").
					OnlyAlphaNum().
					ToLower().
					Reverse()).
					String(),
				Generator: a,
				GeneratorData: EmailAnalyzerMetaColumnInfo{
					LinkedColumn: a.ColumnName,
					ColumnType:   "reverse_login",
					Version:      1,
				}},
			{
				Name:        a.Data.BidirectLogin,
				Type:        FMT_TYPE_STR,
				Tags:        []string{"nullable"},
				MaxLen:      Column.MaxLen * 2,
				IsInvisible: true,
				Generator:   a,
				GeneratorData: EmailAnalyzerMetaColumnInfo{
					LinkedColumn: a.ColumnName,
					ColumnType:   "bidirect_login",
					Version:      1,
				}}},
		[]FormatterIndex{
			{ColumnName: Column.Name,
				IndexName: "__" + Column.Name + "__email",
				Reversed:  false},
			{ColumnName: a.Data.ReverseLogin,
				IndexName: a.Data.ReverseLogin,
				Reversed:  false},
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

	if !a.HasGeneratedAs {
		tmp := reverse_str(login)
		(*row)[a.Data.ReverseLogin] = &tmp
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
