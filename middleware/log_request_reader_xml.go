package middleware

import (
	"encoding/xml"
	"fmt"

	"github.com/happyhippyhippo/servlet"
)

type logRequestReaderXML struct {
	reader LogRequestReader
	model  interface{}
}

// NewLogRequestReaderXML @TODO
func NewLogRequestReaderXML(reader LogRequestReader, model interface{}) (LogRequestReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("Invalid nil 'reader' argument")
	}

	return &logRequestReaderXML{
		reader: reader,
		model:  model,
	}, nil
}

// Get @TODO
func (r logRequestReaderXML) Get(context servlet.Context) map[string]interface{} {
	data := r.reader.Get(context)

	if err := xml.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyXml"] = r.model
	}

	return data
}
