package servlet

import (
	"github.com/golang/mock/gomock"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_NewLogLoader(t *testing.T) {
	t.Run("error when missing the logger", func(t *testing.T) {
		streamFactory := NewLogStreamFactory()

		if loader, err := NewLogLoader(nil, streamFactory); loader != nil {
			t.Errorf("return a valid reference")
		} else if err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'logger' argument" {
			t.Errorf("returned the (%v)) error", err)
		}
	})

	t.Run("error when missing the logger stream factory", func(t *testing.T) {
		logger := NewLog()

		if loader, err := NewLogLoader(logger, nil); loader != nil {
			t.Errorf("return a valid reference")
		} else if err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'streamFactory' argument" {
			t.Errorf("returned the (%v)) error", err)
		}
	})

	t.Run("create loader", func(t *testing.T) {
		logger := NewLog()
		streamFactory := NewLogStreamFactory()

		if loader, err := NewLogLoader(logger, streamFactory); loader == nil {
			t.Errorf("didn't returned a valid reference")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}

func Test_LogLoader_Load(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		if err := loader.Load(nil); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'config' argument" {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("no-op if stream list is missing", func(t *testing.T) {
		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		config, _ := NewConfig(0 * time.Second)

		if err := loader.Load(config); err != nil {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("no-op if stream list is empty", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		conf := ConfigPartial{"log": ConfigPartial{"sources": []interface{}{}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)

		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("source", 0, source)

		if err := loader.Load(config); err != nil {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("error if stream list is not a list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		conf := ConfigPartial{"log": ConfigPartial{"streams": 123}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)

		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("source", 0, source)

		if err := loader.Load(config); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("missing stream id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		streamConfig := ConfigPartial{}
		conf := ConfigPartial{"log": ConfigPartial{"streams": []interface{}{streamConfig}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)

		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("source", 0, source)

		if err := loader.Load(config); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("invalid stream id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		streamConfig := ConfigPartial{"id": 123}
		conf := ConfigPartial{"log": ConfigPartial{"streams": []interface{}{streamConfig}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)
		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("source", 0, source)

		if err := loader.Load(config); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("error creating stream", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		loader, _ := NewLogLoader(logger, streamFactory)

		streamConfig := ConfigPartial{"id": "id"}
		conf := ConfigPartial{"log": ConfigPartial{"streams": []interface{}{streamConfig}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)
		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("source", 0, source)

		if err := loader.Load(config); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if strings.Index(err.Error(), "unrecognized stream config :") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("error storing stream", func(t *testing.T) {
		streamConfig := ConfigPartial{
			"id":       "id",
			"type":     "file",
			"path":     "path",
			"format":   "json",
			"channels": []interface{}{},
			"level":    "debug"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile("path", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatterFactory := NewLogFormatterFactory()
		_ = formatterFactory.Register(NewLogFormatterFactoryStrategyJSON())

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		fileStreamFactoryStrategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)
		_ = streamFactory.Register(fileStreamFactoryStrategy)
		loader, _ := NewLogLoader(logger, streamFactory)

		conf := ConfigPartial{"log": ConfigPartial{"streams": []interface{}{streamConfig}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)

		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("id", 0, source)

		writer := NewMockWriter(ctrl)
		formatter := NewLogFormatterJSON()
		fileLogger, _ := NewLogStreamFile(writer, formatter, []string{}, FATAL)
		_ = logger.AddStream("id", fileLogger)

		if err := loader.Load(config); err == nil {
			t.Errorf("didn't returned the expected error")
		} else if strings.Index(err.Error(), "duplicate id :") != 0 {
			t.Errorf("returned the (%s) error", err)
		}
	})

	t.Run("register stream", func(t *testing.T) {
		streamConfig := ConfigPartial{
			"id":       "id",
			"type":     "file",
			"path":     "path",
			"format":   "json",
			"channels": []interface{}{},
			"level":    "debug"}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile("path", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		formatterFactory := NewLogFormatterFactory()
		_ = formatterFactory.Register(NewLogFormatterFactoryStrategyJSON())

		logger := NewLog()
		streamFactory := NewLogStreamFactory()
		fileStreamFactoryStrategy, _ := NewLogStreamFactoryStrategyFile(fileSystem, formatterFactory)
		_ = streamFactory.Register(fileStreamFactoryStrategy)
		loader, _ := NewLogLoader(logger, streamFactory)

		conf := ConfigPartial{"log": ConfigPartial{"streams": []interface{}{streamConfig}}}
		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(conf).Times(1)

		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("id", 0, source)

		if err := loader.Load(config); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if !logger.HasStream("id") {
			t.Error("didn't stored the loaded stream")
		}
	})
}
