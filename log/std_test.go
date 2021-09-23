package log

import (
	"os"
	"testing"
)

func TestStdLogger(t *testing.T) {
	logger := NewStdLogger(os.Stdout)
	logger.Log(Debug, "a", 123)
	logger.Log(Info, "a", 123, "b", "asd")
	logger.Log(Warning, "a", 123, "c")
	logger.Log(Error, "a", 123)
	logger.Log(Fatal, "a", 123)

}
