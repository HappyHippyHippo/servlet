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

// NewLogResponseReaderXML will instantiate a new response event context
// reader XML decorator used to parse the response body as a XML and add
// the parsed content into the logging data.
func NewLogResponseReaderXML(reader LogResponseReader, model interface{}) (LogResponseReader, error) {
	if reader == nil {
		return nil, fmt.Errorf("Invalid nil 'reader' argument")
	}

	return &logResponseReaderXML{
		reader: reader,
		model:  model,
	}, nil
}

// Get process the context response and add the extra bodyJson if the body
// content can be parsed as XML.
func (r logResponseReaderXML) Get(context servlet.Context) map[string]interface{} {
	data := r.reader.Get(context)

	if err := xml.Unmarshal([]byte(data["body"].(string)), &r.model); err == nil {
		data["bodyXml"] = r.model
	}

	return data
}
