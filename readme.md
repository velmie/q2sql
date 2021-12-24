# q2sql - Convert URL query to SQL

It turns this `?fields[articles]=id,title,createdAt&filter[title]contains:bitcoin&sort=-createdAt`

into this `SELECT id, title, created_at FROM articles WHERE title LIKE "%bitcoin%" ORDER BY created_at DESC`

## Use cases

The package is not an ORM or SQL builder, but rather an add-on. The main task of the package is
simplification of creating an API for requesting a list of some resources.

It can be used in conjunction with any ORM that supports raw sql execution or with "native" golang sql database handle.

## Note

The package uses [github.com/velmie/qparser](https://github.com/velmie/qparser) to work with the request query.
See the QParser documentation for query string syntax and other details.

## Usage

### Initializing

The first step is to provide a "dictionary" to translate from URL string terms to SQL terms.
For example, a table field may have the name "created_at", but in the URL it will be indicated as "createdAt".
```go
	translator := q2sql.MapTranslator(map[string]string{
		"id":        "id",
		"title":     "title",
		"body":      "body",
		"author":    "author",
		"createdAt": "created_at",
	})
```


Nothing is allowed by default. Everything must be specified explicitly.
The second step is to specify default fields to use in the SELECT SQL statement.
```go
// the list will be used latter
defaultFields := []string{"id", "title", "created_at"}
```

Default fields are used in case if there is no specified fields in the URL query string.
The default fields will also be used as a list of allowed for selection fields. 
It is possible to allow additional fields for selection. It will be explained below.

We now have everything we need to create a query builder.

```go
	const resourceName = "articles"

	builder := q2sql.NewResourceSelectBuilder(
		resourceName,
		translator,
		q2sql.WithDefaultFields(defaultFields),
	)
```
 
The constructor takes 2 required arguments:
 
* resource name - the name of the database table
* translator - for translating terms as explained above

Everything else is options that expand the capabilities of the builder.

### The builder options

#### WithDefaultFields - sets the default field list

These fields are used if no fields are specified in the request query 
like this `?fields[articles]=id,name`.

#### AllowFiltering - enables filtering, specifies filtering rules

In order to use filtering, it should be specified: which filters can be applied to which fields, how to parse the filter query string
and how to create conditions based on the filter names. The package provides everything needed for this.

For example, let's consider the following task: it is required to implement filtering by the list of ids (id field).
```go
    //...
import (
	"github.com/velmie/q2sql"
	"github.com/velmie/condition"
)
    //...

    // 1. how to create conditions based on the filter names
	conditions := q2sql.ConditionMap{
		"any": condition.In,    
	}
    // 2. which filters can be applied to which fields
	allowedConditionsByField := q2sql.AllowedConditions{
		"id": []string{"any"},
	}
    // 3. how to parse the filter query string
    // pseudocode: "filter[id]=any:1,2,3" -> { Field: "id", Filter: "any", Args: ["1", "2", "3"] }
    parser := q2sql.NewDelimitedArgsParser(':', ',')

    // use with builder
	const resourceName = "articles"
	
	builder := q2sql.NewResourceSelectBuilder(
		resourceName,
		translator,
		q2sql.WithDefaultFields(defaultFields),
		q2sql.AllowFiltering(allowedConditionsByField, conditions, parser),
	)     
```

#### AllowSelectFields - adds a list of allowed fields to the selection

This option is used to explicitly specify which fields are allowed to be used in the build SELECT SQL statement.
The option could be applied multiple times.

```go
	builder := q2sql.NewResourceSelectBuilder(
		resourceName,
		translator,
		q2sql.WithDefaultFields(defaultFields),
		q2sql.AllowSelectFields(defaultFields),
        q2sql.AllowSelectFields([]string{"body", "author"}),
	)   
```

#### AllowSortingByFields - sets a list of allowed fields for sorting

This option allows you to specify a list of fields that can be used for sorting i.e. for building the ORDER BY statement.

```go
	builder := q2sql.NewResourceSelectBuilder(
		resourceName,
		translator,
		q2sql.WithDefaultFields(defaultFields),
        q2sql.AllowSortingByFields([]string{"created_id", "title"}),
	)   
```

#### Extend - this special option allows you to extend the functionality of the builder

For example, the builder does not implement the pagination functionality. Different projects may have their own requirements
to the implementation of pagination. In such a case, pagination is a good candidate for creating a builder extension.

The package includes a simple implementation of "LIMIT/OFFSET" pagination as an extension.  

```go
    //...
import (
	"github.com/velmie/q2sql"
	"github.com/velmie/extension"
)
    //...
    var maxLimit, maxOffset int64 = 10, 1000
	builder := q2sql.NewResourceSelectBuilder(
		resourceName,
		translator,
		q2sql.WithDefaultFields(defaultFields),
        q2sql.Extend(extension.LimitOffsetPagination(maxLimit, maxOffset)),

	)   
```

The Extension accesses * q2sql.SelectBuilder and can use it to modify the result query.

## Usage example

```go
	translator := q2sql.MapTranslator(map[string]string{
		"id":        "id",
		"title":     "title",
		"body":      "body",
		"author":    "author",
		"createdAt": "created_at",
	})

	allowedConditionsByField := q2sql.AllowedConditions{
		"id": []string{condition.NameIn},
	}
	builder := q2sql.NewResourceSelectBuilder(
		"articles",
		translator,
		q2sql.AllowSelectFields([]string{"id", "title", "author"}),
		q2sql.AllowFiltering(
			allowedConditionsByField,
			condition.DefaultConditionMap,
			q2sql.DefaultFilterExpressionParser,
		),
		q2sql.AllowSortingByFields([]string{"created_at", "title"}),
		q2sql.Extend(extension.LimitOffsetPagination(extension.Unlimited, extension.Unlimited)),
	)

	qstr := "?fields[articles]=id,title,author&filter[id]=in:1,2,3,4,5&sort=-createdAt,title&page[limit]=100"
	query, _ := qparser.ParseQuery(qstr)

	sqlizer, err := builder.Build(context.Background(), query)
	if err != nil {
		log.Fatal("failed to build query", err)
	}

	sqlStr, args, err := sqlizer.ToSql()
	if err != nil {
		log.Fatal("failed to build SQL query", err)
	}
	fmt.Println(sqlStr)
	fmt.Println(args)
	// prints
	/*
		SELECT id, title, author FROM articles WHERE id IN (?,?,?,?,?) ORDER BY created_at DESC, title ASC LIMIT 100
		[1 2 3 4 5]
	 */
```

In order to apply custom conditions you may pass select builder to the resource select builder.

```go
    sb := new(q2sql.SelectBuilder)
	sb.Where(&q2sql.Eq{"author", "Alan Turing"})
    
    _, err := builder.Build(context.Background(), query, sb)
    ...
    
```
