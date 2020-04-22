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

// NewLogResponseReaderJSON @TODO
func NewLogResponseReaderJSON(reader LogResponseReader, model interface{}) (LogResponseReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("Invalid nil 'reader' argument")
	}

	return &logResponseReaderJSON{
		reader: reader,
		model:  model,
	}, nil
}

// Get @TODO
func (r logResponseReaderJSON) Get(context servlet.Context) map[string]interface{} {
	data := r.reader.Get(context)

	if err := json.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyJson"] = r.model
	}

	return data
}
