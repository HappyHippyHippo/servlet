package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_NewLogStreamFactory(t *testing.T) {
	t.Run("new config stream factory", func(t *testing.T) {
		if NewLogStreamFactory() == nil {
			t.Errorf("didn't returned a valid reference")
		}
	})
}

func Test_LogStreamFactory_Register(t *testing.T) {
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

		var factory *LogStreamFactory
		_ = factory.Register(nil)
	})

	t.Run("nil strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogStreamFactory()

		if err := factory.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if check := err.Error(); check != "invalid nil 'strategy' argument" {
			t.Errorf("return the (%v) error", check)
		}
	})

	t.Run("register the stream factory strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockLogStreamFactoryStrategy(ctrl)
		factory := NewLogStreamFactory()

		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.strategies[0] != strategy {
			t.Error("didn't stored the strategy")
		}
	})
}

func Test_LogStreamFactory_Create(t *testing.T) {
	t.Run("unrecognized format", func(t *testing.T) {
		sourceType := "type"
		path := "path"
		format := "format"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogStreamFactory()

		strategy := NewMockLogStreamFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(sourceType, path, format).Return(false).Times(1)
		_ = factory.Register(strategy)

		if stream, err := factory.Create(sourceType, path, format); stream != nil {
			t.Error("returned an valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized stream type : type" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the config stream", func(t *testing.T) {
		sourceType := "type"
		path := "path"
		format := "format"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogStreamFactory()

		stream := NewMockLogStream(ctrl)
		strategy := NewMockLogStreamFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(sourceType, path, format).Return(true).Times(1)
		strategy.EXPECT().Create(path, format).Return(stream, nil).Times(1)
		_ = factory.Register(strategy)

		if stream, err := factory.Create(sourceType, path, format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(stream, stream) {
			t.Error("didn't returned the created stream")
		}
	})
}

func Test_LogStreamFactory_CreateConfig(t *testing.T) {
	t.Run("unrecognized type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogStreamFactory()

		conf := ConfigPartial{}
		strategy := NewMockLogStreamFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(false).Times(1)
		_ = factory.Register(strategy)

		if stream, err := factory.CreateConfig(conf); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != fmt.Sprintf("unrecognized stream config : %v", conf) {
			t.Errorf("returned the (%v) error", err)
		} else if stream != nil {
			t.Error("returned a config stream")
		}
	})

	t.Run("create the config stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewLogStreamFactory()

		conf := ConfigPartial{}
		stream := NewMockLogStream(ctrl)
		strategy := NewMockLogStreamFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(true).Times(1)
		strategy.EXPECT().CreateConfig(conf).Return(stream, nil).Times(1)
		_ = factory.Register(strategy)

		if stream, err := factory.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(stream, stream) {
			t.Error("didn't returned the created stream")
		}
	})
}
