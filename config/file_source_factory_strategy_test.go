package config

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewFileSourceFactoryStrategy(t *testing.T) {
	t.Run("should return nil when missing file system adapter", func(t *testing.T) {
		action := "Creating a file source factory strategy without a file system adapter reference"

		expected := "Invalid nil 'fileSystem' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		decoderFactory := NewMockDecoderFactory(ctrl)

		strategy, err := NewFileSourceFactoryStrategy(nil, decoderFactory)

		if strategy != nil {
			t.Errorf("%s returned a valid file source factory strategy, expected nil", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return nil when missing decoder factory", func(t *testing.T) {
		action := "Creating a file source factory strategy without a decoder factory reference"

		expected := "Invalid nil 'decoderFactory' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		strategy, err := NewFileSourceFactoryStrategy(fileSystem, nil)

		if strategy != nil {
			t.Errorf("%s returned a valid file source factory strategy, expected nil", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("creates a new file source factory strategy", func(t *testing.T) {
		action := "Creating a file source factory strategy"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)

		strategy, err := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if strategy == nil {
			t.Errorf("%s didn't return a valid reference to a new file source factory strategy", action)
		}
	})
}

func Test_FileSourceFactoryStrategy_Accept(t *testing.T) {
	t.Run("should accept only file type", func(t *testing.T) {
		action := "Checking the accepting type"

		scenarios := []struct {
			stype    string
			expected bool
		}{
			{ // test file type
				stype:    SourceTypeFile,
				expected: true,
			},
			{ // test non-file type (observable_file)
				stype:    SourceTypeObservableFile,
				expected: false,
			},
		}

		for _, scn := range scenarios {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			fileSystem := NewMockFs(ctrl)
			decoderFactory := NewMockDecoderFactory(ctrl)

			strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

			if check := strategy.Accept(scn.stype, "path", "format"); check != scn.expected {
				t.Errorf("%s didn't returned the expected (%v) for the type (%s), returned (%v)", action, scn.expected, scn.stype, check)
			}
		}
	})

	t.Run("should not accept if at least 2 extra arguments are passed (the path and format)", func(t *testing.T) {
		action := "Checking the acceptance with less than 2 extra arguments"

		stype := SourceTypeFile
		path := "path"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)

		strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.Accept(stype, path) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if the path extra argument is not a string", func(t *testing.T) {
		action := "Checking the acceptance when the path extra argument not a string"

		stype := SourceTypeFile
		path := 1
		format := DecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)

		strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.Accept(stype, path, format) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if the format extra argument is not a string", func(t *testing.T) {
		action := "Checking the acceptance when the format extra argument not a string"

		stype := SourceTypeFile
		path := "path"
		format := 1

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)

		strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.Accept(stype, path, format) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})
}

func Test_FileSourceFactoryStrategy_AcceptConfig(t *testing.T) {
	t.Run("should not accept if there is not a type config entry", func(t *testing.T) {
		action := "Checking the acceptance of a config without a type field"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("type").DoAndReturn(func(key string) string {
			panic("invalid convertion")
		})

		fileSystem := NewMockFs(ctrl)
		decoderFactory := NewMockDecoderFactory(ctrl)

		strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if there is not a path config entry", func(t *testing.T) {
		action := "Checking the acceptance of a config without a path field"

		stype := SourceTypeFile

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

		strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if there is not a format config entry", func(t *testing.T) {
		action := "Checking the acceptance of a config without a format field"

		stype := SourceTypeFile
		path := "__path__"

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

		strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should not accept if the type field does not have the file value", func(t *testing.T) {
		action := "Checking the acceptance of a config when the type is not the value file"

		stype := SourceTypeObservableFile
		path := "__path__"
		format := DecoderFormatYAML

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

		strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if strategy.AcceptConfig(partial) {
			t.Errorf("%s didn't returned the expected false", action)
		}
	})

	t.Run("should accept if all the mandatory fields are present and have the file type", func(t *testing.T) {
		action := "Checking the acceptance of a config"

		stype := SourceTypeFile
		path := "__path__"
		format := DecoderFormatYAML

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

		strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

		if !strategy.AcceptConfig(partial) {
			t.Errorf("%s didn't returned the expected true value", action)
		}
	})
}

func Test_FileSourceFactoryStrategy_Create(t *testing.T) {
	t.Run("should create the requested file source", func(t *testing.T) {
		action := "Creating a new file source"

		path := "__path__"
		format := DecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)

		file := NewMockFile(ctrl)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(partial, nil).Times(1)
		decoder.EXPECT().Close().Times(1)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

		source, err := strategy.Create(path, format)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if source == nil {
			t.Errorf("%s didn't returned a valid file source reference", action)
		} else {
			switch source.(type) {
			case *fileSource:
			default:
				t.Errorf("%s didn't return a valid reference to a new file source reference", action)
			}
		}
	})
}

func Test_FileSourceFactoryStrategy_CreateConfig(t *testing.T) {
	t.Run("should create the requested file source", func(t *testing.T) {
		action := "Creating a new file source"

		path := "__path__"
		format := DecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		conf := NewMockPartial(ctrl)
		gomock.InOrder(
			conf.EXPECT().String("path").Return(path).Times(1),
			conf.EXPECT().String("format").Return(format).Times(1),
		)
		partial := NewMockPartial(ctrl)

		file := NewMockFile(ctrl)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(partial, nil).Times(1)
		decoder.EXPECT().Close().Times(1)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		strategy, _ := NewFileSourceFactoryStrategy(fileSystem, decoderFactory)

		source, err := strategy.CreateConfig(conf)
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}

		if source == nil {
			t.Errorf("%s didn't returned a valid file source reference", action)
		} else {
			switch source.(type) {
			case *fileSource:
			default:
				t.Errorf("%s didn't return a valid reference to a new file source reference", action)
			}
		}
	})
}
