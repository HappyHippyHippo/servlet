package servlet

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_NewConfigDecoderFactory(t *testing.T) {
	t.Run("create a new config decoder factory", func(t *testing.T) {
		if factory := NewConfigDecoderFactory(); factory == nil {
			t.Error("didn't returned a valid reference")
		} else if factory.strategies == nil {
			t.Errorf("didn't instantiated the strategies storing array")
		}
	})
}

func Test_ConfigDecoderFactory_Register(t *testing.T) {
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

		var factory *ConfigDecoderFactory
		_ = factory.Register(nil)
	})

	t.Run("nil strategy", func(t *testing.T) {
		factory := NewConfigDecoderFactory()
		if err := factory.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'strategy' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the decoder factory strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockConfigDecoderFactoryStrategy(ctrl)

		factory := NewConfigDecoderFactory()
		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.strategies[0] != strategy {
			t.Error("didn't stored the strategy")
		}
	})
}

func Test_ConfigDecoderFactory_Create(t *testing.T) {
	format := "format"

	t.Run("error if the format is unrecognized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigDecoderFactory()

		reader := NewMockReader(ctrl)
		strategy := NewMockConfigDecoderFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format, reader).Return(false).Times(1)
		_ = factory.Register(strategy)

		if result, err := factory.Create(format, reader); result != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized format type : format" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("should create the requested yaml config decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigDecoderFactory()

		reader := NewMockReader(ctrl)
		decoder := NewMockConfigDecoder(ctrl)
		strategy := NewMockConfigDecoderFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format, reader).Return(true).Times(1)
		strategy.EXPECT().Create(reader).Return(decoder, nil).Times(1)
		_ = factory.Register(strategy)

		if check, err := factory.Create(format, reader); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(check, decoder) {
			t.Error("didn't returned the created strategy")
		}
	})
}
