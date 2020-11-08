package servlet

import (
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_NewLogFormatterFactory(t *testing.T) {
	t.Run("new log formatter factory", func(t *testing.T) {
		if NewLogFormatterFactory() == nil {
			t.Error("didn't returned a valid reference")
		}
	})
}

func Test_LogFormatterFactory_Register(t *testing.T) {
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

		var factory *LogFormatterFactory
		_ = factory.Register(nil)
	})

	t.Run("nil strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogFormatterFactory()

		if err := factory.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'strategy' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the formatter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockLogFormatterFactoryStrategy(ctrl)
		factory := NewLogFormatterFactory()

		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.strategies[0] != strategy {
			t.Errorf("didn't stored the strategy")
		}
	})
}

func Test_LogFormatterFactory_Create(t *testing.T) {
	t.Run("unrecognized format", func(t *testing.T) {
		format := "invalid format"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogFormatterFactory()

		strategy := NewMockLogFormatterFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format).Return(false).Times(1)
		_ = factory.Register(strategy)

		if result, err := factory.Create(format); result != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized format type : invalid format" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the formatter", func(t *testing.T) {
		format := "format"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogFormatterFactory()

		formatter := NewLogFormatterJSON()
		strategy := NewMockLogFormatterFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format).Return(true).Times(1)
		strategy.EXPECT().Create().Return(formatter, nil).Times(1)
		_ = factory.Register(strategy)

		if formatter, err := factory.Create(format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(formatter, formatter) {
			t.Errorf("didn't returned the formatter")
		}
	})
}
