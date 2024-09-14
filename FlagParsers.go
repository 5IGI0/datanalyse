package main

import (
	"errors"
	"strings"
)

type MultipleVars struct {
	vars *[]string
}

func (v MultipleVars) String() string {
	return ""
}

func (v MultipleVars) Set(s string) error {
	*v.vars = append(*v.vars, s)
	return nil
}

type ColumnTagsVar struct {
	vars *map[string]*ColumnInfo
}

func (v ColumnTagsVar) String() string {
	return ""
}

func (v ColumnTagsVar) Set(s string) error {
	tmp := strings.Split(s, ":")
	if len(tmp) < 2 {
		return errors.New("invalid")
	}

	if vv, e := (*v.vars)[tmp[0]]; e {
		vv.Tags = append(vv.Tags, tmp[1:]...)
	} else {
		t := ColumnInfo{}
		t.Init()
		t.Tags = tmp[1:]
		(*v.vars)[tmp[0]] = &t
	}

	return nil
}

type DatasetMetaVar struct {
	vars *map[string]string
}

func (v DatasetMetaVar) String() string {
	return ""
}

func (v DatasetMetaVar) Set(s string) error {
	tmp := strings.IndexByte(s, ':')
	if tmp < 0 {
		return errors.New("invalid")
	}

	(*v.vars)[s[:tmp]] = s[tmp+1:]

	return nil
}

type GroupVar struct {
	groups *[]GroupInfo
}

func (v GroupVar) String() string {
	return ""
}

func (v GroupVar) Set(s string) error {
	tmp := strings.Split(s, ":")
	if len(tmp) < 3 {
		return errors.New("invalid")
	}

	*v.groups = append(*v.groups, GroupInfo{
		Kind:   tmp[0],
		Fields: tmp[1:]})

	return nil
}
