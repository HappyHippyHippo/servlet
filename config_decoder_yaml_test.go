package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_NewConfigDecoderYaml(t *testing.T) {
	t.Run("nil reader", func(t *testing.T) {
		if decoder, err := NewConfigDecoderYaml(nil); decoder != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'reader' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new yaml decoder adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl)
		reader.EXPECT().Close().Times(1)

		if decoder, err := NewConfigDecoderYaml(reader); decoder == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			defer decoder.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			} else if decoder.reader != reader {
				t.Error("didn't store the reader reference")
			}
		}
	})
}

func Test_ConfigDecoderYaml_Close(t *testing.T) {
	t.Run("nil pointer receiver", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("didn't panic")
			} else {
				switch e := r.(type) {
				case error:
					if e.Error() != "nil pointer receiver" {
						t.Errorf("panic with the (%v) error", e)
					}
				default:
					t.Error("didn't panic with an error")
				}
			}
		}()

		var decoder *ConfigDecoderYaml
		decoder.Close()
	})

	t.Run("call close method on reader only once", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl)
		reader.EXPECT().Close().Times(1)
		decoder, _ := NewConfigDecoderYaml(reader)

		decoder.Close()
		decoder.Close()
	})
}

func Test_ConfigDecoderYaml_Decode(t *testing.T) {
	t.Run("return decode error", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl)
		reader.EXPECT().Close().Times(1)
		decoder, _ := NewConfigDecoderYaml(reader)
		defer decoder.Close()

		underlyingDecoder := NewMockUnderlyingConfigDecoderYaml(ctrl)
		underlyingDecoder.EXPECT().Decode(&ConfigPartial{}).DoAndReturn(func(p interface{}) error {
			return fmt.Errorf(expectedError)
		}).Times(1)
		decoder.decoder = underlyingDecoder

		if result, err := decoder.Decode(); result != nil {
			t.Error("returned an reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("redirect to the underlying decoder", func(t *testing.T) {
		value := ConfigPartial{"node": "value"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl)
		reader.EXPECT().Close().Times(1)
		decoder, _ := NewConfigDecoderYaml(reader)
		defer decoder.Close()

		underlyingDecoder := NewMockUnderlyingConfigDecoderYaml(ctrl)
		underlyingDecoder.EXPECT().Decode(&ConfigPartial{}).DoAndReturn(func(p interface{}) error {
			p = p.(*ConfigPartial).merge(value)
			return nil
		}).Times(1)
		decoder.decoder = underlyingDecoder

		if result, err := decoder.Decode(); result == nil {
			t.Error("returned a nil value")
		} else if !reflect.DeepEqual(result, value) {
			t.Errorf("returned (%v)", result)
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}
