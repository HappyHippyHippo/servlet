package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewYamlDecoder(t *testing.T) {
	t.Run("error when missing reader", func(t *testing.T) {
		if decoder, err := NewYamlDecoder(nil); decoder != nil {
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'reader' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("creates a new yaml decoder adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockFile(ctrl)
		reader.EXPECT().Close().Times(1)

		if decoder, err := NewYamlDecoder(reader); decoder == nil {
			t.Errorf("didn't return a valid reference")
		} else {
			decoder.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})

	t.Run("instantiate with the injected reader", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockFile(ctrl)
		reader.EXPECT().Close().Times(1)

		decoder, _ := NewYamlDecoder(reader)
		defer decoder.Close()

		if decoder.(*yamlDecoder).reader != reader {
			t.Errorf("didn't store the passed reader")
		}
	})
}

func Test_YamlDecoder_Close(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	reader := NewMockFile(ctrl)
	reader.EXPECT().Close().Times(1)

	decoder, _ := NewYamlDecoder(reader)

	t.Run("call close method on reader only once", func(t *testing.T) {
		decoder.Close()
		decoder.Close()
	})
}

func Test_YamlDecoder_Decode(t *testing.T) {
	value := partial{"node": "value"}
	expectedError := "error"

	t.Run("return decode error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockFile(ctrl)
		reader.EXPECT().Close().Times(1)
		decoder, _ := NewYamlDecoder(reader)
		defer decoder.Close()

		underlyingDecoder := NewMockUnderlyingYamlDecoder(ctrl)
		underlyingDecoder.EXPECT().Decode(&partial{}).DoAndReturn(func(p interface{}) error {
			return fmt.Errorf(expectedError)
		}).Times(1)
		decoder.(*yamlDecoder).decoder = underlyingDecoder

		if result, err := decoder.Decode(); result != nil {
			t.Errorf("returned an reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("redirect to the underlying decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockFile(ctrl)
		reader.EXPECT().Close().Times(1)
		decoder, _ := NewYamlDecoder(reader)
		defer decoder.Close()

		underlyingDecoder := NewMockUnderlyingYamlDecoder(ctrl)
		underlyingDecoder.EXPECT().Decode(&partial{}).DoAndReturn(func(p interface{}) error {
			p = p.(*partial).merge(value)
			return nil
		}).Times(1)
		decoder.(*yamlDecoder).decoder = underlyingDecoder

		if result, err := decoder.Decode(); result == nil {
			t.Errorf("returned a nil value")
		} else if !reflect.DeepEqual(result, value) {
			t.Errorf("returned (%v)", result)
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}
