package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewLoader(t *testing.T) {
	t.Run("should return an error when missing the config", func(t *testing.T) {
		action := "Creating a new config loader without the config"

		expected := "Invalid nil 'config' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sourceFactory := NewMockSourceFactory(ctrl)

		loader, err := NewLoader(nil, sourceFactory)

		if loader != nil {
			t.Errorf("%s return an unexpected valid reference to a new config", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return an error when missing the config source factory", func(t *testing.T) {
		action := "Creating a new config loader without the config source factory"

		expected := "Invalid nil 'sourceFactory' argument"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := NewMockConfig(ctrl)

		loader, err := NewLoader(config, nil)

		if loader != nil {
			t.Errorf("%s return an unexpected valid reference to a new config", action)
		}
		if err == nil {
			t.Errorf("%s didn't return a expected error", action)
		} else {
			if err.Error() != expected {
				t.Errorf("%s didn't return the expected return error (%s), expected (%s)", action, err.Error(), expected)
			}
		}
	})

	t.Run("should return the config loader", func(t *testing.T) {
		action := "Creating a new config loader"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := NewMockConfig(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)

		loader, err := NewLoader(config, sourceFactory)

		if loader == nil {
			t.Errorf("%s didn't return the expected valid reference to a new config", action)
		}
		if err != nil {
			t.Errorf("%s return the unexpected error : %v", action, err)
		}
	})
}

func Test_Loader_Load(t *testing.T) {
	t.Run("should return the error that may occure while getting the base source", func(t *testing.T) {
		action := "Loading the configuration when erroring loading the base source"

		sourceID := "__dummy_base_source_id__"
		sourcePath := "__dummy_base_source_path__"
		sourceFormat := DecoderFormatYAML
		expected := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		config := NewMockConfig(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(nil, fmt.Errorf(expected)).Times(1)

		loader, _ := NewLoader(config, sourceFactory)

		check := loader.Load(sourceID, sourcePath, sourceFormat)

		if check == nil {
			t.Errorf("%s didn't returned the expected error instance", action)
		} else {
			if check.Error() != expected {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expected)
			}
		}
	})

	t.Run("should return the error that may occure while storing the base source into the config", func(t *testing.T) {
		action := "Storing the loaded base config source into the config with error"

		sourceID := "__dummy_base_source_id__"
		sourcePath := "__dummy_base_source_path__"
		sourceFormat := DecoderFormatYAML
		expectedError := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		source := NewMockSource(ctrl)

		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, source).Return(fmt.Errorf(expectedError)).Times(1)

		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(source, nil).Times(1)

		loader, _ := NewLoader(config, sourceFactory)
		check := loader.Load(sourceID, sourcePath, sourceFormat)

		if check == nil {
			t.Errorf("%s didn't returned the expected error instance", action)
		} else {
			if check.Error() != expectedError {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should add the loaded base source into the config", func(t *testing.T) {
		action := "Loading the base config source into the config"

		sourceID := "__dummy_base_source_id__"
		sourcePath := "__dummy_base_source_path__"
		sourceFormat := DecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		source := NewMockSource(ctrl)

		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, source).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return(nil).Times(1)

		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(source, nil).Times(1)

		loader, _ := NewLoader(config, sourceFactory)

		if check := loader.Load(sourceID, sourcePath, sourceFormat); check != nil {
			t.Errorf("%s returned the unexpected error : %v", action, check.Error())
		}
	})

	t.Run("should add the loaded base source into the config and retrieve the list of loaded sources", func(t *testing.T) {
		action := "Loading an empty list of loaded sources"

		sourceID := "__dummy_base_source_id__"
		sourcePath := "__dummy_base_source_path__"
		sourceFormat := DecoderFormatYAML

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		source := NewMockSource(ctrl)

		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, source).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return([]interface{}{}).Times(1)

		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(source, nil).Times(1)

		loader, _ := NewLoader(config, sourceFactory)

		if check := loader.Load(sourceID, sourcePath, sourceFormat); check != nil {
			t.Errorf("%s returned the unexpected error : %v", action, check.Error())
		}
	})

	t.Run("should return a error when trying to retrieve a invalid list of loading sources", func(t *testing.T) {
		action := "Loading an invalid list of loading sources"

		sourceID := "__dummy_base_source_id__"
		sourcePath := "__dummy_base_source_path__"
		sourceFormat := DecoderFormatYAML
		invalidList := "__dummy_invalid_list_"
		expectedError := "Error while parsing the list of sources"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		source := NewMockSource(ctrl)

		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, source).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return(invalidList).Times(1)

		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(source, nil).Times(1)

		loader, _ := NewLoader(config, sourceFactory)
		check := loader.Load(sourceID, sourcePath, sourceFormat)

		if check == nil {
			t.Errorf("%s didn't returned the expected error instance", action)
		} else {
			if check.Error() != expectedError {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should return a error when trying to retrieve the a invalid id of the loaded source config", func(t *testing.T) {
		action := "Loading the config source into the config when erroring while retrieving the source id from config"

		sourceID := "__dummy_base_source_id__"
		sourcePath := "__dummy_base_source_path__"
		sourceFormat := DecoderFormatYAML
		expectedError := "Error while parsing the config entry : "

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").DoAndReturn(func(key string) string {
			panic("invalid convertion")
		}).Times(1)

		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return([]interface{}{partial}).Times(1)

		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)

		loader, _ := NewLoader(config, sourceFactory)
		check := loader.Load(sourceID, sourcePath, sourceFormat)

		if check == nil {
			t.Errorf("%s didn't returned the expected error instance", action)
		} else {
			if strings.Index(check.Error(), expectedError) != 0 {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should return a error when trying to retrieve the a invalid priority of the loaded source config", func(t *testing.T) {
		action := "Loading the config source into the config when erroring while retrieving the source priority from config"

		sourceID := "__dummy_base_source_id__"
		sourcePath := "__dummy_base_source_path__"
		sourceFormat := DecoderFormatYAML
		loadedID := "__dummy_loaded_id__"
		expectedError := "Error while parsing the config entry : "

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(loadedID).Times(1)
		partial.EXPECT().Int("priority").DoAndReturn(func(key string) string {
			panic("invalid convertion")
		}).Times(1)

		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return([]interface{}{partial}).Times(1)

		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)

		loader, _ := NewLoader(config, sourceFactory)
		check := loader.Load(sourceID, sourcePath, sourceFormat)

		if check == nil {
			t.Errorf("%s didn't returned the expected error instance", action)
		} else {
			if strings.Index(check.Error(), expectedError) != 0 {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should return the error that the source factory may return", func(t *testing.T) {
		action := "Loading the config source and the source factory returns an error"

		sourceID := "__dummy_base_source_id__"
		sourcePath := "__dummy_base_source_path__"
		sourceFormat := DecoderFormatYAML
		loadedID := "__dummy_loaded_id__"
		loadedPriority := 1
		expectedError := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(loadedID).Times(1)
		partial.EXPECT().Int("priority").Return(loadedPriority).Times(1)

		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return([]interface{}{partial}).Times(1)

		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		sourceFactory.EXPECT().CreateConfig(partial).Return(nil, fmt.Errorf(expectedError)).Times(1)

		loader, _ := NewLoader(config, sourceFactory)
		check := loader.Load(sourceID, sourcePath, sourceFormat)

		if check == nil {
			t.Errorf("%s didn't returned the expected error instance", action)
		} else {
			if strings.Index(check.Error(), expectedError) != 0 {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should return the error that the config may return on source registration", func(t *testing.T) {
		action := "Loading the config source and the config returns an error on registration"

		sourceID := "__dummy_base_source_id__"
		sourcePath := "__dummy_base_source_path__"
		sourceFormat := DecoderFormatYAML
		loadedID := "__dummy_loaded_id__"
		loadedPriority := 1
		expectedError := "__dummy_error__"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)
		loadedSource := NewMockSource(ctrl)

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(loadedID).Times(1)
		partial.EXPECT().Int("priority").Return(loadedPriority).Times(1)

		config := NewMockConfig(ctrl)
		config.EXPECT().Get("config.sources").Return([]interface{}{partial}).Times(1)
		gomock.InOrder(
			config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1),
			config.EXPECT().AddSource(loadedID, loadedPriority, loadedSource).Return(fmt.Errorf(expectedError)).Times(1),
		)

		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		sourceFactory.EXPECT().CreateConfig(partial).Return(loadedSource, nil).Times(1)

		loader, _ := NewLoader(config, sourceFactory)
		check := loader.Load(sourceID, sourcePath, sourceFormat)

		if check == nil {
			t.Errorf("%s didn't returned the expected error instance", action)
		} else {
			if strings.Index(check.Error(), expectedError) != 0 {
				t.Errorf("%s returned the (%v) error, expected (%v)", action, check.Error(), expectedError)
			}
		}
	})

	t.Run("should correctly register the loaded source", func(t *testing.T) {
		action := "Loading the config source"

		sourceID := "__dummy_base_source_id__"
		sourcePath := "__dummy_base_source_path__"
		sourceFormat := DecoderFormatYAML
		loadedID := "__dummy_loaded_id__"
		loadedPriority := 1

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)
		loadedSource := NewMockSource(ctrl)

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(loadedID).Times(1)
		partial.EXPECT().Int("priority").Return(loadedPriority).Times(1)

		config := NewMockConfig(ctrl)
		config.EXPECT().Get("config.sources").Return([]interface{}{partial}).Times(1)
		gomock.InOrder(
			config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1),
			config.EXPECT().AddSource(loadedID, loadedPriority, loadedSource).Return(nil).Times(1),
		)

		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		sourceFactory.EXPECT().CreateConfig(partial).Return(loadedSource, nil).Times(1)

		loader, _ := NewLoader(config, sourceFactory)
		check := loader.Load(sourceID, sourcePath, sourceFormat)

		if check != nil {
			t.Errorf("%s returned a unexpected error : %v", action, check)
		}
	})
}
