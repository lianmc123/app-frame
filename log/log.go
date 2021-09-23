package log

import "log"

var DefaultLogger = NewStdLogger(log.Writer())

type Logger interface {
	Log(level Level, kvPairs ...interface{}) error
}
