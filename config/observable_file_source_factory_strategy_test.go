package config

import (
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_NewObservableFileSourceFactoryStrategy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	decoderFactory := NewMockDecoderFactory(ctrl)

	t.Run("error when missing file system adapter", func(t *testing.T) {
		if strategy, err := NewObservableFileSourceFactoryStrategy(nil, decoderFactory); strategy != nil {
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error when missing decoder factory", func(t *testing.T) {
		if strategy, err := NewObservableFileSourceFactoryStrategy(fileSystem, nil); strategy != nil {
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'decoderFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("creates a new observable file source factory strategy", func(t *testing.T) {
		if strategy, err := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory); err != nil {
			t.Errorf("return the (%v) error", err)
		} else if strategy == nil {
			t.Errorf("didn't return a valid reference")
		}
	})
}

func Test_ObservableFileSourceFactoryStrategy_Accept(t *testing.T) {
	stype := SourceTypeObservableFile
	path := "path"
	format := DecoderFormatYAML

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	fileSystem := NewMockFs(ctrl)
	decoderFactory := NewMockDecoderFactory(ctrl)
	strategy, _ := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

	t.Run("don't accept if at least 2 extra arguments are passed", func(t *testing.T) {
		if strategy.Accept(stype, path) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if the path is not a string", func(t *testing.T) {
		if strategy.Accept(stype, 1, format) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if the format is not a string", func(t *testing.T) {
		if strategy.Accept(stype, path, 1) {
			t.Errorf("returned true")
		}
	})

	t.Run("accept only observable file type", func(t *testing.T) {
		scenarios := []struct {
			stype    string
			expected bool
		}{
			{ // test file type
				stype:    SourceTypeObservableFile,
				expected: true,
			},
			{ // test non-file type (file)
				stype:    SourceTypeFile,
				expected: false,
			},
		}

		for _, scn := range scenarios {
			if check := strategy.Accept(scn.stype, path, format); check != scn.expected {
				t.Errorf("for the type (%s), returned (%v)", scn.stype, check)
			}
		}
	})
}

func Test_ObservableFileSourceFactoryStrategy_AcceptConfig(t *testing.T) {
	stype := SourceTypeObservableFile
	path := "path"
	format := DecoderFormatYAML

	t.Run("don't accept if type is missing or is not a string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("type").DoAndReturn(func(key string) string {
			panic("invalid convertion")
		})

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)
		strategy, _ := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if path is missing or is not a string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return(stype).Times(1),
			partial.EXPECT().String("path").DoAndReturn(func(key string) string {
				panic("invalid convertion")
			}),
		)

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)
		strategy, _ := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if format is missing or is not a string", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return(stype).Times(1),
			partial.EXPECT().String("path").Return(path).Times(1),
			partial.EXPECT().String("format").DoAndReturn(func(key string) string {
				panic("invalid convertion")
			}),
		)

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)
		strategy, _ := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("don't accept if invalid type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return("invalid-type").Times(1),
			partial.EXPECT().String("path").Return(path).Times(1),
			partial.EXPECT().String("format").Return(format).Times(1),
		)

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)
		strategy, _ := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("returned true")
		}
	})

	t.Run("accept config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		gomock.InOrder(
			partial.EXPECT().String("type").Return(stype).Times(1),
			partial.EXPECT().String("path").Return(path).Times(1),
			partial.EXPECT().String("format").Return(format).Times(1),
		)

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)
		strategy, _ := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if !strategy.AcceptConfig(partial) {
			t.Errorf("returned false")
		}
	})
}

func Test_ObservableFileSourceFactoryStrategy_Create(t *testing.T) {
	path := "path"
	format := DecoderFormatYAML

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	file := NewMockFile(ctrl)
	fileInfo := NewMockFileInfo(ctrl)
	fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
	fileSystem := NewMockFs(ctrl)
	fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
	fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

	partial := NewMockPartial(ctrl)
	decoder := NewMockDecoder(ctrl)
	decoder.EXPECT().Decode().Return(partial, nil).Times(1)
	decoder.EXPECT().Close().Times(1)

	decoderFactory := NewMockDecoderFactory(ctrl)
	decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

	strategy, _ := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

	t.Run("create the observable file source", func(t *testing.T) {
		if source, err := strategy.Create(path, format); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if source == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			switch source.(type) {
			case *observableFileSource:
			default:
				t.Errorf("didn't return a new observable file source")
			}
		}
	})
}

func Test_ObservableFileSourceFactoryStrategy_CreateConfig(t *testing.T) {
	path := "path"
	format := DecoderFormatYAML

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	file := NewMockFile(ctrl)
	fileInfo := NewMockFileInfo(ctrl)
	fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
	fileSystem := NewMockFs(ctrl)
	fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
	fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

	partial := NewMockPartial(ctrl)
	decoder := NewMockDecoder(ctrl)
	decoder.EXPECT().Decode().Return(partial, nil).Times(1)
	decoder.EXPECT().Close().Times(1)
	decoderFactory := NewMockDecoderFactory(ctrl)
	decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

	conf := NewMockPartial(ctrl)
	gomock.InOrder(
		conf.EXPECT().String("path").Return(path).Times(1),
		conf.EXPECT().String("format").Return(format).Times(1),
	)

	strategy, _ := NewObservableFileSourceFactoryStrategy(fileSystem, decoderFactory)

	t.Run("create the observable file source", func(t *testing.T) {
		if source, err := strategy.CreateConfig(conf); err != nil {
			t.Errorf("returned the (%v) error", err)
		} else if source == nil {
			t.Errorf("didn't returned a valid reference")
		} else {
			switch source.(type) {
			case *observableFileSource:
			default:
				t.Errorf("didn't return a new observable file source")
			}
		}
	})
}
