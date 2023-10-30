package es

import (
	"es-content-export/settings"
	"fmt"
	"github.com/olivere/elastic/v7"
)

func NewESClient() (client *elastic.Client, err error) {
	url := fmt.Sprintf("http://%s:%s", settings.Config.Host, settings.Config.Port)
	clientConf := []elastic.ClientOptionFunc{
		elastic.SetURL(url), elastic.SetSniff(false), elastic.SetGzip(true)}
	if settings.Config.User != "" && settings.Config.Pass != "" {
		clientConf = append(clientConf, elastic.SetBasicAuth(settings.Config.User, settings.Config.Pass))
	}
	client, err = elastic.NewClient(clientConf...)
	if err != nil {
		return nil, err
	}
	return
}
