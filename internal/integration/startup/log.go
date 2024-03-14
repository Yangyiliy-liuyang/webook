package startup

import "webook/pkg/logger"

func InitLog() logger.Logger {
	return logger.NewNopLogger()
}
