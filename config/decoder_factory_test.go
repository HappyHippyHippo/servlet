package config

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewDecoderFactory(t *testing.T) {
	t.Run("create a new config decoder factory", func(t *testing.T) {
		action := "Creating a new config decoder factory"

		factory := NewDecoderFactory()

		if factory == nil {
			t.Errorf("%s didn't return a valid reference to a new config decoder factory", action)
		}
	})
}

func Test_DecoderFactory_Register(t *testing.T) {
	t.Run("should return a error if passing a nil strategy", func(t *testing.T) {
		action := "Registering a nil strategy"

		expected := "Invalid nil 'strategy' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewDecoderFactory()

		err := factory.Register(nil)
		if err == nil {
			t.Errorf("%s didn't return the expected error", action)
		} else {
			if check := err.Error(); check != expected {
				t.Errorf("%s return the error (%s) when expecting (%s)", action, check, expected)
			}
		}
	})

	t.Run("should correctly register the decoder factory strategy", func(t *testing.T) {
		action := "Registering a decoder factory strategy"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockDecoderFactoryStrategy(ctrl)

		factory := NewDecoderFactory()

		err := factory.Register(strategy)
		if err != nil {
			t.Errorf("%s return a unexpected error : %s", action, err.Error())
		}

		if factory.(*decoderFactory).strategies[0] != strategy {
			t.Errorf("%s didn't stored the strategy in the factory", action)
		}
	})
}

func Test_DecoderFactory_Create(t *testing.T) {
	t.Run("should return a error signaling that the format is unrecognized", func(t *testing.T) {
		action := "Creating a invalid format decoder"

		format := "__format__"
		expected := "Unrecognized format type : __format__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl, "{}")
		strategy := NewMockDecoderFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format, reader).Return(false).Times(1)

		factory := NewDecoderFactory()
		factory.Register(strategy)

		decoder, err := factory.Create(format, reader)
		if err == nil {
			t.Errorf("%s didn't return the expected error", action)
		} else {
			if check := err.Error(); check != expected {
				t.Errorf("%s returned the error (%s) when expected (%s)", action, check, expected)
			}
		}
		if decoder != nil {
			t.Errorf("%s returned an unexpected yaml config decoder reference", action)
		}
	})

	t.Run("should create the requested yaml config decoder", func(t *testing.T) {
		action := "Creating a new yaml decoder"

		format := "__format__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		reader := NewMockReader(ctrl, "{}")
		decoder := NewMockDecoder(ctrl)

		strategy := NewMockDecoderFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format, reader).Return(true).Times(1)
		strategy.EXPECT().Create(reader).Return(decoder, nil).Times(1)

		factory := NewDecoderFactory()
		factory.Register(strategy)

		check, err := factory.Create(format, reader)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if !reflect.DeepEqual(check, decoder) {
			t.Errorf("%s didn't returned the created strategy", action)
		}
	})
}
