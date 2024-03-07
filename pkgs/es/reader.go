package es

import (
	"context"
	"es-content-export/settings"
	"fmt"
	"github.com/olivere/elastic/v7"
	"strings"
)

func formatIndex(index string) string {
	if strings.HasSuffix(index, "*") {
		return index
	}
	return fmt.Sprintf("%s*", index)
}

func QueryLogCount(esClient *settings.EsClient, queryData *settings.QueryData) (int64, error) {
	client, err := NewESClient(esClient)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	index := elastic.NewWildcardQuery("_index", formatIndex(queryData.IndexPrefix))
	searchQuery := elastic.NewBoolQuery()
	searchQuery.Must(elastic.NewRangeQuery("@timestamp").Gte(
		fmt.Sprintf("now-%dm", queryData.Cycle)).Lte("now"))
	searchQuery.Must(elastic.NewMatchPhraseQuery(queryData.Field, queryData.Content))
	//for _, val := range settings.Config.Must {
	//	searchQuery.Must(elastic.NewMatchQuery(settings.Config.Field, val))
	//}
	//if len(settings.Config.MustNot) > 0 {
	//	for _, noVal := range settings.Config.MustNot {
	//		searchQuery.MustNot(elastic.NewMatchQuery(settings.Config.Field, noVal))
	//	}
	//}
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
