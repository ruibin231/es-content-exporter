package es

import (
	"context"
	"es-content-export/settings"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strings"
)

func formatIndex() string {
	if strings.HasSuffix(settings.Config.IndexPrefix, "*") {
		return settings.Config.IndexPrefix
	}
	return fmt.Sprintf("%s*", settings.Config.IndexPrefix)
}

func QueryLogCount() (int64, error) {
	client, err := NewESClient()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	index := elastic.NewWildcardQuery("_index", formatIndex())
	searchQuery := elastic.NewBoolQuery()
	searchQuery.Must(elastic.NewRangeQuery("@timestamp").Gte(
		fmt.Sprintf("now-%dm", settings.Config.Cycle)).Lte("now"))
	for _, val := range settings.Config.Must {
		searchQuery.Must(elastic.NewMatchQuery(settings.Config.Field, val))
	}
	if len(settings.Config.MustNot) > 0 {
		for _, noVal := range settings.Config.MustNot {
			searchQuery.MustNot(elastic.NewMatchQuery(settings.Config.Field, noVal))
		}
	}
	count, err := client.Count().
		Query(index).
		Query(searchQuery).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return count, nil
}
