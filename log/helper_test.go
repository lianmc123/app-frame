package log

import (
	"log"
	"testing"
)

func TestNewHelper(t *testing.T) {
	helper := NewHelper(NewStdLogger(log.Writer()))
	helper.DebugKv("msg", "123456")
	helper.InfoKv("msg", "123456")
	helper.WarningKv("msg", "123456")
	helper.ErrorKv("msg", "123456")
	//helper.FatalKv("msg", "123456")

	helper.DebugF("%s %d", "123456", 11223344)
}
