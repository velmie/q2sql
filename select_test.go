package q2sql

import (
	"fmt"
	"reflect"
	"testing"
)

type selectBuilderTest struct {
	b     *SelectBuilder
	query string
	args  []interface{}
	err   bool
}

var selectBuilderTests = []selectBuilderTest{
	{
		b:     new(SelectBuilder).Select([]string{"f1", "f2"}),
		query: "SELECT f1, f2",
		args:  []interface{}{},
		err:   false,
	},
	{
		b:     new(SelectBuilder),
		query: "",
		args:  []interface{}{},
		err:   true,
	},
	{
		b: new(SelectBuilder).
			Select([]string{"*"}).
			From("tbl").
			Where(&Eq{"f1", "val"}),
		query: "SELECT * FROM tbl WHERE f1 = ?",
		args:  []interface{}{"val"},
		err:   false,
	},
	{
		b: new(SelectBuilder).
			Select([]string{"*"}).
			From("tbl").
			Join(RawSQL("JOIN tbl2 ON tbl2.id = tbl1.some_id")).
			Where(RawSQL("field1 = value")).
			Where(&RawSQLWithArgs{"field2 = ?", []interface{}{"value2"}}).
			GroupBy("field").
			Having(RawSQL("field = value")).
			OrderBy(RawSQL("created_at DESC")).
			Limit(10).
			Offset(20),
		query: "SELECT * FROM tbl JOIN tbl2 ON tbl2.id = tbl1.some_id WHERE field1 = value AND field2 = ? GROUP BY field HAVING field = value ORDER BY created_at DESC LIMIT 10 OFFSET 20",
		args:  []interface{}{"value2"},
		err:   false,
	},
	{
		b: new(SelectBuilder).
			Select([]string{"first_name", "last_name"}).
			Distinct().
			From("users").
			Where(&Eq{"id", 1}),
		query: "SELECT DISTINCT first_name, last_name FROM users WHERE id = ?",
		args:  []interface{}{1},
		err:   false,
	},
}

func TestSelectBuilder(t *testing.T) {
	for i, tt := range selectBuilderTests {
		meta := fmt.Sprintf("test %d", i)
		sql, args, err := tt.b.ToSQL()
		if err != nil && tt.err == false {
			t.Errorf("%s: SelectBuilder.ToSql() returned unexpected error: %s", meta, err)
			continue
		}
		if sql != tt.query {
			t.Errorf("%s: expected SelectBuilder.ToSql() to return %q query, got %q", meta, tt.query, sql)
			continue
		}
		if !reflect.DeepEqual(args, tt.args) {
			t.Errorf("%s: expected SelectBuilder.ToSql() to return args  %+v, got %+v", meta, tt.args, args)
			continue
		}
	}
}
