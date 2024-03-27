package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/prometheus/common/log"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"testing"
	"time"
)

func TestCustomer(t *testing.T) {
	cfg := sarama.NewConfig()
	consumer, err := sarama.NewConsumerGroup(addr, "demo", cfg)
	assert.NoError(t, err)
	defer consumer.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 启动消费者
	//start := time.Now()
	err = consumer.Consume(ctx, []string{"test_topic"}, &ConsumerHandler{})
	assert.NoError(t, err)

	// 等待直到超时10s
	//t.Log(time.Since(start).String())
}

type ConsumerHandler struct {
}

func (c *ConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	log.Info("setup")
	partitions := session.Claims()["test_topic"]
	for _, part := range partitions {
		session.ResetOffset("test_topic", part, sarama.OffsetOldest, "")
	}
	return nil
}

func (c *ConsumerHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Info("cleanup")
	return nil
}

// ConsumeClaim 异步消费批量消费
// 这个是异步消费，批量提交的例子
func (c *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	ch := claim.Messages()
	batchSize := 10
	for {
		var eg errgroup.Group
		msgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				// 这一批次已经超时了，
				// 或者，整个 consumer 被关闭了
				// 不再尝试凑够一批了
				done = true
			case msg, ok := <-ch:
				if !ok {
					cancel()
					// channel 被关闭了
					return nil
				}
				msgs = append(msgs, msg)
				eg.Go(func() error {
					log.Infof("offset %d", msg.Offset)
					// 标记为消费成功
					time.Sleep(time.Second * 3)
					return nil
				})
			}
		}
		err := eg.Wait()
		if err == nil {
			// 这边就要都提交了
			for _, msg := range msgs {
				session.MarkMessage(msg, "")
			}
		} else {
			// 这里可以考虑重试，也可以在具体的业务逻辑里面重试
			// 也就是 eg.Go 里面重试
		}
		cancel()
	}
}

// ConsumeClaimV2 批量消费批量消费
func (c *ConsumerHandler) ConsumeClaimV2(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	const batchSize = 10
	batch := make([]*sarama.ConsumerMessage, 0, batchSize)
	for i := 0; i < 10; i++ {
		msg := <-msgs
		batch = append(batch, msg)
	}
	for _, msg := range batch {
		log.Infof("msg: %s", string(msg.Value))
	}
	for _, msg := range batch {
		session.MarkMessage(msg, "")
	}
	return nil
}

func (c *ConsumerHandler) ConsumeClaimV1(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		log.Infof("msg: %s", string(msg.Value))
		// 提交
		session.MarkMessage(msg, "")
	}
	return nil
}
