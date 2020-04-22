package middleware

import (
	"encoding/xml"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewLogRequestReaderXML(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockLogRequestReader(ctrl)

	t.Run("error if reader is nil", func(t *testing.T) {
		if reader, err := NewLogRequestReaderXML(nil, nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'reader' argument" {
			t.Errorf("returned the (%v) error", err)
		} else if reader != nil {
			t.Errorf("returned a valid reference")
		}
	})

	t.Run("creates a new decorator", func(t *testing.T) {
		if reader, err := NewLogRequestReaderXML(reader, nil); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Errorf("didn't return a valid reference")
		}
	})

	t.Run("creates a new decorator with model", func(t *testing.T) {
		model := struct{ data string }{data: "bing"}

		if reader, err := NewLogRequestReaderXML(reader, model); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Errorf("didn't return a valid reference")
		} else if check := reader.(*logRequestReaderXML).model; check != model {
			t.Errorf("stored the (%v) model", check)
		}
	})
}

func Test_LogRequestReaderXML_Get(t *testing.T) {
	model := struct {
		XMLName xml.Name `xml:"message"`
		Field   string   `xml:"field"`
	}{}

	t.Run("non-xml body does not add decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "{"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		reader := NewMockLogRequestReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		xmlReader, _ := NewLogRequestReaderXML(reader, &model)
		result := xmlReader.Get(context)

		if _, ok := result["bodyXml"]; ok {
			t.Errorf("added the bodyXml field")
		}
	})

	t.Run("xml body adds decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "\u003cmessage\u003e\u003cfield\u003evalue\u003c/field\u003e\u003c/message\u003e"}
		expected := struct {
			XMLName xml.Name `xml:"message"`
			Field   string   `xml:"field"`
		}{XMLName: xml.Name{Local: "message"}, Field: "value"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		reader := NewMockLogRequestReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		xmlReader, _ := NewLogRequestReaderXML(reader, &model)
		result := xmlReader.Get(context)

		if body, ok := result["bodyXml"]; !ok {
			t.Errorf("didn't added the bodyXml field")
		} else if !reflect.DeepEqual(body, &expected) {
			t.Errorf("added the (%v) xml content)", body)
		}
	})
}
