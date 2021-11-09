package extension

import (
	"context"
	"fmt"
	"strconv"

	"github.com/velmie/qparser"

	"github.com/velmie/q2sql"
)

const Unlimited = int64(-1)

// LimitOffsetPagination is the extension
// that sets limit and offset based on the corresponding fields of the given query.Page
func LimitOffsetPagination(maxLimit, maxOffset int64) q2sql.Extension {
	return func(_ context.Context, query *qparser.Query, builder *q2sql.SelectBuilder) error {
		if page := query.Page; page != nil {
			if page.Limit != "" {
				limit, err := strconv.ParseUint(page.Limit, 10, 32)
				if err != nil {
					return fmt.Errorf("page limit must be unsigned integer, got %q", page.Limit)
				}
				if maxLimit != Unlimited && int64(limit) > maxLimit {
					return fmt.Errorf("page limit cannot be greater than %d", maxLimit)
				}
				builder.Limit(limit)
			}
			if page.Offset != "" && page.Limit == "" {
				return fmt.Errorf("offset cannot be used without specifying limit")
			}
			if page.Offset != "" {
				offset, err := strconv.ParseUint(page.Offset, 10, 32)
				if err != nil {
					return fmt.Errorf("page offset must be unsigned integer, got %q", page.Offset)
				}
				if maxOffset != Unlimited && int64(offset) > maxOffset {
					return fmt.Errorf("page offset cannot be greater than %d", maxOffset)
				}
				builder.Offset(offset)
			}
		}
		return nil
	}
}
