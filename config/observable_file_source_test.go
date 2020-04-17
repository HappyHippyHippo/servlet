package config

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func Test_NewObservableFileSource(t *testing.T) {
	path := "path"
	format := DecoderFormatYAML
	expectedError := "error"

	t.Run("error when missing the file system adapter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		decoderFactory := NewMockDecoderFactory(ctrl)

		if stream, err := NewObservableFileSource("path", DecoderFormatYAML, nil, decoderFactory); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'fileSystem' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error when missing the decoder factory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)

		if stream, err := NewObservableFileSource("path", DecoderFormatYAML, fileSystem, nil); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'decoderFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when retrieving the file info", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(nil, fmt.Errorf(expectedError)).Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)

		if stream, err := NewObservableFileSource(path, DecoderFormatYAML, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when opening the file", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)).Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)

		if stream, err := NewObservableFileSource(path, DecoderFormatYAML, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when creating the decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(nil, fmt.Errorf(expectedError)).Times(1)

		if stream, err := NewObservableFileSource(path, format, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error that may be raised when running the decoder", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(nil, fmt.Errorf(expectedError)).Times(1)
		decoder.EXPECT().Close().Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		if stream, err := NewObservableFileSource(path, format, fileSystem, decoderFactory); stream != nil {
			defer stream.Close()
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("create the config observable file source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		conf := NewMockPartial(ctrl)
		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(conf, nil).Times(1)
		decoder.EXPECT().Close().Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		if stream, err := NewObservableFileSource(path, format, fileSystem, decoderFactory); stream == nil {
			t.Errorf("didn't return a valid reference")
		} else {
			defer stream.Close()
			if err != nil {
				t.Errorf("returned the (%v) error", err)
			}
		}
	})

	t.Run("store the decoded partial", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(1)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		conf := NewMockPartial(ctrl)
		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(conf, nil).Times(1)
		decoder.EXPECT().Close().Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, _ := NewObservableFileSource(path, format, fileSystem, decoderFactory)
		defer stream.Close()

		if check := stream.(*observableFileSource).partial; check != conf {
			t.Errorf("didn't correctly stored the decoded partial")
		}
	})
}

func Test_ObservableFileSource_Reload(t *testing.T) {
	path := "path"
	format := DecoderFormatYAML
	expectedError := "error"

	t.Run("error if fail to retrieving the file info", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(1)
		fileSystem := NewMockFs(ctrl)
		gomock.InOrder(
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
			fileSystem.EXPECT().Stat(path).Return(nil, fmt.Errorf(expectedError)),
		)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		conf := NewMockPartial(ctrl)
		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(conf, nil).Times(1)
		decoder.EXPECT().Close().Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, _ := NewObservableFileSource(path, format, fileSystem, decoderFactory)
		defer stream.Close()

		if reloaded, err := stream.Reload(); reloaded {
			t.Errorf("flagged that was reloaded")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error if fails to load the file content", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileInfo := NewMockFileInfo(ctrl)
		gomock.InOrder(
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)),
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 2)),
		)
		fileSystem := NewMockFs(ctrl)
		gomock.InOrder(
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
			fileSystem.EXPECT().Stat(path).Return(fileInfo, nil),
		)
		gomock.InOrder(
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil),
			fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(nil, fmt.Errorf(expectedError)),
		)

		conf := NewMockPartial(ctrl)
		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(conf, nil).Times(1)
		decoder.EXPECT().Close().Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, _ := NewObservableFileSource(path, format, fileSystem, decoderFactory)
		defer stream.Close()

		if reloaded, err := stream.Reload(); reloaded {
			t.Errorf("flagged that was reloaded")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("prevent reload of a unchanged source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileInfo := NewMockFileInfo(ctrl)
		fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)).Times(2)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(2)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(1)

		conf := NewMockPartial(ctrl)
		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(conf, nil).Times(1)
		decoder.EXPECT().Close().Times(1)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(1)

		stream, _ := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		if reloaded, err := stream.Reload(); reloaded {
			t.Errorf("flagged that was reloaded")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("should reload a changed source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		file := NewMockFile(ctrl)
		fileInfo := NewMockFileInfo(ctrl)
		gomock.InOrder(
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 1)),
			fileInfo.EXPECT().ModTime().Return(time.Unix(0, 2)),
		)
		fileSystem := NewMockFs(ctrl)
		fileSystem.EXPECT().Stat(path).Return(fileInfo, nil).Times(2)
		fileSystem.EXPECT().OpenFile(path, os.O_RDONLY, os.FileMode(0644)).Return(file, nil).Times(2)

		conf := NewMockPartial(ctrl)
		decoder := NewMockDecoder(ctrl)
		decoder.EXPECT().Decode().Return(conf, nil).Times(2)
		decoder.EXPECT().Close().Times(2)
		decoderFactory := NewMockDecoderFactory(ctrl)
		decoderFactory.EXPECT().Create(format, file).Return(decoder, nil).Times(2)

		stream, _ := NewObservableFileSource(path, format, fileSystem, decoderFactory)

		if reloaded, err := stream.Reload(); !reloaded {
			t.Errorf("flagged that was not reloaded")
		} else if err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}
