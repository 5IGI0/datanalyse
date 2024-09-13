package main

import (
	"bytes"
	"fmt"
)

type SqlExprGenerator struct {
	CurrentVal string
}

func newSqlExpr(input_column string) SqlExprGenerator {
	return SqlExprGenerator{
		CurrentVal: "`" + input_column + "`",
	}
}

func CheckNull(input_column string, expr SqlExprGenerator) SqlExprGenerator {
	return newSqlExpr(input_column).IfNotNull(expr)
}

func EmptyIfNull(input_column string) SqlExprGenerator {
	return newSqlExpr(input_column).IfNull(SqlExprGenerator{CurrentVal: "''"})
}

func SqlConcat(vals ...SqlExprGenerator) SqlExprGenerator {
	var output bytes.Buffer

	output.WriteString("CONCAT(")

	for i, val := range vals {
		if i != 0 {
			output.WriteByte(',')
		}

		output.WriteString(val.String())
	}

	output.WriteByte(')')
	return SqlExprGenerator{CurrentVal: output.String()}
}

func (e SqlExprGenerator) IfNull(in SqlExprGenerator) SqlExprGenerator {
	return SqlExprGenerator{
		CurrentVal: "IFNULL(" + e.CurrentVal + ", " + in.String() + ")",
	}
}

func (e SqlExprGenerator) SplitBefore(text string) SqlExprGenerator {
	return SqlExprGenerator{
		CurrentVal: "SUBSTRING_INDEX(" + e.CurrentVal + ", " + EscapeString(text) + " , 1)",
	}
}

func (e SqlExprGenerator) Reverse() SqlExprGenerator {
	return SqlExprGenerator{
		CurrentVal: "REVERSE(" + e.CurrentVal + ")",
	}
}

func (e SqlExprGenerator) ToLower() SqlExprGenerator {
	return SqlExprGenerator{
		CurrentVal: "LOWER(" + e.CurrentVal + ")",
	}
}

func (e SqlExprGenerator) IfNotNull(expr SqlExprGenerator) SqlExprGenerator {
	return SqlExprGenerator{
		CurrentVal: "IF (" + e.CurrentVal + " IS NOT NULL, " + expr.String() + ", NULL)",
	}
}

func (e SqlExprGenerator) OnlyAlphaNum() SqlExprGenerator {
	return SqlExprGenerator{
		CurrentVal: "REGEXP_REPLACE(" + e.CurrentVal + ", '[^a-zA-Z0-9]', '')",
	}
}

func (e SqlExprGenerator) OnlyNum() SqlExprGenerator {
	return SqlExprGenerator{
		CurrentVal: "REGEXP_REPLACE(" + e.CurrentVal + ", '[^0-9]', '')",
	}
}

func (e SqlExprGenerator) MaxLen(len uint) SqlExprGenerator {
	return SqlExprGenerator{
		CurrentVal: "LEFT(" + e.CurrentVal + ", " + fmt.Sprint(len) + ")",
	}
}

func (e SqlExprGenerator) String() string {
	return e.CurrentVal
}
