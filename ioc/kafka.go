package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
)

func InitSaramaClient() sarama.Client {
	type Config struct {
		Addrs []string
	}
	var cfg Config
	err := viper.UnmarshalKey("kafkakafka", &cfg)
	if err != nil {
		panic(err)
	}
	scfg := sarama.NewConfig()
	scfg.Producer.Return.Successes = true
	// 创建Sarama客户端
	client, err := sarama.NewClient(cfg.Addrs, scfg)
	if err != nil {
		panic(err)
	}
	return client
}

func InitSyncProducer(c sarama.Client) sarama.SyncProducer {
	p, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return p
}

//
//func InitConsumers(c article.InteractiveReadEventConsumer) []events.Consumer {
//	return []events.Consumer{c}
//}
