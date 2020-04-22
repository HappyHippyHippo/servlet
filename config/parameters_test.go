package config

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/happyhippyhippo/servlet/sys"
)

func Test_NewParameters(t *testing.T) {
	configID := "config_id"
	fileSystemID := "file_system_id"
	sourceFactoryID := "source_factory_id"
	decoderFactoryID := "decoder_factory_id"
	loaderID := "loader_id"
	observeFrequency := time.Second
	baseSourceID := "base_config_id"
	baseSourcePath := "base_config_path"
	baseSourceFormat := "base_config_format"

	t.Run("creates a new parameters", func(t *testing.T) {
		parameters := NewParameters(
			configID,
			fileSystemID,
			sourceFactoryID,
			decoderFactoryID,
			loaderID,
			observeFrequency,
			baseSourceID,
			baseSourcePath,
			baseSourceFormat,
		)

		if value := parameters.ConfigID; value != configID {
			t.Errorf("stored (%v) config ID", value)
		} else if value := parameters.FileSystemID; value != fileSystemID {
			t.Errorf("stored (%v) file sytem ID", value)
		} else if value := parameters.SourceFactoryID; value != sourceFactoryID {
			t.Errorf("stored (%v) source factory ID", value)
		} else if value := parameters.DecoderFactoryID; value != decoderFactoryID {
			t.Errorf("stored (%v) decoder factory ID", value)
		} else if value := parameters.LoaderID; value != loaderID {
			t.Errorf("stored (%v) loader ID", value)
		} else if value := parameters.ObserveFrequency; value != observeFrequency {
			t.Errorf("stored (%v) observe frequecy", value)
		} else if value := parameters.BaseSourceID; value != baseSourceID {
			t.Errorf("stored (%v) base source id", value)
		} else if value := parameters.BaseSourcePath; value != baseSourcePath {
			t.Errorf("stored (%v) base source path", value)
		} else if value := parameters.BaseSourceFormat; value != baseSourceFormat {
			t.Errorf("stored (%v) base source format", value)
		}
	})
}

func Test_NewDefaultParameters(t *testing.T) {
	t.Run("creates a new parameters with the default values", func(t *testing.T) {
		parameters := NewDefaultParameters()

		if value := parameters.ConfigID; value != ContainerConfigID {
			t.Errorf("stored (%v) config ID", value)
		} else if value := parameters.FileSystemID; value != sys.ContainerFileSystemID {
			t.Errorf("stored (%v) file sytem ID", value)
		} else if value := parameters.SourceFactoryID; value != ContainerSourceFactoryID {
			t.Errorf("stored (%v) source factory ID", value)
		} else if value := parameters.DecoderFactoryID; value != ContainerDecoderFactoryID {
			t.Errorf("stored (%v) decoder factory ID", value)
		} else if value := parameters.LoaderID; value != ContainerLoaderID {
			t.Errorf("stored (%v) loader ID", value)
		} else if value := parameters.ObserveFrequency; value != ObserveFrequency {
			t.Errorf("stored (%v) observe frequecy", value)
		} else if value := parameters.BaseSourceID; value != BaseSourceID {
			t.Errorf("stored (%v) base source id", value)
		} else if value := parameters.BaseSourcePath; value != BaseSourcePath {
			t.Errorf("stored (%v) base source path", value)
		} else if value := parameters.BaseSourceFormat; value != BaseSourceFormat {
			t.Errorf("stored (%v) base source format", value)
		}
	})

	configID := "config_id"
	fileSystemID := "file_system_id"
	sourceFactoryID := "source_factory_id"
	decoderFactoryID := "decoder_factory_id"
	loaderID := "loader_id"
	observeFrequency := time.Second * 10
	baseSourceID := "base_config_id"
	baseSourcePath := "base_config_path"
	baseSourceFormat := "base_config_format"

	t.Run("creates a new parameters with the env config ID", func(t *testing.T) {
		os.Setenv(EnvContainerConfigID, configID)
		defer os.Setenv(EnvContainerConfigID, "")

		parameters := NewDefaultParameters()
		if value := parameters.ConfigID; value != configID {
			t.Errorf("stored (%v) config ID", value)
		}
	})

	t.Run("creates a new parameters with the env file system ID", func(t *testing.T) {
		os.Setenv(sys.EnvContainerFileSystemID, fileSystemID)
		defer os.Setenv(sys.EnvContainerFileSystemID, "")

		parameters := NewDefaultParameters()
		if value := parameters.FileSystemID; value != fileSystemID {
			t.Errorf("stored (%v) file system ID", value)
		}
	})

	t.Run("creates a new parameters with the env source factory ID", func(t *testing.T) {
		os.Setenv(EnvContainerSourceFactoryID, sourceFactoryID)
		defer os.Setenv(EnvContainerSourceFactoryID, "")

		parameters := NewDefaultParameters()
		if value := parameters.SourceFactoryID; value != sourceFactoryID {
			t.Errorf("stored (%v) source factory ID", value)
		}
	})

	t.Run("creates a new parameters with the env decoder factory ID", func(t *testing.T) {
		os.Setenv(EnvContainerDecoderFactoryID, decoderFactoryID)
		defer os.Setenv(EnvContainerDecoderFactoryID, "")

		parameters := NewDefaultParameters()
		if value := parameters.DecoderFactoryID; value != decoderFactoryID {
			t.Errorf("stored (%v) decoder factory ID", value)
		}
	})

	t.Run("creates a new parameters with the env loader ID", func(t *testing.T) {
		os.Setenv(EnvContainerLoaderID, loaderID)
		defer os.Setenv(EnvContainerLoaderID, "")

		parameters := NewDefaultParameters()
		if value := parameters.LoaderID; value != loaderID {
			t.Errorf("stored (%v) loader ID", value)
		}
	})

	t.Run("creates a new parameters with the env observer frequency", func(t *testing.T) {
		os.Setenv(EnvObserveFrequency, strconv.Itoa(int(observeFrequency.Seconds())))
		defer os.Setenv(EnvObserveFrequency, "")

		parameters := NewDefaultParameters()
		if value := parameters.ObserveFrequency; value != observeFrequency {
			t.Errorf("stored (%v) observe frequency", value)
		}
	})

	t.Run("creates a new parameters with the base source ID", func(t *testing.T) {
		os.Setenv(EnvBaseSourceID, baseSourceID)
		defer os.Setenv(EnvBaseSourceID, "")

		parameters := NewDefaultParameters()
		if value := parameters.BaseSourceID; value != baseSourceID {
			t.Errorf("stored (%v) base source ID", value)
		}
	})

	t.Run("creates a new parameters with the env base source path", func(t *testing.T) {
		os.Setenv(EnvBaseSourcePath, baseSourcePath)
		defer os.Setenv(EnvBaseSourcePath, "")

		parameters := NewDefaultParameters()
		if value := parameters.BaseSourcePath; value != baseSourcePath {
			t.Errorf("stored (%v) base source path", value)
		}
	})

	t.Run("creates a new parameters with the env base source format", func(t *testing.T) {
		os.Setenv(EnvBaseSourceFormat, baseSourceFormat)
		defer os.Setenv(EnvBaseSourceFormat, "")

		parameters := NewDefaultParameters()
		if value := parameters.BaseSourceFormat; value != baseSourceFormat {
			t.Errorf("stored (%v) base source format", value)
		}
	})
}
