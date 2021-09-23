package log

import "fmt"

//Helper 可以直接使用Help.[LEVEL]方法打印日志
type Helper struct {
	logger Logger
}

func NewHelper(logger Logger) *Helper {
	return &Helper{
		logger: logger,
	}
}

func (h *Helper) Log(level Level, kvParis ...interface{}) {
	_ = h.logger.Log(level, kvParis...)
}

func (h *Helper) Debug(msg interface{}) {
	h.Log(Debug, "msg", msg)
}
func (h *Helper) Info(msg interface{}) {
	h.Log(Info, "msg", msg)
}
func (h *Helper) Warning(msg interface{}) {
	h.Log(Warning, "msg", msg)
}
func (h *Helper) Error(msg interface{}) {
	h.Log(Error, "msg", msg)
}
func (h *Helper) Fatal(msg interface{}) {
	h.Log(Fatal, "msg", msg)
}

func (h *Helper) DebugKv(kvParis ...interface{}) {
	h.Log(Debug, kvParis...)
}
func (h *Helper) InfoKv(kvParis ...interface{}) {
	h.Log(Info, kvParis...)
}
func (h *Helper) WarningKv(kvParis ...interface{}) {
	h.Log(Warning, kvParis...)
}
func (h *Helper) ErrorKv(kvParis ...interface{}) {
	h.Log(Error, kvParis...)
}
func (h *Helper) FatalKv(kvParis ...interface{}) {
	h.Log(Fatal, kvParis...)
}

func (h *Helper) DebugF(format string, values ...interface{}) {
	h.Log(Debug, "msg", fmt.Sprintf(format, values...))
}
func (h *Helper) InfoF(format string, values ...interface{}) {
	h.Log(Info, "msg", fmt.Sprintf(format, values...))
}
func (h *Helper) WarningF(format string, values ...interface{}) {
	h.Log(Warning, "msg", fmt.Sprintf(format, values...))
}
func (h *Helper) ErrorF(format string, values ...interface{}) {
	h.Log(Error, "msg", fmt.Sprintf(format, values...))
}
func (h *Helper) FatalF(format string, values ...interface{}) {
	h.Log(Fatal, "msg", fmt.Sprintf(format, values...))
}
