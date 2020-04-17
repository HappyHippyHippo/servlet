package config

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewSourceFactory(t *testing.T) {
	t.Run("create a new config source factory", func(t *testing.T) {
		if factory := NewSourceFactory(); factory == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_SourceFactory_Register(t *testing.T) {
	t.Run("should return a error if passing a nil strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewSourceFactory()

		if err := factory.Register(nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'strategy' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the source factory strategy", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		strategy := NewMockSourceFactoryStrategy(ctrl)
		factory := NewSourceFactory()

		if err := factory.Register(strategy); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if factory.(*sourceFactory).strategies[0] != strategy {
			t.Errorf("didn't stored the strategy")
		}
	})
}

func Test_SourceFactory_Create(t *testing.T) {
	stype := "type"
	path := "path"
	format := "format"

	t.Run("error on unrecognized format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewSourceFactory()

		strategy := NewMockSourceFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(stype, path, format).Return(false).Times(1)
		factory.Register(strategy)

		if source, err := factory.Create(stype, path, format); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Unrecognized source type : type" {
			t.Errorf("returned the (%v) error", err)
		} else if source != nil {
			t.Errorf("didn't return the source")
		}
	})

	t.Run("create the requested config source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewSourceFactory()

		source := NewMockSource(ctrl)
		strategy := NewMockSourceFactoryStrategy(ctrl)
		strategy.EXPECT().Accept(stype, path, format).Return(true).Times(1)
		strategy.EXPECT().Create(path, format).Return(source, nil).Times(1)
		factory.Register(strategy)

		if check, err := factory.Create(stype, path, format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(check, source) {
			t.Errorf("didn't returned the created source")
		}
	})
}

func Test_SourceFactory_CreateConfig(t *testing.T) {
	t.Run("error on unrecognized format", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewSourceFactory()

		conf := NewMockPartial(ctrl)
		strategy := NewMockSourceFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(false).Times(1)
		factory.Register(strategy)

		if source, err := factory.CreateConfig(conf); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != fmt.Sprintf("Unrecognized source config : %v", conf) {
			t.Errorf("returned the (%v) error", err)
		} else if source != nil {
			t.Errorf("returned a valid reference")
		}
	})

	t.Run("create the config source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		factory := NewSourceFactory()

		conf := NewMockPartial(ctrl)
		source := NewMockSource(ctrl)
		strategy := NewMockSourceFactoryStrategy(ctrl)
		strategy.EXPECT().AcceptConfig(conf).Return(true).Times(1)
		strategy.EXPECT().CreateConfig(conf).Return(source, nil).Times(1)
		factory.Register(strategy)

		if check, err := factory.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !reflect.DeepEqual(check, source) {
			t.Errorf("didn't return the created source")
		}
	})
}
