package config

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewDecoderFactory(t *testing.T) {
	t.Run("create a new config decoder factory", func(t *testing.T) {
		if factory := NewDecoderFactory(); factory == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_DecoderFactory_Register(t *testing.T) {
	factory := NewDecoderFactory()

	t.Run("error if passing a nil strategy", func(t *testing.T) {
		if err := factory.Register(nil); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "Invalid nil 'strategy' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the decoder factory strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockDecoderFactoryStrategy(ctrl)

		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.(*decoderFactory).strategies[0] != strategy {
			t.Errorf("didn't stored the strategy")
		}
	})
}

func Test_DecoderFactory_Create(t *testing.T) {
	format := "format"
	expectedError := "Unrecognized format type : format"

	t.Run("error if the format is unrecognized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewDecoderFactory()

		reader := NewMockReader(ctrl, "{}")
		strategy := NewMockDecoderFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format, reader).Return(false).Times(1)
		factory.Register(strategy)

		if result, err := factory.Create(format, reader); result != nil {
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("should create the requested yaml config decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewDecoderFactory()

		reader := NewMockReader(ctrl, "{}")
		decoder := NewMockDecoder(ctrl)
		strategy := NewMockDecoderFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format, reader).Return(true).Times(1)
		strategy.EXPECT().Create(reader).Return(decoder, nil).Times(1)
		factory.Register(strategy)

		if check, err := factory.Create(format, reader); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(check, decoder) {
			t.Errorf("didn't returned the created strategy")
		}
	})
}
