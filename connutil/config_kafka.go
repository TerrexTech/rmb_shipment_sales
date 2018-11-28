package connutil

import (
	"fmt"
	"os"

	"github.com/TerrexTech/go-commonutils/commonutil"
	"github.com/TerrexTech/go-eventspoll/poll"
	"github.com/TerrexTech/go-kafkautils/kafka"
)

// LoadKafkaConfig is a convenient function to load Kafkaconfig for EventsPoll.
func LoadKafkaConfig(aggID int8) (*poll.KafkaConfig, error) {
	kafkaBrokers := *commonutil.ParseHosts(
		os.Getenv("KAFKA_BROKERS"),
	)

	cEventGroup := os.Getenv("KAFKA_CONSUMER_EVENT_GROUP")
	cEventQueryGroup := os.Getenv("KAFKA_CONSUMER_EVENT_QUERY_GROUP")
	cEventTopic := os.Getenv("KAFKA_CONSUMER_EVENT_TOPIC")
	cEventQueryTopic := os.Getenv("KAFKA_CONSUMER_EVENT_QUERY_TOPIC")
	pEventQueryTopic := os.Getenv("KAFKA_PRODUCER_EVENT_QUERY_TOPIC")

	cEventTopic = fmt.Sprintf("%s.%d", cEventTopic, aggID)
	cEventQueryTopic = fmt.Sprintf("%s.%d", cEventQueryTopic, aggID)

	kc := &poll.KafkaConfig{
		EventCons: &kafka.ConsumerConfig{
			KafkaBrokers: kafkaBrokers,
			GroupName:    cEventGroup,
			Topics:       []string{cEventTopic},
		},
		ESQueryResCons: &kafka.ConsumerConfig{
			KafkaBrokers: kafkaBrokers,
			GroupName:    cEventQueryGroup,
			Topics:       []string{cEventQueryTopic},
		},

		ESQueryReqProd: &kafka.ProducerConfig{
			KafkaBrokers: kafkaBrokers,
		},
		ESQueryReqTopic: pEventQueryTopic,
	}

	return kc, nil
}
