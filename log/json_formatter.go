package log

import (
	"encoding/json"
	"strings"
	"time"
)

type jsonFormatter struct{}

// NewJSONFormatter will instantiate a new JSON formatter that will take the
// logging entry request and create the output JSON string.
func NewJSONFormatter() Formatter {
	return &jsonFormatter{}
}

// Format will create the output JSON string message formatted with the content
// of the passed level, fields and message
func (f jsonFormatter) Format(level Level, message string, fields F) string {
	if fields == nil {
		fields = F{}
	}

	fields["time"] = time.Now().Format("2006-01-02T15:04:05.000-0700")
	fields["level"] = strings.ToUpper(LevelNameMap[level])
	fields["message"] = message

	bytes, _ := json.Marshal(fields)
	return string(bytes)
}
