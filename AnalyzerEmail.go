package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/idna"
)

type EmailAnalyzer struct {
	ColumnName string
}

type EmailAnalyzerMetaColumnInfo struct {
	LinkedColumn string `json:"linked_col"`
	ColumnType   string `json:"coltyp"`
	Version      uint32 `json:"ver"`
}

func (a *EmailAnalyzer) Init(Column FormatterColumn) ([]FormatterColumn, []FormatterIndex, error) {
	a.ColumnName = Column.Name

	return []FormatterColumn{
			{
				Name:        "__" + Column.Name + "__email_sanitized",
				Type:        FMT_TYPE_STR,
				Tags:        []string{"nullable"},
				MaxLen:      Column.MaxLen,
				IsInvisible: true,
				Generator:   a,
				GeneratorData: EmailAnalyzerMetaColumnInfo{
					LinkedColumn: a.ColumnName,
					ColumnType:   "sanitized",
					Version:      1,
				}},
			{
				Name:        "__" + Column.Name + "__email_reverse_login",
				Type:        FMT_TYPE_STR,
				Tags:        []string{"nullable"},
				MaxLen:      Column.MaxLen,
				IsInvisible: true,
				Generator:   a,
				GeneratorData: EmailAnalyzerMetaColumnInfo{
					LinkedColumn: a.ColumnName,
					ColumnType:   "reverse_login",
					Version:      1,
				}},
			{
				Name:        "__" + Column.Name + "__email_bidirect_login",
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
			{ColumnName: "__" + Column.Name + "__email_sanitized",
				IndexName: "__" + Column.Name + "__email_reverse_sanitized",
				Reversed:  true},
			{ColumnName: "__" + Column.Name + "__email_reverse_login",
				IndexName: "__" + Column.Name + "__email_reverse_login",
				Reversed:  false},
			{ColumnName: "__" + Column.Name + "__email_sanitized",
				IndexName: "__" + Column.Name + "__email_sanitized",
				Reversed:  false},
			{ColumnName: "__" + Column.Name + "__email_bidirect_login",
				IndexName: "__" + Column.Name + "__email_bidirect_login",
				Reversed:  false}},
		nil
}

func (a *EmailAnalyzer) Analyze(row *map[string]*string) error {
	v, e := (*row)[a.ColumnName]

	if !e || v == nil || strings.IndexByte(*v, '@') == -1 {
		return nil
	}

	vv := *v

	// TODO: error
	domain, _ := idna.ToASCII(vv[strings.IndexByte(vv, '@')+1:])
	domain = strings.ToLower(domain)
	var login string
	if strings.IndexByte(vv, '+') != -1 && strings.IndexByte(vv, '+') < strings.IndexByte(vv, '@') {
		login = vv[:strings.IndexByte(vv, '+')]
	} else {
		login = vv[:strings.IndexByte(vv, '@')]
	}
	login = OnlyAlphaNum(login)

	{
		tmp := fmt.Sprint(login, "@", domain)
		(*row)["__"+a.ColumnName+"__email_sanitized"] = &tmp
	}

	{
		tmp := reverse_str(login)
		(*row)["__"+a.ColumnName+"__email_reverse_login"] = &tmp
	}

	{
		tmp := BidirectionalizeTextA(login)
		(*row)["__"+a.ColumnName+"__email_bidirect_login"] = &tmp
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
