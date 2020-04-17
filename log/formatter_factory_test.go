package log

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewFormatterFactory(t *testing.T) {
	t.Run("create a new log formatter factory", func(t *testing.T) {
		if NewFormatterFactory() == nil {
			t.Errorf("return a valid reference")
		}
	})
}

func Test_FormatterFactory_Register(t *testing.T) {
	t.Run("error if passing a nil strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewFormatterFactory()

		if err := factory.Register(nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'strategy' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the formatter factory strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockFormatterFactoryStrategy(ctrl)
		factory := NewFormatterFactory()

		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.(*formatterFactory).strategies[0] != strategy {
			t.Errorf("didn't stored the strategy")
		}
	})
}

func Test_FormatterFactory_Create(t *testing.T) {
	format := "format"
	expectedError := "Unrecognized format type : format"

	t.Run("error if the format is unrecognized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewFormatterFactory()

		strategy := NewMockFormatterFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format).Return(false).Times(1)
		factory.Register(strategy)

		if result, err := factory.Create(format); result != nil {
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the formatter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewFormatterFactory()

		formatter := NewMockFormatter(ctrl)
		strategy := NewMockFormatterFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format).Return(true).Times(1)
		strategy.EXPECT().Create().Return(formatter, nil).Times(1)
		factory.Register(strategy)

		if formatter, err := factory.Create(format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(formatter, formatter) {
			t.Errorf("didn't returned the formatter")
		}
	})
}
