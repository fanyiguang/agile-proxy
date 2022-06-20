package log

import (
	"agile-proxy/pkg/beego-log/core/logs"
	"encoding/json"
)

var (
	log *logs.BeeLogger
)

type config struct {
	Filename string `json:"filename"`
	Level    int    `json:"level"`
	MaxLines int    `json:"maxlines"`
	MaxSize  int    `json:"maxsize"`
	MaxDays  int    `json:"maxdays"`
	Daily    bool   `json:"daily"`
}

func init() {
	log = logs.NewLogger()
}

func New(logPath string, level string) {
	_config := config{
		Filename: logPath,
		Level:    getLogsLevel(level),
		MaxSize:  1 << 28,
		MaxDays:  5,
		Daily:    true,
	}
	jsonConfig, _ := json.Marshal(_config)
	_ = log.SetLogger(logs.AdapterFile, string(jsonConfig))
	log.EnableFuncCallDepth(true)
}

func getLogsLevel(level string) (l int) {
	switch level {
	case "debug":
		l = logs.LevelDebug
	case "warning":
		l = logs.LevelWarning
	case "warn":
		l = logs.LevelWarn
	case "info":
		l = logs.LevelInfo
	default:
		l = logs.LevelDebug
	}
	return
}

func WarnF(format string, args ...interface{}) {
	log.Warn(format, args...)
}

func InfoF(format string, args ...interface{}) {
	log.Info(format, args...)
}

func DebugF(format string, args ...interface{}) {
	log.Debug(format, args...)
}

func ErrorF(format string, args ...interface{}) {
	log.Error(format, args...)
}

func Info(msg interface{}) {
	log.Info("%v", msg)
}

func Debug(msg interface{}) {
	log.Debug("%v", msg)
}

func Trace(msg interface{}) {
	log.Trace("%v", msg)
}

func Warn(msg interface{}) {
	log.Warn("%v", msg)
}

func Error(msg interface{}) {
	log.Error("%v", msg)
}
