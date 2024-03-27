package sarama

import (
	"github.com/IBM/sarama"
	"testing"
)

var addr = []string{"localhost:9094"}

func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	// 设置等待服务器确认消息成功写入的时间设置生产者返回成功的消息
	cfg.Producer.Return.Successes = true
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	producer, err := sarama.NewSyncProducer(addr, cfg)
	if err != nil {
		panic("Failed to start Sarama producer: producer err" + err.Error())
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("这是一条消息"),
		// 生产者和消费者都可以使用Headers来传递额外的信息
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		t.Error(err)
	}
	t.Logf("partition:%d,offset:%d", partition, offset)
}
func TestAsyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	cfg.Producer.Return.Errors = true
	cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	producer, err := sarama.NewAsyncProducer(addr, cfg)
	if err != nil {
		panic("Failed to start Sarama producer: " + err.Error())
	}
	defer producer.Close()
	msg := producer.Input()
	msg <- &sarama.ProducerMessage{
		Topic: "test_topic",
		Value: sarama.StringEncoder("这是一条消息"),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("key1"),
				Value: []byte("value1"),
			},
		},
	}
	select {
	case err := <-producer.Errors():
		t.Error(err)
	case msg := <-producer.Successes():
		t.Logf("partition:%d,offset:%d", msg.Partition, msg.Offset)
		t.Logf("msg:%s", msg.Value.(sarama.StringEncoder))
	}
}
