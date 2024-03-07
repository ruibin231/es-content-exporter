package es

import (
	"es-content-export/settings"
	"fmt"
	"github.com/olivere/elastic/v7"
)

func NewESClient(clientData *settings.EsClient) (client *elastic.Client, err error) {
	url := fmt.Sprintf("http://%s:%s", clientData.Host, clientData.Port)
	clientConf := []elastic.ClientOptionFunc{
		elastic.SetURL(url), elastic.SetSniff(false), elastic.SetGzip(true)}
	if clientData.User != "" && clientData.Pass != "" {
		clientConf = append(clientConf, elastic.SetBasicAuth(clientData.User, clientData.Pass))
	}
	client, err = elastic.NewClient(clientConf...)
	if err != nil {
		return nil, err
	}
	return
}
