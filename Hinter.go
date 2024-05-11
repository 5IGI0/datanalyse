package main

import (
	"sort"
	"strconv"
	"strings"
)

type Hinter struct {
	TotalRow  uint64
	NullRow   uint64
	Int8Row   uint64
	Int16Row  uint64
	Int32Row  uint64
	Int64Row  uint64
	UInt8Row  uint64
	UInt16Row uint64
	UInt32Row uint64
	Email     uint64
	// TODO: UInt64Row uint64
}

type TypeHint struct {
	TypeName string
	Priority uint8
	Match    uint64
}

// TODO: boolean / sex / country / ...
// TODO: optimisation
func (h *Hinter) Analyze(row *string) {
	h.TotalRow++

	if row == nil {
		h.NullRow++
		return
	}

	intval, err := strconv.ParseInt(*row, 10, 64)
	if err == nil {
		h.Int64Row++ // if there is no error, then it is in int64 range.
		if intval >= -0x80 && intval <= 0x7F {
			h.Int8Row++
		}
		if intval >= 0 && intval <= 0xFF {
			h.UInt8Row++
		}
		if intval >= -0x8000 && intval <= 0x7FFF {
			h.Int16Row++
		}
		if intval >= 0 && intval <= 0xFFFF {
			h.UInt16Row++
		}
		if intval >= -0x8000_0000 && intval <= 0x7FFF_FFFF {
			h.Int32Row++
		}
		if intval >= 0 && intval <= 0xFFFF_FFFF {
			h.UInt32Row++
		}
	} else { // if it's number, it's obviously not an email or something like that

		// email IMPROVEMENT: detect anti-bot format (webmaster [at] example <dot> com)
		if strings.IndexByte(*row, '@') != -1 &&
			strings.IndexByte(*row, ' ') == -1 {
			h.Email++
		}
	}
}

func (h *Hinter) GetType() string {
	if h.NullRow == h.TotalRow {
		return "int8"
	}
	return h.GetTypeHints()[0].TypeName
}

func (h *Hinter) GetTypeHints() []TypeHint {
	var tmp = []TypeHint{
		{TypeName: "int8", Priority: 110, Match: h.Int8Row},
		{TypeName: "uint8", Priority: 100, Match: h.UInt8Row},
		{TypeName: "int16", Priority: 90, Match: h.Int16Row},
		{TypeName: "uint16", Priority: 80, Match: h.UInt16Row},
		{TypeName: "int32", Priority: 70, Match: h.Int32Row},
		{TypeName: "uint32", Priority: 60, Match: h.Int32Row},
		{TypeName: "int64", Priority: 50, Match: h.Int64Row},
		{TypeName: "str", Priority: 0, Match: h.TotalRow}}

	sort.Slice(tmp, func(a, b int) bool {
		if tmp[a].Match == tmp[b].Match {
			return uint64(tmp[a].Priority) > uint64(tmp[b].Priority)
		} else {
			return tmp[a].Match > tmp[b].Match
		}
	})

	return tmp
}

func (h *Hinter) get_tags() []string {
	ret := []string{}

	if h.NullRow == h.TotalRow {
		ret = append(ret, "unused")
		return ret
	} else if h.NullRow == 0 {
		ret = append(ret, "nonnull")
	} else {
		ret = append(ret, "nullable")
	}

	// TODO: find the right trigger value
	if h.Email*2 >= (h.TotalRow - h.NullRow) {
		ret = append(ret, "email")
	}

	return ret
}
