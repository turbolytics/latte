package kafka

import (
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/turbolytics/latte/internal/sink"
)

type config struct {
	Brokers                []string
	Topic                  string
	AllowAutoTopicCreation bool `mapstructure:"allow_auto_topic_creation"`
}

type Kafka struct {
	config config
	writer *kafka.Writer
}

func (k *Kafka) Close() error {
	return k.writer.Close()
}

func (k *Kafka) Type() sink.Type {
	return sink.TypeKafka
}

func (k *Kafka) Write(bs []byte) (int, error) {
	err := k.writer.WriteMessages(context.TODO(),
		kafka.Message{
			Value: bs,
		},
	)
	return len(bs), err
}

func NewFromGenericConfig(m map[string]any) (*Kafka, error) {
	var conf config
	if err := mapstructure.Decode(m, &conf); err != nil {
		return nil, err
	}

	var w *kafka.Writer
	w = &kafka.Writer{
		Addr:                   kafka.TCP(conf.Brokers...),
		Topic:                  conf.Topic,
		AllowAutoTopicCreation: conf.AllowAutoTopicCreation,
	}

	return &Kafka{
		config: conf,
		writer: w,
	}, nil
}
