package main

import (
	"context"
	"fmt"
	"log"

	"github.com/velmie/qparser"

	"github.com/velmie/q2sql"
	"github.com/velmie/q2sql/condition"
	"github.com/velmie/q2sql/extension"
)

func main() {
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

	// It's possible to pass select builder
	sb := new(q2sql.SelectBuilder)
	// and define custom conditions
	sb.Where(&q2sql.Eq{"author", "Alan Turing"})

	// not necessary to declare new variable since the "sb" already points to the select builder
	_, err := builder.Build(context.Background(), query, sb)
	if err != nil {
		log.Fatal("failed to build query", err)
	}
	// or rewrite some parts, for example the line below clears sorting
	sb.OrderByParts = nil

	sqlStr, args, err := sb.ToSql()
	if err != nil {
		log.Fatal("failed to build SQL query", err)
	}
	fmt.Println(sqlStr)
	fmt.Println(args)
	// prints
	/*
		SELECT id, title, author FROM articles WHERE author = ? AND id IN (?,?,?,?,?) LIMIT 100
		[Alan Turing 1 2 3 4 5]
	*/
}
