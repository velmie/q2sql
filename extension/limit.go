package extension

import (
	"context"
	"github.com/velmie/q2sql"
	"github.com/velmie/qparser"
)

// DefaultLimit is the extension that sets the query limit if it has not been set
func DefaultLimit(limit uint64) q2sql.Extension {
	return func(_ context.Context, _ *qparser.Query, builder *q2sql.SelectBuilder) error {
		if builder.LimitPart == "" {
			builder.Limit(limit)
		}
		return nil
	}
}
