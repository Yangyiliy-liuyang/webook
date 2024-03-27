package article

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type Producer interface {
	ProduceReadEvent(event ReadEvent) error
}

type ReadEvent struct {
	ArtId int64 // 文章id
	Uid   int64 // 用户id
}

type SaramaSyncProducer struct {
	producer       sarama.SyncProducer
	TopicReadEvent string
}

func NewSaramaSyncProducer(producer sarama.SyncProducer) Producer {
	return &SaramaSyncProducer{producer: producer,
		TopicReadEvent: "article_read",
	}
}

func (s *SaramaSyncProducer) ProduceReadEvent(event ReadEvent) error {
	val, err := json.Marshal(event)

	msg := &sarama.ProducerMessage{
		Topic: s.TopicReadEvent, //"article_read"
		Value: sarama.ByteEncoder(val),
	}
	_, _, err = s.producer.SendMessage(msg)
	return err
}
