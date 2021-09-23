package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var _ Logger = (*stdLogger)(nil)

type stdLogger struct {
	log  *log.Logger
	pool *sync.Pool
}

func NewStdLogger(writer io.Writer) Logger {
	return &stdLogger{
		log:  log.New(writer, "", 0),
		pool: &sync.Pool{New: func() interface{} { return new(bytes.Buffer)}},
	}
}

func (s *stdLogger) Log(level Level, kvPairs ...interface{}) error {
	if len(kvPairs)&1 == 1 {
		kvPairs = append(kvPairs, " [Pairs Error]")
	}
	buf := s.pool.Get().(*bytes.Buffer)
	_, err := buf.WriteString(level.String())
	if err != nil {
		return err
	}
	for i := 0; i < len(kvPairs); i += 2 {
		_, err := fmt.Fprintf(buf, " %s=%v", kvPairs[i], kvPairs[i+1])
		if err != nil {
			return err
		}
	}
	s.log.Output(4, buf.String())
	buf.Reset()
	s.pool.Put(buf)
	if level == Fatal {
		os.Exit(1)
	}
	return nil
}
