package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/happyhippyhippo/servlet"
)

type logRequestReaderJSON struct {
	reader LogRequestReader
	model  interface{}
}

// NewLogRequestReaderJSON will instantiate a new request event context
// reader JSON decorator used to parse the request body as a JSON and add
// the parsed content into the logging data.
func NewLogRequestReaderJSON(reader LogRequestReader, model interface{}) (LogRequestReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("Invalid nil 'reader' argument")
	}

	return &logRequestReaderJSON{
		reader: reader,
		model:  model,
	}, nil
}

// Get process the context request and add the extra bodyJson if the body
// content can be parsed as JSON.
func (r logRequestReaderJSON) Get(context servlet.Context) map[string]interface{} {
	data := r.reader.Get(context)

	if err := json.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyJson"] = r.model
	}

	return data
}