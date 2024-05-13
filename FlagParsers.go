package main

import (
	"errors"
	"strconv"
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

type ColumnLenVar struct {
	Lens *map[string]*ColumnLen
}

func (v ColumnLenVar) String() string {
	return ""
}

func (v ColumnLenVar) Set(s string) error {
	splitted := strings.Split(s, ":")
	var val ColumnLen
	var err error

	if len(splitted) != 2 && len(splitted) != 3 {
		return errors.New("invalid")
	}

	if val.Max, err = strconv.Atoi(splitted[1]); err != nil {
		return err
	}

	if len(splitted) == 2 {
		val.Min = -1
	} else {
		if val.Min, err = strconv.Atoi(splitted[2]); err != nil {
			return err
		}
	}

	(*v.Lens)[splitted[0]] = &val
	return nil
}
