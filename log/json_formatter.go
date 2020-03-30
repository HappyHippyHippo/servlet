package log

import (
	"encoding/json"
	"time"
)

var levelNameMap = map[Level]string{
	FATAL:   "fatal",
	ERROR:   "error",
	WARNING: "warning",
	NOTICE:  "notice",
	DEBUG:   "debug",
}

type jsonFormatter struct{}

// NewJSONFormatter method
func NewJSONFormatter() Formatter {
	return &jsonFormatter{}
}

func (f jsonFormatter) Format(level Level, fields F, message string) string {
	if fields == nil {
		fields = F{}
	}

	fields["time"] = time.Now().Format("2006-01-02T15:04:05-0700")
	fields["level"] = levelNameMap[level]
	fields["message"] = message

	bytes, _ := json.Marshal(fields)
	return string(bytes)
}
