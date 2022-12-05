package q2sql

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

// this builder is simplified version of the select
// builder from the github.com/Masterminds/squirrel project
// thanks to the project authors

type SelectBuilder struct {
	Columns      []Sqlizer
	FromPart     Sqlizer
	Joins        []Sqlizer
	WhereParts   []Sqlizer
	GroupBys     []string
	HavingParts  []Sqlizer
	OrderByParts []Sqlizer
	LimitPart    string
	OffsetPart   string
}

func (s *SelectBuilder) Select(columns []string) *SelectBuilder {
	if len(columns) == 0 {
		return s
	}
	s.Columns = append(s.Columns, Columns(columns))
	return s
}

func (s *SelectBuilder) From(from string) *SelectBuilder {
	s.FromPart = RawSql(from)
	return s
}

func (s *SelectBuilder) Join(clause Sqlizer) *SelectBuilder {
	s.Joins = append(s.Joins, clause)
	return s
}

func (s *SelectBuilder) Where(conditions ...Sqlizer) *SelectBuilder {
	s.WhereParts = append(s.WhereParts, conditions...)
	return s
}

func (s *SelectBuilder) OrderBy(clause Sqlizer) *SelectBuilder {
	s.OrderByParts = append(s.OrderByParts, clause)
	return s
}

func (s *SelectBuilder) GroupBy(groupBys ...string) *SelectBuilder {
	s.GroupBys = append(s.GroupBys, groupBys...)
	return s
}

func (s *SelectBuilder) Having(clause Sqlizer) *SelectBuilder {
	s.HavingParts = append(s.HavingParts, clause)
	return s
}

func (s *SelectBuilder) Limit(limit uint64) *SelectBuilder {
	s.LimitPart = strconv.FormatUint(limit, 10)
	return s
}

func (s *SelectBuilder) Offset(offset uint64) *SelectBuilder {
	s.OffsetPart = strconv.FormatUint(offset, 10)
	return s
}

//nolint:stylecheck // skip because this is an Interface implementation
func (s *SelectBuilder) ToSql() (sqlStr string, args []interface{}, err error) {
	sql := new(bytes.Buffer)
	args = make([]interface{}, 0)
	if len(s.Columns) == 0 {
		err = errors.New("select must have at least one column")
		return
	}
	sql.WriteString("SELECT ")
	args, err = appendToSQL(s.Columns, sql, ",", args)
	if err != nil {
		return "", nil, err
	}

	if s.FromPart != nil {
		sql.WriteString(" FROM ")
		args, err = appendToSQL([]Sqlizer{s.FromPart}, sql, "", args)
		if err != nil {
			return
		}
	}

	if len(s.Joins) > 0 {
		sql.WriteString(" ")
		args, err = appendToSQL(s.Joins, sql, " ", args)
		if err != nil {
			return
		}
	}

	if len(s.WhereParts) > 0 {
		sql.WriteString(" WHERE ")
		args, err = appendToSQL(s.WhereParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(s.GroupBys) > 0 {
		sql.WriteString(" GROUP BY ")
		sql.WriteString(strings.Join(s.GroupBys, ", "))
	}

	if len(s.HavingParts) > 0 {
		sql.WriteString(" HAVING ")
		args, err = appendToSQL(s.HavingParts, sql, " AND ", args)
		if err != nil {
			return
		}
	}

	if len(s.OrderByParts) > 0 {
		sql.WriteString(" ORDER BY ")
		args, err = appendToSQL(s.OrderByParts, sql, ", ", args)
		if err != nil {
			return
		}
	}

	if s.LimitPart != "" {
		sql.WriteString(" LIMIT ")
		sql.WriteString(s.LimitPart)
	}

	if s.OffsetPart != "" {
		sql.WriteString(" OFFSET ")
		sql.WriteString(s.OffsetPart)
	}
	sqlStr = sql.String()

	return sqlStr, args, err
}
