package main

import (
	"crypto/md5"
	"encoding/hex"
)

type RealnamesGroupAnalyzer struct {
	SubAnalyzer    Analyzer
	NeedRealColumn bool
	Data           struct {
		VirtualColName string   `json:"virtcol_name"`
		Columns        []string `json:"cols"`
	}
}

func (ga *RealnamesGroupAnalyzer) Init(cols []FormatterColumn, f Formatter) ([]FormatterColumn, []FormatterIndex, error) {

	/* generate column info */
	max_len := 0
	min_len := 0
	column_slug := "realnames"
	for _, column := range cols {
		ga.Data.Columns = append(ga.Data.Columns, column.Name)
		column_slug += ":" + column.Name

		if max_len >= 0 {
			max_len += column.MaxLen
			if column.MaxLen == -1 {
				max_len = -1
			}
		}

		if min_len >= 0 {
			min_len += column.MinLen
			if column.MinLen == -1 {
				max_len = -1
			}
		}

	}

	sum := md5.Sum([]byte(column_slug))
	ga.Data.VirtualColName = "__virtcol_" + hex.EncodeToString(sum[:5])

	VirtCol := FormatterColumn{
		Name:          ga.Data.VirtualColName,
		Type:          FMT_TYPE_STR,
		Tags:          []string{"nullable"},
		MaxLen:        max_len,
		MinLen:        min_len,
		IsLenFixed:    max_len == min_len,
		IsInvisible:   true,
		Generator:     ga,
		GeneratorData: ga.Data}

	/* pass it to the realnames analyzer */
	ga.SubAnalyzer = GetAnalyzer("realnames")
	Assert(ga.SubAnalyzer != nil)
	output_columns, output_indexes, err := ga.SubAnalyzer.Init(VirtCol, f)
	if err != nil {
		return nil, nil, err
	}

	/* now check if the virtual column actually need to be generated */
	if f.GetFeatures()&FMT_FEATURE_GENERATED_AS != 0 {
		/*	if the formatter support generated columns
			and one of the output column have generated value, then we need to generate it */
		for _, column := range output_columns {
			if column.AlwaysGeneratedAs != "" {
				ga.NeedRealColumn = true
				break
			}
		}
	}

	for _, index := range output_indexes {
		/* if one of the indexes indexes the virtual column, then we need to generate it */
		if index.ColumnName == ga.Data.VirtualColName {
			ga.NeedRealColumn = true
			break
		}
	}

	/* if we need to generate it, add it to the returned columns */
	if ga.NeedRealColumn {
		output_columns = append(output_columns, VirtCol)

		/* if it supports always generated columns, then just use it */
		if f.GetFeatures()&FMT_FEATURE_GENERATED_AS != 0 {
			exprs := []SqlExprGenerator{}
			for _, column := range cols {
				exprs = append(exprs, EmptyIfNull(column.Name))
			}

			VirtCol.AlwaysGeneratedAs = SqlConcat(exprs...).String()
		}
	}
	return output_columns, output_indexes, nil
}

func (ga *RealnamesGroupAnalyzer) Analyze(row *map[string]*string) error {
	tmp_val := ""

	for _, column_name := range ga.Data.Columns {
		if v, e := (*row)[column_name]; e && v != nil {
			tmp_val += *v
		}
	}

	(*row)[ga.Data.VirtualColName] = &tmp_val
	return ga.SubAnalyzer.Analyze(row)
}

func (ga *RealnamesGroupAnalyzer) GetGeneratorInfo() GeneratorInfo {
	return GeneratorInfo{
		Name:          "realname_group_analyzer",
		VersionString: "1.0.0",
		VersionId:     0x010000,
	}
}

func (ga *RealnamesGroupAnalyzer) GetAnalyzerData() any {
	return ga.Data
}

func (ga *RealnamesGroupAnalyzer) GetVirtualColumnMap() map[string][]string {
	return map[string][]string{
		ga.Data.VirtualColName: ga.Data.Columns}
}