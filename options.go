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
		b.allowedSelectFieldsSlc = append(b.allowedSelectFieldsSlc, fields...)
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

// AlwaysSelectFields sets fields which that will always be included  in the SELECT
// SQL statement regardless if specific fields are requested or not
func AlwaysSelectFields(fields []string) ResourceSelectBuilderOption {
	return func(b *ResourceSelectBuilder) {
		b.alwaysSelectFields = fields
	}
}

// AlwaysSelectAllFields will always include all allowed fields in the SELECT
// SQL statement regardless if specific fields are requested or not
// Overrides AlwaysSelectFields
func AlwaysSelectAllFields(flag bool) ResourceSelectBuilderOption {
	return func(b *ResourceSelectBuilder) {
		b.alwaysSelectAllFields = flag
	}
}
