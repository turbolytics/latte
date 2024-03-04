package kafka

import (
	"bytes"
	"context"
	"github.com/mitchellh/mapstructure"
	"github.com/segmentio/kafka-go"
	"github.com/turbolytics/latte/internal/encoding"
	"github.com/turbolytics/latte/internal/record"
	"github.com/turbolytics/latte/internal/sink"
)

type config struct {
	URI                    string
	Encoding               encoding.Config
	Topic                  string
	AllowAutoTopicCreation bool `mapstructure:"allow_auto_topic_creation"`
}

type Kafka struct {
	config config

	encoder encoding.Encoder
	writer  *kafka.Writer
}

func (k *Kafka) Close() error {
	return k.writer.Close()
}

func (k *Kafka) Flush(ctx context.Context) error {
	return nil
}

func (k *Kafka) Type() sink.Type {
	return sink.TypeKafka
}

func (k *Kafka) Write(ctx context.Context, r record.Record) (int, error) {
	buf := &bytes.Buffer{}
	if err := k.encoder.Init(buf); err != nil {
		return 0, nil
	}

	if err := k.encoder.Write(r.Map()); err != nil {
		return 0, err
	}

	bs := buf.Bytes()

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

	w := &kafka.Writer{
		Addr:                   kafka.TCP(conf.URI),
		Topic:                  conf.Topic,
		AllowAutoTopicCreation: conf.AllowAutoTopicCreation,
	}

	e, err := encoding.NewEncoder(conf.Encoding)

	if err != nil {
		return nil, err
	}

	return &Kafka{
		config:  conf,
		encoder: e,
		writer:  w,
	}, nil
}
