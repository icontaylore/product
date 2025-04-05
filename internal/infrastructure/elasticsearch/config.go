package elasticsearch

import (
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

type Config struct {
	Addresses []string
	Username  string
	Password  string
	CloudID   string
	APIKey    string
}

type ElasticService struct {
	Client *elasticsearch.Client
	Mapping string
}

func NewClientElastic(add []string) *ElasticService {
	cfg := elasticsearch.Config{
		Addresses: add,
	}

	client,err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Println("elastic:ошибка создания клиента")
	}
	res,err := client.Ping()
	if err != nil {
		log.Println("elastic:нет пинга")
	}
	if res.IsError() {
		fmt.Errorf("ping error: %s", res.String())
        return nil
    }
	return &ElasticService{
		Client: client,
	}
}

