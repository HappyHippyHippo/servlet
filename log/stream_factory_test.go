package log

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewStreamFactory(t *testing.T) {
	t.Run("create a new config stream factory", func(t *testing.T) {
		if NewStreamFactory() == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_StreamFactory_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	strategy := NewMockStreamFactoryStrategy(ctrl)
	factory := NewStreamFactory()

	t.Run("error if passing a nil strategy", func(t *testing.T) {
		if err := factory.Register(nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if check := err.Error(); check != "Invalid nil 'strategy' argument" {
			t.Errorf("return the (%v) error", check)
		}
	})

	t.Run("register the stream factory strategy", func(t *testing.T) {
		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.(*streamFactory).strategies[0] != strategy {
			t.Errorf("didn't stored the strategy")
		}
	})
}

func Test_StreamFactory_Create(t *testing.T) {
	stype := "type"
	path := "path"
	format := "format"

	t.Run("error if the format is unrecognized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewStreamFactory()

		strategy := NewMockStreamFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(stype, path, format).Return(false).Times(1)
		factory.Register(strategy)

		if stream, err := factory.Create(stype, path, format); stream != nil {
			t.Errorf("returned an valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Unrecognized stream type : type" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the config stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewStreamFactory()

		stream := NewMockStream(ctrl)
		strategy := NewMockStreamFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(stype, path, format).Return(true).Times(1)
		strategy.EXPECT().Create(path, format).Return(stream, nil).Times(1)
		factory.Register(strategy)

		if stream, err := factory.Create(stype, path, format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(stream, stream) {
			t.Errorf("didn't returned the created stream")
		}
	})
}

func Test_StreamFactory_CreateConfig(t *testing.T) {
	t.Run("error on unrecognized type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewStreamFactory()

		conf := NewMockPartial(ctrl)
		strategy := NewMockStreamFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(false).Times(1)
		factory.Register(strategy)

		if stream, err := factory.CreateConfig(conf); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != fmt.Sprintf("Unrecognized stream config : %v", conf) {
			t.Errorf("returned the (%v) error", err)
		} else if stream != nil {
			t.Errorf("returned a config stream")
		}
	})

	t.Run("create the config stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewStreamFactory()

		conf := NewMockPartial(ctrl)
		stream := NewMockStream(ctrl)
		strategy := NewMockStreamFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(true).Times(1)
		strategy.EXPECT().CreateConfig(conf).Return(stream, nil).Times(1)
		factory.Register(strategy)

		if stream, err := factory.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(stream, stream) {
			t.Errorf("didn't returned the created stream")
		}
	})
}
