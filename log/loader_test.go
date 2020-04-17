package log

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewLoader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := NewMockLogger(ctrl)
	streamFactory := NewMockStreamFactory(ctrl)

	t.Run("error when missing the logger", func(t *testing.T) {
		if loader, err := NewLoader(nil, streamFactory); loader != nil {
			t.Errorf("return a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'logger' argument" {
			t.Errorf("returned the (%v)) error", err)
		}
	})

	t.Run("error when missing the logger stream factory", func(t *testing.T) {
		if loader, err := NewLoader(logger, nil); loader != nil {
			t.Errorf("return a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'streamFactory' argument" {
			t.Errorf("returned the (%v)) error", err)
		}
	})

	t.Run("create loader", func(t *testing.T) {
		if loader, err := NewLoader(logger, streamFactory); loader == nil {
			t.Errorf("didn't return a valid reference")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_Loader_Load(t *testing.T) {
	id := "id"
	expectedError := "error"

	t.Run("error if nil config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		loader, _ := NewLoader(logger, streamFactory)

		if err := loader.Load(nil); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'config' argument" {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("no-op if stream list is missing", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockConfig(ctrl)
		conf.EXPECT().Get("log.streams").Return(nil).Times(1)

		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		loader, _ := NewLoader(logger, streamFactory)

		if err := loader.Load(conf); err != nil {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("no-op if stream list is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockConfig(ctrl)
		conf.EXPECT().Get("log.streams").Return([]interface{}{}).Times(1)

		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		loader, _ := NewLoader(logger, streamFactory)

		if err := loader.Load(conf); err != nil {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("error if stream list is not a list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockConfig(ctrl)
		conf.EXPECT().Get("log.streams").Return("invalid_entry_list").Times(1)

		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		loader, _ := NewLoader(logger, streamFactory)

		if err := loader.Load(conf); err == nil {
			t.Errorf("didn't return the expected error")
		} else if strings.Index(err.Error(), "Error while parsing the logger stream entry : ") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("error retrieving stream id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").DoAndReturn(func(key string) interface{} {
			panic("invalid convertion")
		}).Times(1)

		conf := NewMockConfig(ctrl)
		conf.EXPECT().Get("log.streams").Return([]interface{}{partial}).Times(1)

		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		loader, _ := NewLoader(logger, streamFactory)

		if err := loader.Load(conf); err == nil {
			t.Errorf("didn't return the expected error")
		} else if strings.Index(err.Error(), "Error while parsing the logger stream entry : ") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("error creating stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(id).Times(1)

		conf := NewMockConfig(ctrl)
		conf.EXPECT().Get("log.streams").Return([]interface{}{partial}).Times(1)

		streamFactory := NewMockStreamFactory(ctrl)
		streamFactory.EXPECT().CreateConfig(partial).Return(nil, fmt.Errorf(expectedError)).Times(1)
		logger := NewMockLogger(ctrl)

		loader, _ := NewLoader(logger, streamFactory)

		if err := loader.Load(conf); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("error storing stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(id).Times(1)

		conf := NewMockConfig(ctrl)
		conf.EXPECT().Get("log.streams").Return([]interface{}{partial}).Times(1)

		stream := NewMockStream(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		streamFactory.EXPECT().CreateConfig(partial).Return(stream, nil).Times(1)
		logger := NewMockLogger(ctrl)
		logger.EXPECT().AddStream(id, stream).Return(fmt.Errorf(expectedError)).Times(1)

		loader, _ := NewLoader(logger, streamFactory)

		if err := loader.Load(conf); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("register stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(id).Times(1)

		conf := NewMockConfig(ctrl)
		conf.EXPECT().Get("log.streams").Return([]interface{}{partial}).Times(1)

		stream := NewMockStream(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		streamFactory.EXPECT().CreateConfig(partial).Return(stream, nil).Times(1)
		logger := NewMockLogger(ctrl)
		logger.EXPECT().AddStream(id, stream).Return(nil).Times(1)

		loader, _ := NewLoader(logger, streamFactory)

		if err := loader.Load(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}
