package config

import (
	"strings"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

// NewKafkaConsumerGroup creates a new Kafka consumer group with the given configuration
func NewKafkaConsumerGroup(config *Config, log *logrus.Logger) sarama.ConsumerGroup {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true

	// Set offset reset behavior
	if config.Kafka.AutoOffsetReset == "earliest" {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	} else {
		saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	}

	// Parse brokers from comma-separated list
	brokers := strings.Split(config.Kafka.BootstrapServers, ",")
	groupID := config.Kafka.GroupID

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, saramaConfig)
	if err != nil {
		log.Fatalf("failed to create consumer group: %v", err)
	}

	log.Infof("Kafka consumer group created with brokers: %v, group: %s", brokers, groupID)
	return consumerGroup
}

// NewKafkaProducer creates a new Kafka producer with the given configuration
func NewKafkaProducer(config *Config, log *logrus.Logger) sarama.SyncProducer {
	if !config.Kafka.ProducerEnabled {
		log.Info("Kafka producer is disabled")
		return nil
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.RequiredAcks = sarama.WaitForAll
	saramaConfig.Producer.Retry.Max = 3

	// Parse brokers from comma-separated list
	brokers := strings.Split(config.Kafka.BootstrapServers, ",")

	producer, err := sarama.NewSyncProducer(brokers, saramaConfig)
	if err != nil {
		log.Fatalf("failed to create producer: %v", err)
	}

	log.Infof("Kafka producer created with brokers: %v", brokers)
	return producer
}
