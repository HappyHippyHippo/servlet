package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/happyhippyhippo/servlet"
)

type logResponseReaderJSON struct {
	reader LogResponseReader
	model  interface{}
}

// NewLogResponseReaderJSON will instantiate a new response event context
// reader JSON decorator used to parse the response   body as a JSON and add
// the parsed content into the logging data.
func NewLogResponseReaderJSON(reader LogResponseReader, model interface{}) (LogResponseReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("Invalid nil 'reader' argument")
	}

	return &logResponseReaderJSON{
		reader: reader,
		model:  model,
	}, nil
}

// Get process the context response and add the extra bodyJson if the body
// content can be parsed as JSON.
func (r logResponseReaderJSON) Get(context servlet.Context) map[string]interface{} {
	data := r.reader.Get(context)

	if err := json.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyJson"] = r.model
	}

	return data
}
