package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewSourceFactory(t *testing.T) {
	t.Run("create a new config source factory", func(t *testing.T) {
		action := "Creating a new config source factory"

		factory := NewSourceFactory()

		if factory == nil {
			t.Errorf("%s didn't return a valid reference to a new source factory", action)
		}
	})
}

func Test_SourceFactory_Register(t *testing.T) {
	t.Run("should return a error if passing a nil strategy", func(t *testing.T) {
		action := "Registering a nil strategy"

		expected := "Invalid nil 'strategy' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewSourceFactory()

		err := factory.Register(nil)
		if err == nil {
			t.Errorf("%s didn't return the expected error", action)
		} else {
			if check := err.Error(); check != expected {
				t.Errorf("%s return the error (%s) when expecting (%s)", action, check, expected)
			}
		}
	})

	t.Run("should correctly register the source factory strategy", func(t *testing.T) {
		action := "Registering a source factory strategy"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockSourceFactoryStrategy(ctrl)

		factory := NewSourceFactory()

		err := factory.Register(strategy)
		if err != nil {
			t.Errorf("%s return a unexpected error : %s", action, err.Error())
		}

		if factory.(*sourceFactory).strategies[0] != strategy {
			t.Errorf("%s didn't stored the strategy in the factory", action)
		}
	})
}

func Test_SourceFactory_Create(t *testing.T) {
	t.Run("should return a error signaling that the format is unrecognized", func(t *testing.T) {
		action := "Creating a invalid format source"

		stype := "__type__"
		path := "__path__"
		format := "__format__"
		expected := "Unrecognized source type : __type__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockSourceFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(stype, path, format).Return(false).Times(1)

		factory := NewSourceFactory()
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
			t.Errorf("%s returned an unexpected config source reference", action)
		}
	})

	t.Run("should create the requested config source", func(t *testing.T) {
		action := "Creating a new source"

		stype := "__type__"
		path := "__path__"
		format := "__format__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		source := NewMockSource(ctrl)

		strategy := NewMockSourceFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(stype, path, format).Return(true).Times(1)
		strategy.EXPECT().Create(path, format).Return(source, nil).Times(1)

		factory := NewSourceFactory()
		factory.Register(strategy)

		check, err := factory.Create(stype, path, format)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if !reflect.DeepEqual(check, source) {
			t.Errorf("%s didn't returned the created source", action)
		}
	})
}

func Test_SourceFactory_CreateConfig(t *testing.T) {
	t.Run("should return a error signaling that the format is unrecognized", func(t *testing.T) {
		action := "Creating a invalid format source"

		expected := "Unrecognized source config : "

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)

		strategy := NewMockSourceFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(false).Times(1)

		factory := NewSourceFactory()
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
			t.Errorf("%s returned an unexpected config source reference", action)
		}
	})

	t.Run("should create the requested config source", func(t *testing.T) {
		action := "Creating a new source"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		source := NewMockSource(ctrl)

		strategy := NewMockSourceFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(true).Times(1)
		strategy.EXPECT().CreateConfig(conf).Return(source, nil).Times(1)

		factory := NewSourceFactory()
		factory.Register(strategy)

		check, err := factory.CreateConfig(conf)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if !reflect.DeepEqual(check, source) {
			t.Errorf("%s didn't returned the created source", action)
		}
	})
}
