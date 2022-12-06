package q2sql

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/velmie/qparser"
)

type resourceBuilderTest struct {
	title             string
	query             string
	sql               string
	args              []interface{}
	expectedErr       bool
	sb                func() *SelectBuilder
	additionalOptions []ResourceSelectBuilderOption
}

const (
	resourceName           = "articles"
	resourceFieldID        = "id"
	resourceFieldTitle     = "title"
	resourceFieldBody      = "body"
	resourceFieldCreatedAt = "created_at"
	resourceFieldUpdatedAt = "updated_at"
	filterEq               = "eq"
	filterAny              = "any"
	filterContains         = "contains"
)

var resourceBuilderTests = []resourceBuilderTest{
	{
		title: "Default fields, no conditions",
		query: "",
		sql:   fmt.Sprintf("SELECT * FROM %s", resourceName),
		args:  []interface{}{},
	},
	{
		title: "Specific fields, no conditions",
		query: fmt.Sprintf("fields[%s]=%s,%s", resourceName, resourceFieldID, resourceFieldTitle),
		sql:   fmt.Sprintf("SELECT %s, %s FROM %s", resourceFieldID, resourceFieldTitle, resourceName),
		args:  []interface{}{},
	},
	{
		title: "Filter by id 42",
		query: fmt.Sprintf("filter[%s]=%s:42", resourceFieldID, filterEq),
		sql:   fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", resourceName, resourceFieldID),
		args:  []interface{}{"42"},
	},
	{
		title: "Filter by any id 1,2,3,4,5 and body must contain word 'bitcoin'",
		query: fmt.Sprintf("filter[%s]=%s:1,2,3,4,5&filter[%s]=%s:bitcoin", resourceFieldID, filterAny, resourceFieldBody, filterContains),
		sql:   fmt.Sprintf("SELECT * FROM %s WHERE %s IN (?,?,?,?,?) AND %s LIKE ?", resourceName, resourceFieldID, resourceFieldBody),
		args:  []interface{}{"1", "2", "3", "4", "5", "%bitcoin%"},
	},
	{
		title: "Apply sorting to the result set where title is 'NewYear'",
		query: fmt.Sprintf("sort=-createdAt&filter[%s]=%s:NewYear", resourceFieldTitle, filterEq),
		sql:   fmt.Sprintf("SELECT * FROM %s WHERE %s = ? ORDER BY %s DESC", resourceName, resourceFieldTitle, resourceFieldCreatedAt),
		args:  []interface{}{"NewYear"},
	},
	{
		title:       "Undefined field",
		query:       fmt.Sprintf("fields[%s]=unknown", resourceName),
		expectedErr: true,
	},
	{
		title:       "Undefined filter",
		query:       fmt.Sprintf("filter[%s]=unknown:value", resourceFieldID),
		expectedErr: true,
	},
	{
		title:       "Not allowed filter",
		query:       fmt.Sprintf("filter[%s]=%s:value", resourceFieldCreatedAt, filterEq),
		expectedErr: true,
	},
	{
		title:       "Not allowed sorting",
		query:       "sort=id",
		expectedErr: true,
	},
	{
		title:       "Field updated_at is not allowed",
		query:       fmt.Sprintf("fields[%s]=updatedAt", resourceName),
		expectedErr: true,
	},
	{
		title: "With extra params",
		query: fmt.Sprintf("fields[%s]=%s,%s", resourceName, resourceFieldID, resourceFieldTitle),
		sql:   fmt.Sprintf("SELECT %s, %s FROM %s WHERE id = ?", resourceFieldID, resourceFieldTitle, resourceName),
		args:  []interface{}{"1"},
		sb: func() *SelectBuilder {
			sb := new(SelectBuilder)
			return sb.Where(&RawSQLWithArgs{"id = ?", []interface{}{"1"}})
		},
	},
	{
		title: "Specific fields, specific always select field",
		query: fmt.Sprintf("fields[%s]=%s,%s", resourceName, resourceFieldID, resourceFieldTitle),
		sql:   fmt.Sprintf("SELECT %s, %s, %s FROM %s", resourceFieldID, resourceFieldTitle, resourceFieldBody, resourceName),
		args:  []interface{}{},
		additionalOptions: []ResourceSelectBuilderOption{
			AlwaysSelectFields([]string{resourceFieldBody}),
		},
	},
	{
		title: "Specific fields, always select all fields",
		query: fmt.Sprintf("fields[%s]=%s,%s", resourceName, resourceFieldID, resourceFieldTitle),
		sql:   fmt.Sprintf("SELECT *, %s, %s, %s, %s FROM %s", resourceFieldID, resourceFieldTitle, resourceFieldBody, resourceFieldCreatedAt, resourceName),
		args:  []interface{}{},
		additionalOptions: []ResourceSelectBuilderOption{
			AlwaysSelectAllFields(true),
		},
	},
	{
		title: "Second AllowSelectFields specification with duplicate *, always select all fields",
		query: fmt.Sprintf("fields[%s]=%s,%s", resourceName, resourceFieldID, resourceFieldTitle),
		sql:   fmt.Sprintf("SELECT *, %s, %s, %s, %s, %s FROM %s", resourceFieldID, resourceFieldTitle, resourceFieldBody, resourceFieldCreatedAt, resourceFieldUpdatedAt, resourceName),
		args:  []interface{}{},
		additionalOptions: []ResourceSelectBuilderOption{
			AllowSelectFields(
				[]string{
					"*",
					resourceFieldUpdatedAt,
				},
			),
			AlwaysSelectAllFields(true),
		},
	},
}

