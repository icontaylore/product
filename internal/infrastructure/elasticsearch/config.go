package elasticsearch

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"log"
)

type Config struct {
	Addresses []string
	Username  string
	Password  string
	CloudID   string
	APIKey    string
}

func NewClient(cfg Config) (*elasticsearch.Client, error) {
	esCfg := elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
		CloudID:   cfg.CloudID,
		APIKey:    cfg.APIKey,
	}
	client, err := elasticsearch.NewClient(esCfg)
	if err != nil {
		return nil, err
	}

	res, err := client.Ping()
	if err != nil {
		return nil, fmt.Errorf("elastic:failed to ping Elasticsearch: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("elastic:elasticsearch ping error: %s", res.String())
	}

	log.Println("elastic:successfully connected to Elasticsearch")
	return client, nil
}
