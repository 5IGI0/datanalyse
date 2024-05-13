package main

import "strings"

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

type ColumnTypes struct {
	types *map[string][]string
}

func (v ColumnTypes) String() string {
	return ""
}

func (v ColumnTypes) Set(s string) error {
	tmp := strings.Split(s, ":")
	(*v.types)[tmp[0]] = tmp[1:]
	return nil
}
