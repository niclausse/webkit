package typex

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"strings"
)

// JSONQueryExpression json query expression, implements clause.Expression interface to use as querier
type JSONQueryExpression struct {
	column      string
	keys        []string
	hasKeys     bool
	equals      bool
	equalsValue interface{}
	extract     bool
	path        string
}

func JSONQuery(column string) *JSONQueryExpression {
	return &JSONQueryExpression{column: column}
}

func (jq *JSONQueryExpression) Extract(path string) *JSONQueryExpression {
	jq.extract = true
	jq.path = path
	return jq
}

func (jq *JSONQueryExpression) HasKey(keys ...string) *JSONQueryExpression {
	jq.keys = keys
	jq.hasKeys = true
	return jq
}

func (jq *JSONQueryExpression) Equals(value interface{}, keys ...string) *JSONQueryExpression {
	jq.keys = keys
	jq.equalsValue = value
	jq.equals = true
	return jq
}

func (jq *JSONQueryExpression) Build(builder clause.Builder) {
	stmt, ok := builder.(*gorm.Statement)
	if !ok {
		return
	}

	switch stmt.Dialector.Name() {
	case "mysql", "sqlite":
		switch {
		case jq.extract:
			builder.WriteString("JSON_EXTRACT(")
			builder.WriteQuoted(jq.column)
			builder.WriteByte(',')
			builder.AddVar(stmt, jq.path)
			builder.WriteString(")")
		case jq.hasKeys:
			if len(jq.keys) > 0 {
				builder.WriteString("JSON_EXTRACT(")
				builder.WriteQuoted(jq.column)
				builder.WriteByte(',')
				builder.AddVar(stmt, jsonQueryJoin(jq.keys))
				builder.WriteString(") IS NOT NULL")
			}
		case jq.equals:
			if len(jq.keys) > 0 {
				builder.WriteString("JSON_EXTRACT(")
				builder.WriteQuoted(jq.column)
				builder.WriteByte(',')
				builder.AddVar(stmt, jsonQueryJoin(jq.keys))
				builder.WriteString(") =")
				if value, ok := jq.equalsValue.(bool); ok {
					builder.WriteString(strconv.FormatBool(value))
				} else {
					stmt.AddVar(builder, jq.equalsValue)
				}
			}
		}
	}
}

const prefix = "$."

func jsonQueryJoin(keys []string) string {
	if len(keys) == 1 {
		return prefix + keys[0]
	}

	n := len(prefix)
	n += len(keys) - 1
	for i := 0; i < len(keys); i++ {
		n += len(keys[i])
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(prefix)
	b.WriteString(keys[0])
	for _, key := range keys[1:] {
		b.WriteString(".")
		b.WriteString(key)
	}
	return b.String()
}
