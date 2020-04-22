package middleware

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewLogResponseWriter(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	writer := NewMockGinResponseWriter(ctrl)

	t.Run("error when missing writer", func(t *testing.T) {
		if value, err := NewLogResponseWriter(nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'writer' argument" {
			t.Errorf("returned the (%v) error", err)
		} else if value != nil {
			t.Errorf("returned a valid reference")
		}
	})

	t.Run("creates a new log response writer", func(t *testing.T) {
		if value, err := NewLogResponseWriter(writer); err != nil {
			t.Errorf("return the (%v) error", err)
		} else if value == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_LogResponseWriter_Write(t *testing.T) {
	b := []byte{12, 34, 56}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ginWriter := NewMockGinResponseWriter(ctrl)
	ginWriter.EXPECT().Write(b).Times(1)

	writer := &logResponseWriter{body: &bytes.Buffer{}, ResponseWriter: ginWriter}

	t.Run("write to buffer and underlying writer", func(t *testing.T) {
		writer.Write(b)
		if !reflect.DeepEqual(writer.body.Bytes(), b) {
			t.Errorf("written (%v) bytes on buffer", writer.body)
		}
	})
}

func Test_LogResponseWriter_Body(t *testing.T) {
	b := []byte{12, 34, 56}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ginWriter := NewMockGinResponseWriter(ctrl)
	writer := &logResponseWriter{body: bytes.NewBuffer(b), ResponseWriter: ginWriter}

	t.Run("write to buffer and underlying writer", func(t *testing.T) {
		if !reflect.DeepEqual(writer.Body(), b) {
			t.Errorf("written (%v) bytes on buffer", writer.body)
		}
	})
}
