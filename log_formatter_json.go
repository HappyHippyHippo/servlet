package servlet

import (
	"encoding/json"
	"strings"
	"time"
)

// LogFormatterJSON defines a JSON based log formatter.
type LogFormatterJSON struct{}

// NewLogFormatterJSON will instantiate a new JSON formatter that will take the
// logging entry request and create the output JSON string.
func NewLogFormatterJSON() LogFormatter {
	return &LogFormatterJSON{}
}

// Format will create the output JSON string message formatted with the content
// of the passed level, message and context
func (f LogFormatterJSON) Format(level LogLevel, message string, context map[string]interface{}) string {
	if context == nil {
		context = map[string]interface{}{}
	}

	context["time"] = time.Now().Format("2006-01-02T15:04:05.000-0700")
	context["level"] = strings.ToUpper(LogLevelNameMap[level])
	context["message"] = message

	bytes, _ := json.Marshal(context)
	return string(bytes)
}
