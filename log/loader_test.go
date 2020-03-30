package log

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewLoader(t *testing.T) {
	t.Run("should return nil when missing the logger formatter factory", func(t *testing.T) {
		action := "Creating a new logger loader without the logger formatter factory"

		expected := "Invalid nil 'formatterFactory' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		loader, err := NewLoader(nil, streamFactory, logger)

		if loader != nil {
			t.Errorf("%s return an unexpected valid reference to a new logger loader", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return nil when missing the logger stream factory", func(t *testing.T) {
		action := "Creating a new logger loader without the logger stream factory"

		expected := "Invalid nil 'streamFactory' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)
		logger := NewMockLogger(ctrl)

		loader, err := NewLoader(formatterFactory, nil, logger)

		if loader != nil {
			t.Errorf("%s return an unexpected valid reference to a new logger loader", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return nil when missing logger", func(t *testing.T) {
		action := "Creating a new logger loader without the logger"

		expected := "Invalid nil 'logger' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)

		loader, err := NewLoader(formatterFactory, streamFactory, nil)

		if loader != nil {
			t.Errorf("%s return an unexpected valid reference to a new logger loader", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return logger loader", func(t *testing.T) {
		action := "Creating a new logger loader"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		loader, err := NewLoader(formatterFactory, streamFactory, logger)

		if loader == nil {
			t.Errorf("%s didn't return the expected logger loader", action)
		}
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}
	})
}

func Test_Loader_Load(t *testing.T) {
	t.Run("should return a error if passing a nil reference to a config", func(t *testing.T) {
		action := "Loading the streams when missing the config"

		expectedError := "Invalid nil 'config' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		loader, _ := NewLoader(formatterFactory, streamFactory, logger)

		check := loader.Load(nil)
		if check == nil {
			t.Errorf("%s didn't returned the expected error", action)
		} else {
			if check.Error() != expectedError {
				t.Errorf("%s returned the (%s) error when expecting : (%s)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should no try to register nothing if no stream is present in the config", func(t *testing.T) {
		action := "Loading the streams when there are none in the config"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		config := NewMockConfig(ctrl)
		config.EXPECT().Get("log.streams").Return([]interface{}{}).Times(1)

		loader, _ := NewLoader(formatterFactory, streamFactory, logger)

		check := loader.Load(config)
		if check != nil {
			t.Errorf("%s returned the unexpected error : %v", action, check)
		}
	})

	t.Run("should no try to register nothing if no stream is present in the config", func(t *testing.T) {
		action := "Loading the streams when there are an invalid value in the config list entry"

		invalidList := "__dummy_invalid_list_"
		expectedError := "Error while parsing the list of streams"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		config := NewMockConfig(ctrl)
		config.EXPECT().Get("log.streams").Return(invalidList).Times(1)

		loader, _ := NewLoader(formatterFactory, streamFactory, logger)

		check := loader.Load(config)

		if check == nil {
			t.Errorf("%s didn't returned the expected error instance", action)
		} else {
			if check.Error() != expectedError {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should return the error that may occure when trying to retrieve the id", func(t *testing.T) {
		action := "Loading a logger entry when erroring reading the id"

		expectedError := "Error while parsing the logger stream entry : "

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		logger := NewMockLogger(ctrl)

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").DoAndReturn(func(key string) interface{} {
			panic("invalid convertion")
		}).Times(1)

		config := NewMockConfig(ctrl)
		config.EXPECT().Get("log.streams").Return([]interface{}{partial}).Times(1)

		loader, _ := NewLoader(formatterFactory, streamFactory, logger)

		check := loader.Load(config)

		if check == nil {
			t.Errorf("%s didn't returned an error", action)
		} else {
			if strings.Index(check.Error(), expectedError) != 0 {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should return the error that may occure when trying to create the stream", func(t *testing.T) {
		action := "Loading a logger entry when erroring when creating the stream"

		id := "id"
		expectedError := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(id).Times(1)

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		streamFactory.EXPECT().CreateConfig(partial).Return(nil, fmt.Errorf(expectedError)).Times(1)
		logger := NewMockLogger(ctrl)

		config := NewMockConfig(ctrl)
		config.EXPECT().Get("log.streams").Return([]interface{}{partial}).Times(1)

		loader, _ := NewLoader(formatterFactory, streamFactory, logger)

		check := loader.Load(config)

		if check == nil {
			t.Errorf("%s didn't returned an error", action)
		} else {
			if strings.Index(check.Error(), expectedError) != 0 {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should return the error that may occure when trying to register the created stream", func(t *testing.T) {
		action := "Loading a logger entry when erroring when register the created the stream"

		id := "id"
		expectedError := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(id).Times(1)

		stream := NewMockStream(ctrl)

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		streamFactory.EXPECT().CreateConfig(partial).Return(stream, nil).Times(1)
		logger := NewMockLogger(ctrl)
		logger.EXPECT().AddStream(id, stream).Return(fmt.Errorf(expectedError)).Times(1)

		config := NewMockConfig(ctrl)
		config.EXPECT().Get("log.streams").Return([]interface{}{partial}).Times(1)

		loader, _ := NewLoader(formatterFactory, streamFactory, logger)

		check := loader.Load(config)

		if check == nil {
			t.Errorf("%s didn't returned an error", action)
		} else {
			if strings.Index(check.Error(), expectedError) != 0 {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should correctly register the created stream", func(t *testing.T) {
		action := "Loading a logger stream"

		id := "id"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(id).Times(1)

		stream := NewMockStream(ctrl)

		formatterFactory := NewMockFormatterFactory(ctrl)
		streamFactory := NewMockStreamFactory(ctrl)
		streamFactory.EXPECT().CreateConfig(partial).Return(stream, nil).Times(1)
		logger := NewMockLogger(ctrl)
		logger.EXPECT().AddStream(id, stream).Return(nil).Times(1)

		config := NewMockConfig(ctrl)
		config.EXPECT().Get("log.streams").Return([]interface{}{partial}).Times(1)

		loader, _ := NewLoader(formatterFactory, streamFactory, logger)

		if check := loader.Load(config); check != nil {
			t.Errorf("%s returned a unexpected error : %v", action, check)
		}
	})
}
