package servlet

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

func Test_NewConfigLoader(t *testing.T) {
	config, _ := NewConfig(0 * time.Second)
	sourceFactory := NewConfigSourceFactory()

	t.Run("nil config", func(t *testing.T) {
		if loader, err := NewConfigLoader(nil, sourceFactory); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'config' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("nil config source factory", func(t *testing.T) {
		if loader, err := NewConfigLoader(config, nil); loader != nil {
			t.Error("returned a valid reference")
		} else if err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "invalid nil 'sourceFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("new config loader", func(t *testing.T) {
		if loader, err := NewConfigLoader(config, sourceFactory); loader == nil {
			t.Error("didn't returned a valid reference")
		} else if err != nil {
			t.Errorf("return the (%v) error", err)
		}
	})
}

func Test_ConfigLoader_Load(t *testing.T) {
	sourceID := "base_source_id"
	sourcePath := "base_source_path"
	sourceFormat := ConfigDecoderFormatYAML

	t.Run("error getting the base source", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error storing the base source", func(t *testing.T) {
		content := "field: value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(ConfigPartial{})
		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource(sourceID, 0, source)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "duplicate source id : base_source_id" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("add base source into the config", func(t *testing.T) {
		content := "field: value"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on invalid list of sources", func(t *testing.T) {
		content := "config:\n  sources: 123"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "error while parsing the list of sources" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on loaded invalid id", func(t *testing.T) {
		content := `
config:
  sources:
    - id: 12`

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on loaded invalid priority", func(t *testing.T) {
		content := `
config:
  sources:
    - id: id
      priority: string`

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "interface conversion") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on loaded source factory", func(t *testing.T) {
		content := `
config:
  sources:
    - id: id
      priority: 0`

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file := NewMockFile(ctrl)
		file.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file.EXPECT().Close().Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if strings.Index(err.Error(), "unrecognized source config") != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on source registration", func(t *testing.T) {
		content := `
config:
  sources:
    - id: id
      priority: 0
      type: file
      path: path
      format: yaml`

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		source := NewMockConfigSource(ctrl)
		source.EXPECT().Get("").Return(ConfigPartial{}).AnyTimes()
		config, _ := NewConfig(0 * time.Second)
		_ = config.AddSource("id", 0, source)

		file1 := NewMockFile(ctrl)
		file1.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file1.EXPECT().Close().Times(1)

		file2 := NewMockFile(ctrl)
		file2.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, "field: value")
			return 12, io.EOF
		}).Times(1)
		file2.EXPECT().Close().Times(1)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file1, nil).Times(1)
		fileSystem.EXPECT().OpenFile("path", os.O_RDONLY, os.FileMode(0644)).Return(file2, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Error("didn't returned the expected error")
		} else if err.Error() != "duplicate source id : id" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the loaded source", func(t *testing.T) {
		content := `
config:
  sources:
    - id: id
      priority: 0
      type: file
      path: path
      format: yaml`

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config, _ := NewConfig(0 * time.Second)

		file1 := NewMockFile(ctrl)
		file1.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, content)
			return len(content), io.EOF
		}).Times(1)
		file1.EXPECT().Close().Times(1)

		file2 := NewMockFile(ctrl)
		file2.EXPECT().Read(gomock.Any()).DoAndReturn(func(buf []byte) (int, error) {
			copy(buf, "field: value")
			return 12, io.EOF
		}).Times(1)
		file2.EXPECT().Close().Times(1)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(sourcePath, os.O_RDONLY, os.FileMode(0644)).Return(file1, nil).Times(1)
		fileSystem.EXPECT().OpenFile("path", os.O_RDONLY, os.FileMode(0644)).Return(file2, nil).Times(1)
		decoderFactory := NewConfigDecoderFactory()
		_ = decoderFactory.Register(NewConfigDecoderFactoryStrategyYaml())
		sourceFactory := NewConfigSourceFactory()
		fileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(fileSourceFactoryStrategy)
		observableFileSourceFactoryStrategy, _ := NewConfigSourceFactoryStrategyObservableFile(fileSystem, decoderFactory)
		_ = sourceFactory.Register(observableFileSourceFactoryStrategy)

		loader, _ := NewConfigLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}
