package elasticsearch

import (
	"context"
	"log"
	"strings"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func (c *ElasticService) CreateIndex(indexName string) error {
	res,err := c.Client.Indices.Exists([]string{indexName})
	if err != nil {
		return err
	}
	res.Body.Close()

	if res.StatusCode == 404 {
		mapping := `{
			"settings": {
				"number_of_shards": 1,
				"number_of_replicas": 1
			},
			"mappings": {
				"properties": {
					"title":    { "type": "text" },
					"content": { "type": "text" },
					"created": { "type": "date" }
				}
			}
		}`
		c.Mapping = mapping

		req := esapi.IndicesCreateRequest{
			Index: indexName,
			Body: strings.NewReader(mapping),
		}

		res,err := req.Do(context.Background(),c.Client)
		if err != nil {
			log.Println("elastic:err creating index")
			return err
		}
		defer res.Body.Close()

		if res.IsError() {
			log.Println("elastic:err response")
		} else {
			log.Println("elastic:index created succesfull")
		}
	} else {
		log.Println("elastic:index already exists")
	}

	return nil
}
