package startup

import "GinStart/pkg/logger"

func InitLogger() logger.Logger {
	return logger.NewNopLogger()
}
