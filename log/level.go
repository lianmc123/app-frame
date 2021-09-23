package log

type Level uint8

const (
	Debug Level = iota
	Info
	Warning
	Error
	Fatal
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Error:
		return "ERROR"
	case Warning:
		return "WARNING"
	case Fatal:
		return "FATAL"
	default:
		return "[-UNKNOWN-]"
	}
}
