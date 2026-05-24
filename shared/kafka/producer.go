package kafka

import (
	"context"

	"github.com/IBM/sarama"
)

type Producer struct {
	sync sarama.SyncProducer
}

func NewProducer(cfg Config) (*Producer, error) {
	scfg := sarama.NewConfig()
	scfg.ClientID = cfg.ClientID

	scfg.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(cfg.Brokers, scfg)
	if err != nil {
		return nil, err
	}

	return &Producer{sync: p}, nil
}

func (p *Producer) Send(
	ctx context.Context,
	topic string,
	msg *Message,
) error {
	m := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.ByteEncoder(msg.Key),
		Value: sarama.ByteEncoder(msg.Value),
	}

	for k, v := range msg.Headers {
		m.Headers = append(m.Headers, sarama.RecordHeader{
			Key:   []byte(k),
			Value: []byte(v),
		})
	}

	_, _, err := p.sync.SendMessage(m)
	return err
}

func (p *Producer) Close() error {
	return p.sync.Close()
}
