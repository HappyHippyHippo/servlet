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

// NewLogRequestReaderXML will instantiate a new request event context
// reader XML decorator used to parse the request body as a XML and add
// the parsed content into the logging data.
func NewLogRequestReaderXML(reader LogRequestReader, model interface{}) (LogRequestReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("Invalid nil 'reader' argument")
	}

	return &logRequestReaderXML{
		reader: reader,
		model:  model,
	}, nil
}

// Get process the context request and add the extra bodyJson if the body
// content can be parsed as XML.
func (r logRequestReaderXML) Get(context servlet.Context) map[string]interface{} {
	data := r.reader.Get(context)

	if err := xml.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyXml"] = r.model
	}

	return data
}
