package q2sql

type ResourceSelectBuilderOption func(b *ResourceSelectBuilder)

// WithDefaultFields sets default fields which are used in the SELECT
// SQL statement in case if no specific fields are requested
func WithDefaultFields(fields []string) ResourceSelectBuilderOption {
	return func(b *ResourceSelectBuilder) {
		b.defaultFields = fields
	}
}

// AllowFiltering sets filtering rules
func AllowFiltering(
	allowedFilters AllowedConditions,
	conditions ConditionFactory,
	parser FilterExpressionParser,
) ResourceSelectBuilderOption {
	return func(b *ResourceSelectBuilder) {
		b.allowedConditions = allowedFilters
		b.conditions = conditions
		b.parser = parser
	}
}

// AllowSelectFields allows to "SELECT" given fields
func AllowSelectFields(fields []string) ResourceSelectBuilderOption {
	return func(b *ResourceSelectBuilder) {
		if b.allowedSelectFields == nil {
			b.allowedSelectFields = make(map[string]struct{})
		}
		fillMapKeys(b.allowedSelectFields, fields)
	}
}

// AllowSortingByFields allows to use specified fields in the "ORDER BY" SQL statement
func AllowSortingByFields(fields []string) ResourceSelectBuilderOption {
	return func(b *ResourceSelectBuilder) {
		b.allowedSortFields = fields
	}
}

// Extend adds Extensions to the list
func Extend(extensions ...Extension) ResourceSelectBuilderOption {
	return func(b *ResourceSelectBuilder) {
		b.extensions = append(b.extensions, extensions...)
	}
}
