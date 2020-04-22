package middleware

import (
	"encoding/xml"
	"fmt"

	"github.com/happyhippyhippo/servlet"
)

type logResponseReaderXML struct {
	reader LogResponseReader
	model  interface{}
}

// NewLogResponseReaderXML @TODO
func NewLogResponseReaderXML(reader LogResponseReader, model interface{}) (LogResponseReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("Invalid nil 'reader' argument")
	}

	return &logResponseReaderXML{
		reader: reader,
		model:  model,
	}, nil
}

// Get @TODO
func (r logResponseReaderXML) Get(context servlet.Context) map[string]interface{} {
	data := r.reader.Get(context)

	if err := xml.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyXml"] = r.model
	}

	return data
}
