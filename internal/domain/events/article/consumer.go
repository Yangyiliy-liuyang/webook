package article

import (
	"context"
	"github.com/IBM/sarama"
	"time"
	"webook/internal/repository"
	"webook/pkg/logger"
	"webook/pkg/saramax"
)

type InteractiveReadEventConsumer struct {
	repo   repository.InteractiveRepository
	client sarama.Client
	l      logger.ZapLogger
}

func NewInteractiveReadEventConsumer(repo repository.InteractiveRepository, client sarama.Client, l logger.ZapLogger) *InteractiveReadEventConsumer {
	return &InteractiveReadEventConsumer{repo: repo, client: client, l: l}
}

func (i *InteractiveReadEventConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive", i.client)
	if err != nil {
		return err
	}
	go func() {
		er := cg.Consume(context.Background(),
			[]string{"interactive_read_event"},
			saramax.NewHandler[ReadEvent](i.Consumer, i.l),
		)
		if er != nil {
			i.l.Error("consumer error", logger.Error(er))
		}
	}()
	return nil
}

func (i *InteractiveReadEventConsumer) Consumer(msg *sarama.ConsumerMessage, event ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := i.repo.IncrReadCnt(ctx, "article", event.ArtId)
	return err
}
