package logger

type AppLogger interface {
	LogInfof(format string, args ...interface{})
	LogInfo(message string)
	LogError(message string)
	LogErrorf(format string, args ...interface{})
}
