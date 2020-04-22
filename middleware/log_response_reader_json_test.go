package middleware

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewLogResponseReaderJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockLogResponseReader(ctrl)

	t.Run("error if reader is nil", func(t *testing.T) {
		if reader, err := NewLogResponseReaderJSON(nil, nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'reader' argument" {
			t.Errorf("returned the (%v) error", err)
		} else if reader != nil {
			t.Errorf("returned a valid reference")
		}
	})

	t.Run("creates a new decorator", func(t *testing.T) {
		if reader, err := NewLogResponseReaderJSON(reader, nil); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Errorf("didn't return a valid reference")
		}
	})

	t.Run("creates a new decorator with model", func(t *testing.T) {
		model := struct{ data string }{data: "bing"}

		if reader, err := NewLogResponseReaderJSON(reader, model); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if reader == nil {
			t.Errorf("didn't return a valid reference")
		} else if check := reader.(*logResponseReaderJSON).model; check != model {
			t.Errorf("stored the (%v) model", check)
		}
	})
}

func Test_LogResponseReaderJSON_Get(t *testing.T) {
	t.Run("non-json body does not add decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "{"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		reader := NewMockLogResponseReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		jsonReader, _ := NewLogResponseReaderJSON(reader, nil)
		result := jsonReader.Get(context)

		if _, ok := result["bodyJson"]; ok {
			t.Errorf("added the bodyJson field")
		}
	})

	t.Run("json body adds decorated field", func(t *testing.T) {
		data := map[string]interface{}{"body": "{\"field\":\"value\"}"}
		expected := map[string]interface{}{"field": "value"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		context := NewMockContext(ctrl)
		reader := NewMockLogResponseReader(ctrl)
		reader.EXPECT().Get(context).Return(data).Times(1)

		jsonReader, _ := NewLogResponseReaderJSON(reader, nil)
		result := jsonReader.Get(context)

		if body, ok := result["bodyJson"]; !ok {
			t.Errorf("didn't added the bodyJson field")
		} else if !reflect.DeepEqual(body, expected) {
			t.Errorf("added the (%v) json content)", body)
		}
	})
}
