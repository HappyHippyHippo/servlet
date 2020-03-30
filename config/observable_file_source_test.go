package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_NewObservableFileSource(t *testing.T) {
	t.Run("should return an error when missing the file system", func(t *testing.T) {
		action := "Creating a new observable file source without the file system"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "Invalid nil 'fileSystem' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		decoderFactory := NewMockDecoderFactory(ctrl)

		stream, err := NewObservableFileSource(path, format, nil, decoderFactory)

		if stream != nil {
			t.Errorf("%s return an unexpected valid reference to a new observable file source", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return an error when missing the decoder factory", func(t *testing.T) {
		action := "Creating a new observable file source without the decoder factory"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "Invalid nil 'decoderFactory' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		stream, err := NewObservableFileSource(path, format, fileSystem, nil)

		if stream != nil {
			t.Errorf("%s return an unexpected valid reference to a new observable file source", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return the error that may be raised when retrieving the file info", func(t *testing.T) {
		action := "Creating a new observable file source when erroring when retrieving the file info"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(nil, fmt.Errorf(expected)).Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)

		stream, err := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		if stream != nil {
			stream.Close()
			t.Errorf("%s return an unexpected valid reference to a new observable file source", action)
		}

		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return the error that may be raised when opening the file", func(t *testing.T) {
		action := "Creating a new observable file source when erroring opening the file"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expected)).Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)

		stream, err := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		if stream != nil {
			stream.Close()
			t.Errorf("%s return an unexpected valid reference to a new observable file source", action)
		}

		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return the error that may be raised when creating the decoder", func(t *testing.T) {
		action := "Creating a new observable file source when erroing when creating the decoder"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(nil, fmt.Errorf(expected)).Times(1)

		stream, err := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		if stream != nil {
			stream.Close()
			t.Errorf("%s return an unexpected valid reference to a new observable file source", action)
		}

		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return the error that may be raised when running the decoder", func(t *testing.T) {
		action := "Creating a new observable file source when erroing when running the decoder"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(nil, fmt.Errorf(expected)).Times(1)
		decoder.EXPECT().Close().Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, err := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		if stream != nil {
			stream.Close()
			t.Errorf("%s return an unexpected valid reference to a new observable file source", action)
		}

		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should correctly create the config observable file source", func(t *testing.T) {
		action := "Creating a new file source"

		path := "__dummy_path__"
		format := DecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		partial := NewMockPartial(ctrl)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(partial, nil).Times(1)
		decoder.EXPECT().Close().Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, err := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		if stream == nil {
			t.Errorf("%s didn't returned the expected config observable file source", action)
		}
		defer stream.Close()
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}
	})
}

func Test_ObservableFileSource_Reload(t *testing.T) {
	t.Run("should return false if fail to retrieving the file info", func(t *testing.T) {
		action := "Reloading when failing to retrieve the file info"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		gomock.InOrder(
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
			fileSystem.EXPECT().Stat(path).Return(nil, fmt.Errorf(expected)),
		)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		partial := NewMockPartial(ctrl)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(partial, nil).Times(1)
		decoder.EXPECT().Close().Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, _ := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		reloaded, err := stream.Reload()
		if reloaded {
			t.Errorf("%s unexpectedly flagged that was reloaded", action)
		}
		if err == nil {
			t.Errorf("%s didn't returned the expected error, returned nil", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return false if fail to load the file content", func(t *testing.T) {
		action := "Reloading when failing to load the file content"

		path := "__dummy_path__"
		format := DecoderFormatYAML
		expected := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		gomock.InOrder(
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)),
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 2)),
		)

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		gomock.InOrder(
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
		)
		gomock.InOrder(
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil),
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expected)),
		)

		partial := NewMockPartial(ctrl)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(partial, nil).Times(1)
		decoder.EXPECT().Close().Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, _ := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		reloaded, err := stream.Reload()
		if reloaded {
			t.Errorf("%s unexpectedly flagged that was reloaded", action)
		}
		if err == nil {
			t.Errorf("%s didn't returned the expected error, returned nil", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should prevent reload of an unchanged source", func(t *testing.T) {
		action := "Reloading an unchanged source"

		path := "__dummy_path__"
		format := DecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(2)

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(2)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		partial := NewMockPartial(ctrl)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(partial, nil).Times(1)
		decoder.EXPECT().Close().Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, _ := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		reloaded, err := stream.Reload()
		if reloaded {
			t.Errorf("%s unexpectedly flagged that was reloaded", action)
		}
		if err != nil {
			t.Errorf("%s returned the unexpected error : %v", action, err)
		}
	})

	t.Run("should reload of a changed source", func(t *testing.T) {
		action := "Reloading an changed source"

		path := "__dummy_path__"
		format := DecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		gomock.InOrder(
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)),
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 2)),
		)

		file := NewMockFile(ctrl)

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(2)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(2)

		partial := NewMockPartial(ctrl)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(partial, nil).Times(2)
		decoder.EXPECT().Close().Times(2)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(2)

		stream, _ := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		reloaded, err := stream.Reload()
		if !reloaded {
			t.Errorf("%s unexpectedly flagged that was not reloaded", action)
		}
		if err != nil {
			t.Errorf("%s returned the unexpected error : %v", action, err)
		}
	})
}
