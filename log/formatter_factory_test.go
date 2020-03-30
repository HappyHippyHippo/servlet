package log

import (
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewFormatterFactory(t *testing.T) {
	t.Run("create a new log formatter factory", func(t *testing.T) {
		action := "Creating a new log formatter factory"

		factory := NewFormatterFactory()

		if factory == nil {
			t.Errorf("%s didn't return a valid reference to a new log formatter factory", action)
		}
	})
}

func Test_FormatterFactory_Register(t *testing.T) {
	t.Run("should return a error if passing a nil strategy", func(t *testing.T) {
		action := "Registering a nil strategy"

		expected := "Invalid nil 'strategy' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewFormatterFactory()

		err := factory.Register(nil)
		if err == nil {
			t.Errorf("%s didn't return the expected error", action)
		} else {
			if check := err.Error(); check != expected {
				t.Errorf("%s return the error (%s) when expecting (%s)", action, check, expected)
			}
		}
	})

	t.Run("should correctly register the formatter factory strategy", func(t *testing.T) {
		action := "Registering a formatter factory strategy"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockFormatterFactoryStrategy(ctrl)

		factory := NewFormatterFactory()

		err := factory.Register(strategy)
		if err != nil {
			t.Errorf("%s return a unexpected error : %s", action, err.Error())
		}

		if factory.(*formatterFactory).strategies[0] != strategy {
			t.Errorf("%s didn't stored the strategy in the factory", action)
		}
	})
}

func Test_FormatterFactory_Create(t *testing.T) {
	t.Run("should return a error signaling that the format is unrecognized", func(t *testing.T) {
		action := "Creating a invalid format formatter"

		format := "__format__"
		expected := "Unrecognized format type : __format__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockFormatterFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format).Return(false).Times(1)

		factory := NewFormatterFactory()
		factory.Register(strategy)

		formatter, err := factory.Create(format)
		if err == nil {
			t.Errorf("%s didn't return the expected error", action)
		} else {
			if check := err.Error(); check != expected {
				t.Errorf("%s returned the error (%s) when expected (%s)", action, check, expected)
			}
		}
		if formatter != nil {
			t.Errorf("%s returned an unexpected yaml config formatter reference", action)
		}
	})

	t.Run("should create the requested yaml config formatter", func(t *testing.T) {
		action := "Creating a new yaml formatter"

		format := "__format__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatter := NewMockFormatter(ctrl)

		strategy := NewMockFormatterFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(format).Return(true).Times(1)
		strategy.EXPECT().Create().Return(formatter, nil).Times(1)

		factory := NewFormatterFactory()
		factory.Register(strategy)

		check, err := factory.Create(format)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if !reflect.DeepEqual(check, formatter) {
			t.Errorf("%s didn't returned the created strategy", action)
		}
	})
}
