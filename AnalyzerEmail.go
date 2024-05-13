package main

import (
	"fmt"
	"strings"

	"golang.org/x/net/idna"
)

type EmailAnalyzer struct {
	ColumnName string
}

func (a *EmailAnalyzer) Init(ColumnName string) ([]FormatterColumn, []FormatterIndex, error) {
	a.ColumnName = ColumnName

	return []FormatterColumn{
			{
				Name: "__" + ColumnName + "__email_sanitized",
				Type: FMT_TYPE_STR,
				Tags: []string{"nullable"}},
			{
				Name: "__" + ColumnName + "__email_reverse_login",
				Type: FMT_TYPE_STR,
				Tags: []string{"nullable"}}},
		[]FormatterIndex{
			{ColumnName: "__" + ColumnName + "__email_sanitized",
				IndexName: "__" + ColumnName + "__email_reverse_sanitized",
				Reversed:  true},
			{ColumnName: "__" + ColumnName + "__email_reverse_login",
				IndexName: "__" + ColumnName + "__email_reverse_login",
				Reversed:  false},
			{ColumnName: "__" + ColumnName + "__email_sanitized",
				IndexName: "__" + ColumnName + "__email_sanitized",
				Reversed:  false}},
		nil
}

func (a *EmailAnalyzer) Analyze(row map[string]*string) error {
	v, e := row[a.ColumnName]

	if !e || v == nil || strings.IndexByte(*v, '@') == -1 {
		return nil
	}

	vv := *v

	// TODO: error
	domain, _ := idna.ToASCII(vv[strings.IndexByte(vv, '@')+1:])
	var login string
	if strings.IndexByte(vv, '+') != -1 {
		login = vv[:strings.IndexByte(vv, '+')]
	} else {
		login = vv[:strings.IndexByte(vv, '@')]
	}
	login = OnlyAlphaNum(login)

	{
		tmp := fmt.Sprint(login, "@", domain)
		row["__"+a.ColumnName+"_sanitized"] = &tmp
	}

	{
		tmp := reverse_str(login)
		row["__"+a.ColumnName+"_reverse_login"] = &tmp
	}

	return nil
}
