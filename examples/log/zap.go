package log

import (
	"fmt"
	"github.com/lianmc123/app-frame/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var _ log.Logger = (*ZapLogger)(nil)

type ZapLogger struct {
	log  *zap.Logger
	Sync func() error
}

func NewZapLogger(encoder zapcore.EncoderConfig, level zap.AtomicLevel, opts ...zap.Option) *ZapLogger {
	core := zapcore.NewCore(
		//zapcore.NewConsoleEncoder(encoder),
		zapcore.NewJSONEncoder(encoder),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
		), level)
	zapLogger := zap.New(core, opts...)
	return &ZapLogger{log: zapLogger, Sync: zapLogger.Sync}
}

func (z *ZapLogger) Log(level log.Level, kvPairs ...interface{}) error {
	if len(kvPairs) == 0 || len(kvPairs)%2 != 0 {
		z.log.Warn(fmt.Sprint("kvPairs must appear in pairs: ", kvPairs))
		return nil
	}
	var data []zap.Field
	for i := 0; i < len(kvPairs); i += 2 {
		data = append(data, zap.Any(fmt.Sprint(kvPairs[i]), fmt.Sprint(kvPairs[i+1])))
	}
	switch level {
	case log.Debug:
		z.log.Debug("", data...)
	case log.Info:
		z.log.Info("", data...)
	case log.Warning:
		z.log.Warn("", data...)
	case log.Error:
		z.log.Error("", data...)
	case log.Fatal:
		z.log.Fatal("", data...)
	}
	return nil
}
