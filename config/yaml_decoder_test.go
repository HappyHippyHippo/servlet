package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewYamlDecoder(t *testing.T) {
	t.Run("should return nil when missing reader", func(t *testing.T) {
		action := "Creating a yaml decoder adapter without a reader reference"

		expected := "Invalid nil 'reader' argument"

		decoder, err := NewYamlDecoder(nil)

		if decoder != nil {
			t.Errorf("%s returned a valid yaml decoder adapter, expected nil", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("creates a new yaml decoder adapter", func(t *testing.T) {
		action := "Creating a yaml decoder adapter"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockFile(ctrl)
		reader.EXPECT().Close().Times(1)

		decoder, err := NewYamlDecoder(reader)

		if decoder == nil {
			t.Errorf("%s didn't return a valid reference to a new yaml decoder adapter", action)
		} else {
			decoder.Close()
		}
		if err != nil {
			t.Errorf("%s returned a unexpected error : %v", action, err)
		}
	})

	t.Run("should instantiate the yaml decoder with the injected reader", func(t *testing.T) {
		action := "Creating the yaml decoder adapter"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockFile(ctrl)
		reader.EXPECT().Close().Times(1)

		decoder, _ := NewYamlDecoder(reader)
		defer decoder.Close()

		if decoder.(*yamlDecoder).decoder == nil {
			t.Errorf("%s didn't instantiate the expected yaml.DecoderAdapter", action)
		}
	})
}

func Test_YamlDecoder_Close(t *testing.T) {
	t.Run("should call close method on reader if implements io.Closer interface", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockFile(ctrl)
		reader.EXPECT().Close().Times(1)

		decoder, _ := NewYamlDecoder(reader)

		decoder.Close()
		decoder.Close()
	})
}

func Test_YamlDecoder_Decode(t *testing.T) {
	t.Run("should return decode error if any", func(t *testing.T) {
		action := "Calling the decode method with decode error"

		expected := fmt.Errorf("__dummy_error__")

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		underlyingDecoder := NewMockUnderlyingYamlDecoder(ctrl)
		underlyingDecoder.EXPECT().Decode(&partial{}).DoAndReturn(func(p interface{}) error { return expected }).Times(1)

		reader := NewMockFile(ctrl)
		reader.EXPECT().Close().Times(1)

		decoder, _ := NewYamlDecoder(reader)
		defer decoder.Close()

		decoder.(*yamlDecoder).decoder = underlyingDecoder

		result, err := decoder.Decode()

		if result != nil {
			t.Errorf("%s returned an unexpected value reference to a partial : %v", action, result)
		}
		if err == nil {
			t.Errorf("%s didn't returned any error", action)
		} else {
			if err != expected {
				t.Errorf("%s returned the error (%v), expected (%v)", action, err, expected)
			}
		}
	})

	t.Run("should redirect to the underlying decoder", func(t *testing.T) {
		action := "Calling the decode method"

		expected := partial{"node": "value"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		underlyingDecoder := NewMockUnderlyingYamlDecoder(ctrl)
		underlyingDecoder.EXPECT().Decode(&partial{}).DoAndReturn(func(p interface{}) error { p = p.(*partial).merge(expected); return nil }).Times(1)

		reader := NewMockFile(ctrl)
		reader.EXPECT().Close().Times(1)

		decoder, _ := NewYamlDecoder(reader)
		defer decoder.Close()

		decoder.(*yamlDecoder).decoder = underlyingDecoder

		result, err := decoder.Decode()

		if result == nil {
			t.Errorf("%s returned an unexpected nil reference to the readed partial", action)
		} else {
			if !reflect.DeepEqual(result, expected) {
				t.Errorf("%s returned %v, expected %v", action, result, expected)
			}
		}
		if err != nil {
			t.Errorf("%s returned the unexpected error : %v", action, err)
		}
	})
}
