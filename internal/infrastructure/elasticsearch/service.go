package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
)

type Service struct {
	client    *elasticsearch.Client
	IndexName string
}

func NewService(client *elasticsearch.Client, indexName string) *Service {
	s := &Service{
		client:    client,
		IndexName: indexName,
	}
	s.client.Indices.Create(indexName)

	return s
}

func (s *Service) IndexDocument(ctx context.Context, doc interface{}) error {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(doc); err != nil {
		return fmt.Errorf("elastic:elastic-service:error encoding document: %w", err)
	}
	res, err := s.client.Index(
		s.IndexName,
		&buf,
		s.client.Index.WithContext(ctx),
		s.client.Index.WithRefresh("wait_for"),
	)
	if err != nil {
		return fmt.Errorf("elastic:elastic-service:error indexing document: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("elastic:elastic-service:error response: %s", res.String())
	}

	return nil
}
