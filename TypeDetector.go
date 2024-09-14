package main

import (
	"slices"
	"sort"
	"strconv"
)

const (
	TYPE_DETECTOR_MAX_ENUM = 20
)

type TypeDetector struct {
	TotalRow  uint64
	NullRow   uint64
	Int8Row   uint64
	Int16Row  uint64
	Int32Row  uint64
	Int64Row  uint64
	UInt8Row  uint64
	UInt16Row uint64
	UInt32Row uint64
	// TODO: Uint64

	isASCII  bool
	IsEnum   bool
	EnumVals []string

	MaxLen int
	MinLen int
}

type TypeHint struct {
	Type     int
	Priority uint8
	Match    uint64
}

func (dt *TypeDetector) Init() {
	dt.MinLen = -1
	dt.MaxLen = 0

	dt.isASCII = true
	dt.IsEnum = true
	dt.EnumVals = make([]string, 0, TYPE_DETECTOR_MAX_ENUM)
}

func (dt *TypeDetector) Analyze(row *string) {
	dt.TotalRow++

	if row == nil {
		dt.NullRow++
		return
	}

	if len(*row) > int(dt.MaxLen) {
		dt.MaxLen = len(*row)
	}

	if dt.MinLen == -1 || dt.MinLen > len(*row) {
		dt.MinLen = len(*row)
	}

	/* enum related */
	if dt.IsEnum {
		if !slices.Contains(dt.EnumVals[:], *row) {
			if len(dt.EnumVals) == TYPE_DETECTOR_MAX_ENUM {
				dt.IsEnum = false
				dt.EnumVals = nil
			} else {
				dt.EnumVals = append(dt.EnumVals, *row)
			}
		}
	}

	intval, err := strconv.ParseInt(*row, 10, 64)
	if err == nil {
		dt.Int64Row++ // if there is no error, then it is in int64 range.
		if intval >= -0x80 && intval <= 0x7F {
			dt.Int8Row++
		}
		if intval >= 0 && intval <= 0xFF {
			dt.UInt8Row++
		}
		if intval >= -0x8000 && intval <= 0x7FFF {
			dt.Int16Row++
		}
		if intval >= 0 && intval <= 0xFFFF {
			dt.UInt16Row++
		}
		if intval >= -0x8000_0000 && intval <= 0x7FFF_FFFF {
			dt.Int32Row++
		}
		if intval >= 0 && intval <= 0xFFFF_FFFF {
			dt.UInt32Row++
		}
	} else {
		/* string related checks */

		if dt.isASCII {
			/* UTF-8 check */
			for _, r := range *row {
				if r >= 0x80 {
					dt.isASCII = false
				}
			}
		}
	}
}

func (dt *TypeDetector) GetMaxLen() int {
	return dt.MaxLen
}

func (dt *TypeDetector) GetMinLen() int {
	return dt.MinLen
}

func (dt *TypeDetector) GetType() int {
	if dt.NullRow == dt.TotalRow {
		return FMT_TYPE_INT8
	}
	t := dt.GetTypeHints()[0].Type

	if t == FMT_TYPE_STR && dt.IsEnum && len(dt.EnumVals) > 1 {
		return FMT_TYPE_ENUM
	}

	return t
}

func (dt *TypeDetector) GetEnumVals() []string {
	return dt.EnumVals
}

func (dt *TypeDetector) GetTypeHints() []TypeHint {
	var tmp = []TypeHint{
		{Type: FMT_TYPE_INT8, Priority: 110, Match: dt.Int8Row},
		{Type: FMT_TYPE_UINT8, Priority: 100, Match: dt.UInt8Row},
		{Type: FMT_TYPE_INT16, Priority: 90, Match: dt.Int16Row},
		{Type: FMT_TYPE_UINT16, Priority: 80, Match: dt.UInt16Row},
		{Type: FMT_TYPE_INT32, Priority: 70, Match: dt.Int32Row},
		{Type: FMT_TYPE_UINT32, Priority: 60, Match: dt.Int32Row},
		{Type: FMT_TYPE_INT64, Priority: 50, Match: dt.Int64Row},
		{Type: FMT_TYPE_STR, Priority: 0, Match: dt.TotalRow}}

	sort.Slice(tmp, func(a, b int) bool {
		if tmp[a].Match == tmp[b].Match {
			return uint64(tmp[a].Priority) > uint64(tmp[b].Priority)
		} else {
			return tmp[a].Match > tmp[b].Match
		}
	})

	return tmp
}

func (dt *TypeDetector) GetTags() []string {
	ret := []string{}

	if dt.NullRow == dt.TotalRow {
		ret = append(ret, "unused")
		return ret
	} else if dt.NullRow == 0 {
		ret = append(ret, "nonnull")
	} else {
		ret = append(ret, "nullable")
	}

	return ret
}

func (dt *TypeDetector) IsNullable() bool {
	return dt.NullRow != 0
}

func (dt *TypeDetector) IsConstant() bool {
	return dt.IsEnum && len(dt.EnumVals) <= 1
}

func (dt *TypeDetector) GetConstVal() *string {
	if len(dt.EnumVals) != 0 {
		return &dt.EnumVals[0]
	}
	return nil
}

func (dt *TypeDetector) IsASCII() bool {
	return dt.isASCII
}
