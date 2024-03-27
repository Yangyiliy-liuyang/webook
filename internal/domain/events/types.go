package events

// 领域事件

type Consumer interface {
	Start() error
}
