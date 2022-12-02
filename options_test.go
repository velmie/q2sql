package q2sql

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/velmie/qparser"
)

type optionsTest struct {
	b       *ResourceSelectBuilder
	options []ResourceSelectBuilderOption
}

func mockExtension(_ context.Context, query *qparser.Query, builder *SelectBuilder) error { return nil }

var optionsTests = []optionsTest{
	{
		b: &ResourceSelectBuilder{
			defaultFields: []string{"fieldA", "fieldB"},
		},
		options: []ResourceSelectBuilderOption{
			WithDefaultFields([]string{"fieldA", "fieldB"}),
		},
	},
	{
		b: &ResourceSelectBuilder{
			allowedSortFields: []string{"fieldA", "fieldB"},
		},
		options: []ResourceSelectBuilderOption{
			AllowSortingByFields([]string{"fieldA", "fieldB"}),
		},
	},
	{
		b: &ResourceSelectBuilder{
			allowedSelectFields: map[string]struct{}{
				"x1": {},
				"x2": {},
				"x3": {},
				"y1": {},
				"y2": {},
				"y3": {},
			},
			allowedSelectFieldsSlc: []string{"x1", "x2", "x3", "y1", "y2", "y3"},
		},
		options: []ResourceSelectBuilderOption{
			AllowSelectFields([]string{"x1", "x2", "x3"}),
			AllowSelectFields([]string{"y1", "y2", "y3"}),
		},
	},
	{
		b: &ResourceSelectBuilder{
			allowedSelectFields: map[string]struct{}{
				"x1": {},
				"x2": {},
				"x3": {},
			},
			allowedSelectFieldsSlc: []string{"x1", "x2", "x3"},
			alwaysSelectFields:     []string{"x1"},
		},
		options: []ResourceSelectBuilderOption{
			AllowSelectFields([]string{"x1", "x2", "x3"}),
			AlwaysSelectFields([]string{"x1"}),
		},
	},
	{
		b: &ResourceSelectBuilder{
			allowedSelectFields: map[string]struct{}{
				"x1": {},
				"x2": {},
				"x3": {},
			},
			allowedSelectFieldsSlc: []string{"x1", "x2", "x3"},
			alwaysSelectAllFields:  true,
		},
		options: []ResourceSelectBuilderOption{
			AllowSelectFields([]string{"x1", "x2", "x3"}),
			AlwaysSelectAllFields(true),
		},
	},
}

func TestOptions(t *testing.T) {
	for i, tt := range optionsTests {
		meta := fmt.Sprintf("test %d", i)
		b := &ResourceSelectBuilder{}
		for _, option := range tt.options {
			option(b)
		}
		if !reflect.DeepEqual(b, tt.b) {
			t.Errorf("%s: \n\twant %+v \n\tgot %+v", meta, tt.b, b)
			continue
		}
	}
}

func TestExtend(t *testing.T) {
	b := new(ResourceSelectBuilder)
	Extend(mockExtension)(b)
	if len(b.extensions) != 1 {
		t.Errorf("unexpected extensions length: want 1, got %d", len(b.extensions))
		return
	}
	if fmt.Sprintf("%p", mockExtension) != fmt.Sprintf("%p", b.extensions[0]) {
		t.Error("unexpected extension in the list")
	}
}

type mockExprParser struct {
	FilterExpressionParser
}

func TestAllowFiltering(t *testing.T) {
	b := new(ResourceSelectBuilder)
	const conditionName = "c1"
	conditions := ConditionMap{
		conditionName: func(field string, args ...interface{}) Sqlizer {
			return nil
		},
	}
	allowedCondition := AllowedConditions{
		"field": []string{conditionName},
	}
	exprParser := new(mockExprParser)
	AllowFiltering(allowedCondition, conditions, exprParser)(b)
	if fmt.Sprintf("%p", b.conditions) != fmt.Sprintf("%p", conditions) {
		t.Errorf("unexpected condition factory:\n\twant %+v\n\tgot %+v", conditions, b.conditions)
	}
	if !reflect.DeepEqual(b.allowedConditions, allowedCondition) {
		t.Errorf("unexpected allowed conditions: \n\twant %+v \n\tgot %+v", allowedCondition, b.allowedConditions)
		return
	}
	if b.parser != exprParser {
		t.Errorf("unexpected expression parser:\n\twant %+v\n\tgot %+v", exprParser, b.parser)
	}
}
