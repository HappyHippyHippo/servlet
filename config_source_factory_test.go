package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"reflect"
	"testing"
)

func Test_NewConfigSourceFactory(t *testing.T) {
	t.Run("new config source factory", func(t *testing.T) {
		if factory := NewConfigSourceFactory(); factory == nil {
			t.Error("didn't returned a valid reference")
		} else if factory.strategies == nil {
			t.Error("didn't instantiated the strategies storing array")
		}
	})
}

func Test_ConfigSourceFactory_Register(t *testing.T) {
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

		var factory *ConfigSourceFactory
		_ = factory.Register(nil)
	})

	t.Run("nil strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigSourceFactory()

		if err := factory.Register(nil); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'strategy' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the source factory strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockConfigSourceFactoryStrategy(ctrl)
		factory := NewConfigSourceFactory()

		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.strategies[0] != strategy {
			t.Error("didn't stored the strategy")
		}
	})
}

func Test_ConfigSourceFactory_Create(t *testing.T) {
	t.Run("error on unrecognized format", func(t *testing.T) {
		sourceType := "type"
		path := "path"
		format := "format"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigSourceFactory()

		strategy := NewMockConfigSourceFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(sourceType, path, format).Return(false).Times(1)
		_ = factory.Register(strategy)

		if source, err := factory.Create(sourceType, path, format); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "unrecognized source type : type" {
			t.Errorf("returned the (%v) error", err)
		} else if source != nil {
			t.Error("didn't returned the source")
		}
	})

	t.Run("create the requested config source", func(t *testing.T) {
		sourceType := "type"
		path := "path"
		format := "format"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigSourceFactory()

		source := &ConfigSourceBase{}
		strategy := NewMockConfigSourceFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(sourceType, path, format).Return(true).Times(1)
		strategy.EXPECT().Create(path, format).Return(source, nil).Times(1)
		_ = factory.Register(strategy)

		if check, err := factory.Create(sourceType, path, format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(check, source) {
			t.Error("didn't returned the created source")
		}
	})
}

func Test_ConfigSourceFactory_CreateConfig(t *testing.T) {
	t.Run("error on unrecognized format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigSourceFactory()

		conf := ConfigPartial{}
		strategy := NewMockConfigSourceFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(false).Times(1)
		_ = factory.Register(strategy)

		if source, err := factory.CreateConfig(conf); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != fmt.Sprintf("unrecognized source config : %v", conf) {
			t.Errorf("returned the (%v) error", err)
		} else if source != nil {
			t.Error("returned a valid reference")
		}
	})

	t.Run("create the config source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewConfigSourceFactory()

		conf := ConfigPartial{}
		source := NewMockConfigSource(ctrl)
		strategy := NewMockConfigSourceFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(true).Times(1)
		strategy.EXPECT().CreateConfig(conf).Return(source, nil).Times(1)
		_ = factory.Register(strategy)

		if check, err := factory.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(check, source) {
			t.Error("didn't returned the created source")
		}
	})
}