func TestNewResourceSelectBuilder(t *testing.T) {
	conditionFactory := ConditionMap{
		filterEq: func(field string, args ...interface{}) Sqlizer {
			return &Eq{
				Field: field,
				Value: args[0],
			}
		},
		filterAny: func(field string, args ...interface{}) Sqlizer {
			return &In{
				Field:  field,
				Values: args,
			}
		},
		filterContains: func(field string, args ...interface{}) Sqlizer {
			if len(args) == 0 {
				return RawSQL("")
			}
			val := "%" + args[0].(string) + "%"
			return &Like{
				Field: field,
				Value: val,
			}
		},
	}
	translator := MapTranslator(
		map[string]string{
			"id":        resourceFieldID,
			"title":     resourceFieldTitle,
			"body":      resourceFieldBody,
			"createdAt": resourceFieldCreatedAt,
			"updatedAt": resourceFieldUpdatedAt,
		},
	)
	allowedConditions := AllowedConditions{
		resourceFieldID:    []string{filterEq, filterAny},
		resourceFieldTitle: []string{filterEq, filterContains},
		resourceFieldBody:  []string{filterContains},
	}
	allowedSorting := []string{resourceFieldCreatedAt}
	parser := NewDelimitedArgsParser(':', ',')

	ctx := context.Background()
	for i, tt := range resourceBuilderTests {
		builder := NewResourceSelectBuilder(
			resourceName,
			translator,
			append(
				[]ResourceSelectBuilderOption{
					AllowFiltering(allowedConditions, conditionFactory, parser),
					AllowSortingByFields(allowedSorting),
					WithDefaultFields([]string{"*"}),
					AllowSelectFields(
						[]string{
							"*",
							resourceFieldID,
							resourceFieldTitle,
							resourceFieldBody,
							resourceFieldCreatedAt,
						},
					),
				},
				tt.additionalOptions...,
			)...,
		)
		meta := fmt.Sprintf("test %d\n%s:\n\n\t", i, tt.title)
		query, err := qparser.ParseQuery(tt.query)
		if err != nil {
			t.Fatalf("%sqparser.ParseQuery(%q) returned unexpected error: %s", meta, tt.query, err)
		}
		var sb []*SelectBuilder
		if tt.sb != nil {
			sb = []*SelectBuilder{tt.sb()}
		}
		sqlizer, err := builder.Build(ctx, query, sb...)
		if !tt.expectedErr && err != nil {
			t.Errorf("%sunexpected error %s", meta, err)
			continue
		} else if err != nil {
			continue
		}
		sql, args, err := sqlizer.ToSQL()
		if !tt.expectedErr && err != nil {
			t.Errorf("%sunexpected error %s", meta, err)
			continue
		} else if err != nil {
			continue
		}
		if sql != tt.sql {
			t.Errorf("%sexpected sql %q\n\tgot %q", meta, tt.sql, sql)
			continue
		}
		if !reflect.DeepEqual(args, tt.args) {
			t.Errorf("%sexpected args  %+v\n\tgot %+v", meta, tt.args, args)
			continue
		}
	}
}
