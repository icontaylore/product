package kafka

import "gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"

type Config struct {
	BootstrapServers string
	GroupID          string
	AutoOffsetReset  string
}

func NewConfig(brokers, groupID string) *Config {
	return &Config{
		BootstrapServers: brokers,
		GroupID:          groupID,
		AutoOffsetReset:  "earliest",
	}
}

func (c *Config) CreateConsumerConfig() *kafka.ConfigMap {
	return &kafka.ConfigMap{
		"bootstrap.servers":  c.BootstrapServers,
		"group.id":           c.GroupID,
		"auto.offset.reset":  c.AutoOffsetReset,
		"enable.auto.commit": false,
	}
}
