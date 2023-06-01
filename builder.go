package q2sql

import (
	"context"
	"fmt"

	"github.com/velmie/qparser"
)

type Builder interface {
	// Build builds sql query
	Build(ctx context.Context, query *qparser.Query, sb ...*SelectBuilder) (*SelectBuilder, error)
}

// Extension is used in order to extend a builder
type Extension func(ctx context.Context, query *qparser.Query, builder *SelectBuilder) error

// AllowedConditions maps a string field name to a list of condition aliases.
// It is used in order to specify a list of condition which could be applied to the field.
type AllowedConditions map[string][]string

// ResourceSelectBuilder is default implementation of sql select query builder
type ResourceSelectBuilder struct {
	resourceName           string
	defaultFields          []string
	allowedConditions      AllowedConditions
	allowedSelectFields    map[string]struct{}
	allowedSelectFieldsSlc []string
	allowedSortFields      []string
	translator             Translator
	parser                 FilterExpressionParser
	conditions             ConditionFactory
	extensions             []Extension
	alwaysSelectFields     []string
	alwaysSelectAllFields  bool
}

// NewResourceSelectBuilder is ResourceSelectBuilder constructor
func NewResourceSelectBuilder(
	resourceName string,
	translator Translator,
	options ...ResourceSelectBuilderOption,

) *ResourceSelectBuilder {
	b := &ResourceSelectBuilder{
		resourceName: resourceName,
		translator:   translator,
	}
	for _, option := range options {
		option(b)
	}
	if b.allowedSelectFields == nil {
		b.allowedSelectFields = make(map[string]struct{})
		fillMapKeys(b.allowedSelectFields, b.defaultFields)
		b.allowedSelectFieldsSlc = b.defaultFields
	}
	b.allowedSelectFieldsSlc = removeDuplicateStrings(b.allowedSelectFieldsSlc)
	return b
}

// Build builds sql query which depends on the applied options
func (s *ResourceSelectBuilder) Build(
	ctx context.Context,
	query *qparser.Query,
	sb ...*SelectBuilder,
) (*SelectBuilder, error) {
	var (
		selectFields []string
		b            *SelectBuilder
	)
	if len(sb) > 0 {
		b = sb[0]
	} else {
		b = new(SelectBuilder)
	}
	if s.alwaysSelectAllFields {
		selectFields = s.allowedSelectFieldsSlc
	} else {
		if fields, ok := query.Fields.FieldsByResource(s.resourceName); ok {
			f, err := s.translator(fields)
			if err != nil {
				return nil, err
			}
			selectFields = f
		} else {
			selectFields = s.defaultFields
		}
		selectFields = append(selectFields, s.alwaysSelectFields...)
		selectFields = removeDuplicateStrings(selectFields)
	}
	for _, field := range selectFields {
		if _, ok := s.allowedSelectFields[field]; !ok {
			return nil, fmt.Errorf("field %q not allowed for selection criteria", field)
		}
	}
	b.Select(selectFields).From(s.resourceName)
	conditions, err := s.retrieveFilterConditions(query)
	if err != nil {
		return nil, err
	}
	if len(conditions) > 0 {
		b.Where(conditions...)
	}
	sortList := make([]qparser.Sort, len(query.Sort))
	sortFields := make([]string, len(query.Sort))
	for i := 0; i < len(query.Sort); i++ {
		sortFields[i] = query.Sort[i].FieldName
	}
	sortFields, err = s.translator(sortFields)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(query.Sort); i++ {
		allowed := false
		for _, allowedName := range s.allowedSortFields {
			if sortFields[i] == allowedName {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, fmt.Errorf("field %q not allowed for sorting criteria", query.Sort[i].FieldName)
		}
		sortList[i] = query.Sort[i]
		sortList[i].FieldName = sortFields[i]
	}
	if len(sortList) > 0 {
		b.OrderBy(OrderBy(sortList))
	}
	for _, extension := range s.extensions {
		if err := extension(ctx, query, b); err != nil {
			return nil, err
		}
	}
	return b, nil
}

func (s *ResourceSelectBuilder) retrieveFilterConditions(query *qparser.Query) ([]Sqlizer, error) {
	conditions := make([]Sqlizer, 0)
	for _, filter := range query.Filters {
		allowList, ok := s.allowedConditions[filter.FieldName]
		if !ok {
			return nil, &FilterError{
				Field:   filter.FieldName,
				Message: fmt.Sprintf("filters cannot be applied to the field %q", filter.FieldName),
			}
		}
		f, err := s.translator([]string{filter.FieldName})
		if err != nil {
			return nil, err
		}
		name, args, err := s.parser.ParseFilterExpression(filter.Predicate)
		if err != nil {
			return nil, err
		}
		allowed := false
		for _, allowedName := range allowList {
			if name == allowedName {
				allowed = true
				break
			}
		}
		if !allowed {
			return nil, &FilterError{
				Filter:  name,
				Field:   filter.FieldName,
				Message: fmt.Sprintf("filter %q cannot be applied to the field %q", name, filter.FieldName),
			}
		}
		condition, err := s.conditions.CreateCondition(name)
		if err != nil {
			return nil, err
		}

		cond, err := condition(f[0], toInterfaceSlice(args)...)
		if err != nil {
			return nil, err
		}

		conditions = append(conditions, cond)
	}
	return conditions, nil
}

func toInterfaceSlice(s []string) []interface{} {
	dest := make([]interface{}, len(s))
	for i := 0; i < len(s); i++ {
		dest[i] = s[i]
	}
	return dest
}

func fillMapKeys(m map[string]struct{}, keys []string) {
	for _, key := range keys {
		m[key] = struct{}{}
	}
}

func removeDuplicateStrings(s []string) []string {
	allKeys := make(map[string]struct{})
	var list []string
	for _, item := range s {
		if _, value := allKeys[item]; !value {
			allKeys[item] = struct{}{}
			list = append(list, item)
		}
	}
	return list
}
