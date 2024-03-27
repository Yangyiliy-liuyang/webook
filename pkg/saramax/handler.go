package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"webook/pkg/logger"
)

type Handler[T any] struct {
	fn func(msg *sarama.ConsumerMessage, event T) error
	l  logger.ZapLogger
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		// 业务处理逻辑
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			// 记录错误日志,可以引入重试
			h.l.Error("json unmarshal error",
				logger.String("topic", msg.Topic),
				logger.Int32("partition", msg.Partition),
				logger.Int64("offset", msg.Offset),
				logger.Error(err))
		}
		// 消费完成
		session.MarkMessage(msg, "")
	}
	return nil
}

func NewHandler[T any](fn func(msg *sarama.ConsumerMessage, event T) error, l logger.ZapLogger) *Handler[T] {
	return &Handler[T]{
		fn: fn,
		l:  l,
	}
}
