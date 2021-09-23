package log

import (
	"github.com/lianmc123/app-frame/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestZapLogger(t *testing.T) {
	encoder := zapcore.EncoderConfig{
		TimeKey:        "t",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stack",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	logger := NewZapLogger(
		encoder,
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
		zap.AddStacktrace(
			zap.NewAtomicLevelAt(zapcore.ErrorLevel)),
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.Development(),
	)
	/*log.NewStdLogger(logger.)
	zlog := log.NewHelper(logger)
	zlog.Infow("name", "kratos", "from", "opensource")
	zlog.Infow("name", "kratos", "from")*/

	logger.Log(log.Debug, "msg", "1234567")
	logger.Log(log.Info, "abcdefs", "fucker")
	logger.Log(log.Warning, "abcdefs", "fucker")
	logger.Log(log.Error, "abcdefs", "fucker")
	//logger.Log(log.Fatal, "abcdefs", "fucker")

	zapHelper := log.NewHelper(logger)
	zapHelper.DebugKv("abc", 123123123)
	zapHelper.FatalKv("asdwer", 123123123)
	defer logger.Sync()
}

