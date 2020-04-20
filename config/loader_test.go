package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func Test_NewLoader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	config := NewMockConfig(ctrl)
	sourceFactory := NewMockSourceFactory(ctrl)

	t.Run("error when missing the config", func(t *testing.T) {
		if loader, err := NewLoader(nil, sourceFactory); loader != nil {
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'config' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error when missing the config source factory", func(t *testing.T) {
		if loader, err := NewLoader(config, nil); loader != nil {
			t.Errorf("returned a valid reference")
		} else if err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Invalid nil 'sourceFactory' argument" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("should return the config loader", func(t *testing.T) {
		if loader, err := NewLoader(config, sourceFactory); loader == nil {
			t.Errorf("didn't return a valid reference")
		} else if err != nil {
			t.Errorf("return the (%v) error", err)
		}
	})
}

func Test_Loader_Load(t *testing.T) {
	sourceID := "base_source_id"
	sourcePath := "base_source_path"
	sourceFormat := DecoderFormatYAML
	loadedID := "loaded_id"
	loadedPriority := 1
	expectedError := "Error while parsing the config entry : "

	t.Run("error getting the base source", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(nil, fmt.Errorf(expectedError)).Times(1)
		config := NewMockConfig(ctrl)

		loader, _ := NewLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error storing the base source", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, baseSource).Return(fmt.Errorf(expectedError)).Times(1)

		loader, _ := NewLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("add base source into the config", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return(nil).Times(1)

		loader, _ := NewLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("no error if list of sources isn't present", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return([]interface{}{}).Times(1)

		loader, _ := NewLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on invalid list of sources", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return("invalid_list").Times(1)

		loader, _ := NewLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != "Error while parsing the list of sources" {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on loaded invalid id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").DoAndReturn(func(key string) string {
			panic("invalid convertion")
		}).Times(1)
		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return([]interface{}{partial}).Times(1)

		loader, _ := NewLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Errorf("didn't return the expected error")
		} else if strings.Index(err.Error(), expectedError) != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on loaded invalid priority", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		baseSource := NewMockSource(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(loadedID).Times(1)
		partial.EXPECT().Int("priority").DoAndReturn(func(key string) string {
			panic("invalid convertion")
		}).Times(1)
		config := NewMockConfig(ctrl)
		config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1)
		config.EXPECT().Get("config.sources").Return([]interface{}{partial}).Times(1)

		loader, _ := NewLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Errorf("didn't return the expected error")
		} else if strings.Index(err.Error(), expectedError) != 0 {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on loaded source factory", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(loadedID).Times(1)
		partial.EXPECT().Int("priority").Return(loadedPriority).Times(1)
		baseSource := NewMockSource(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		sourceFactory.EXPECT().CreateConfig(partial).Return(nil, fmt.Errorf(expectedError)).Times(1)
		config := NewMockConfig(ctrl)
		config.EXPECT().Get("config.sources").Return([]interface{}{partial}).Times(1)
		config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1)

		loader, _ := NewLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("error on source registration", func(t *testing.T) {
		expectedError := "error"

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(loadedID).Times(1)
		partial.EXPECT().Int("priority").Return(loadedPriority).Times(1)
		baseSource := NewMockSource(ctrl)
		loadedSource := NewMockSource(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		sourceFactory.EXPECT().CreateConfig(partial).Return(loadedSource, nil).Times(1)
		config := NewMockConfig(ctrl)
		config.EXPECT().Get("config.sources").Return([]interface{}{partial}).Times(1)
		gomock.InOrder(
			config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1),
			config.EXPECT().AddSource(loadedID, loadedPriority, loadedSource).Return(fmt.Errorf(expectedError)).Times(1),
		)

		loader, _ := NewLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err == nil {
			t.Errorf("didn't return the expected error")
		} else if err.Error() != expectedError {
			t.Errorf("returned the (%v) error", err)
		}
	})

	t.Run("register the loaded source", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		partial := NewMockPartial(ctrl)
		partial.EXPECT().String("id").Return(loadedID).Times(1)
		partial.EXPECT().Int("priority").Return(loadedPriority).Times(1)
		baseSource := NewMockSource(ctrl)
		loadedSource := NewMockSource(ctrl)
		sourceFactory := NewMockSourceFactory(ctrl)
		sourceFactory.EXPECT().Create(SourceTypeFile, sourcePath, sourceFormat).Return(baseSource, nil).Times(1)
		sourceFactory.EXPECT().CreateConfig(partial).Return(loadedSource, nil).Times(1)
		config := NewMockConfig(ctrl)
		config.EXPECT().Get("config.sources").Return([]interface{}{partial}).Times(1)
		gomock.InOrder(
			config.EXPECT().AddSource(sourceID, 0, baseSource).Return(nil).Times(1),
			config.EXPECT().AddSource(loadedID, loadedPriority, loadedSource).Return(nil).Times(1),
		)

		loader, _ := NewLoader(config, sourceFactory)

		if err := loader.Load(sourceID, sourcePath, sourceFormat); err != nil {
			t.Errorf("returned the (%v) error", err)
		}
	})
}
