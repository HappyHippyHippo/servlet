package log

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewStreamFactory(t *testing.T) {
	t.Run("create a new config stream factory", func(t *testing.T) {
		action := "Creating a new config stream factory"

		factory := NewStreamFactory()

		if factory == nil {
			t.Errorf("%s didn't return a valid reference to a new stream factory", action)
		}
	})
}

func Test_StreamFactory_Register(t *testing.T) {
	t.Run("should return a error if passing a nil strategy", func(t *testing.T) {
		action := "Registering a nil strategy"

		expected := "Invalid nil 'strategy' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewStreamFactory()

		err := factory.Register(nil)
		if err == nil {
			t.Errorf("%s didn't return the expected error", action)
		} else {
			if check := err.Error(); check != expected {
				t.Errorf("%s return the error (%s) when expecting (%s)", action, check, expected)
			}
		}
	})

	t.Run("should correctly register the stream factory strategy", func(t *testing.T) {
		action := "Registering a stream factory strategy"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockStreamFactoryStrategy(ctrl)

		factory := NewStreamFactory()

		err := factory.Register(strategy)
		if err != nil {
			t.Errorf("%s return a unexpected error : %s", action, err.Error())
		}

		if factory.(*streamFactory).strategies[0] != strategy {
			t.Errorf("%s didn't stored the strategy in the factory", action)
		}
	})
}

func Test_StreamFactory_Create(t *testing.T) {
	t.Run("should return a error signaling that the format is unrecognized", func(t *testing.T) {
		action := "Creating a invalid format stream"

		stype := "__type__"
		path := "__path__"
		format := "__format__"
		expected := "Unrecognized stream type : __type__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockStreamFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(stype, path, format).Return(false).Times(1)

		factory := NewStreamFactory()
		factory.Register(strategy)

		decoder, err := factory.Create(stype, path, format)
		if err == nil {
			t.Errorf("%s didn't return the expected error", action)
		} else {
			if check := err.Error(); check != expected {
				t.Errorf("%s returned the error (%s) when expected (%s)", action, check, expected)
			}
		}
		if decoder != nil {
			t.Errorf("%s returned an unexpected config stream reference", action)
		}
	})

	t.Run("should create the requested config stream", func(t *testing.T) {
		action := "Creating a new stream"

		stype := "__type__"
		path := "__path__"
		format := "__format__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		stream := NewMockStream(ctrl)

		strategy := NewMockStreamFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(stype, path, format).Return(true).Times(1)
		strategy.EXPECT().Create(path, format).Return(stream, nil).Times(1)

		factory := NewStreamFactory()
		factory.Register(strategy)

		check, err := factory.Create(stype, path, format)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if !reflect.DeepEqual(check, stream) {
			t.Errorf("%s didn't returned the created stream", action)
		}
	})
}

func Test_StreamFactory_CreateConfig(t *testing.T) {
	t.Run("should return a error signaling that the format is unrecognized", func(t *testing.T) {
		action := "Creating a invalid format stream"

		expected := "Unrecognized stream config : "

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)

		strategy := NewMockStreamFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(false).Times(1)

		factory := NewStreamFactory()
		factory.Register(strategy)

		decoder, err := factory.CreateConfig(conf)
		if err == nil {
			t.Errorf("%s didn't return the expected error", action)
		} else {
			if check := err.Error(); check != fmt.Sprintf("%s%v", expected, conf) {
				t.Errorf("%s returned the error (%s) when expected (%s)", action, check, expected)
			}
		}
		if decoder != nil {
			t.Errorf("%s returned an unexpected config stream reference", action)
		}
	})

	t.Run("should create the requested config stream", func(t *testing.T) {
		action := "Creating a new stream"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		stream := NewMockStream(ctrl)

		strategy := NewMockStreamFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(true).Times(1)
		strategy.EXPECT().CreateConfig(conf).Return(stream, nil).Times(1)

		factory := NewStreamFactory()
		factory.Register(strategy)

		check, err := factory.CreateConfig(conf)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if !reflect.DeepEqual(check, stream) {
			t.Errorf("%s didn't returned the created stream", action)
		}
	})
}
