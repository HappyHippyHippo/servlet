package config

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/happyhippyhippo/servlet/sys"
)

func Test_NewProviderParameters(t *testing.T) {
	t.Run("should creates a new config provider parameters", func(t *testing.T) {
		action := "Creating the new config provider parameters"

		id := "__dummy_id__"
		fileSystemID := "__dummy_file_system_id__"
		loaderID := "__dummy_loader_id__"
		sourceFactoryID := "__dummy_source_factory_id__"
		decoderFactoryID := "__dummy_decoder_factory_id__"
		observeFrequency := time.Second
		baseSourceID := "__dummy_base_config_id__"
		baseSourcePath := "__dummy_base_config_path__"
		baseSourceFormat := "__dummy_base_config_format__"

		parameters := NewProviderParameters(
			id,
			fileSystemID,
			loaderID,
			sourceFactoryID,
			decoderFactoryID,
			observeFrequency,
			baseSourceID,
			baseSourcePath,
			baseSourceFormat,
		).(*providerParameters)

		if check := parameters.id; check != id {
			t.Errorf("%s didn't store the expected (%v) config id, returned (%v)", action, id, check)
		}
		if check := parameters.fileSystemID; check != fileSystemID {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, fileSystemID, check)
		}
		if check := parameters.loaderID; check != loaderID {
			t.Errorf("%s didn't store the expected (%v) loader id, returned (%v)", action, loaderID, check)
		}
		if check := parameters.sourceFactoryID; check != sourceFactoryID {
			t.Errorf("%s didn't store the expected (%v) source factory id, returned (%v)", action, sourceFactoryID, check)
		}
		if check := parameters.decoderFactoryID; check != decoderFactoryID {
			t.Errorf("%s didn't store the expected (%v) decoder factory id, returned (%v)", action, decoderFactoryID, check)
		}
		if check := parameters.observeFrequency; check != observeFrequency {
			t.Errorf("%s didn't store the expected (%v) observe frequency, returned (%v)", action, observeFrequency, check)
		}
		if check := parameters.baseSourceID; check != baseSourceID {
			t.Errorf("%s didn't store the expected (%v) base config id, returned (%v)", action, baseSourceID, check)
		}
		if check := parameters.baseSourcePath; check != baseSourcePath {
			t.Errorf("%s didn't store the expected (%v) base config path, returned (%v)", action, baseSourcePath, check)
		}
		if check := parameters.baseSourceFormat; check != baseSourceFormat {
			t.Errorf("%s didn't store the expected (%v) base config format, returned (%v)", action, baseSourceFormat, check)
		}
	})
}

func Test_NewDefaultProviderParameters(t *testing.T) {
	t.Run("should creates a new config provider parameters with the default values", func(t *testing.T) {
		action := "Creating the new config provider parameters without values"

		parameters := NewDefaultProviderParameters().(*providerParameters)

		if check := parameters.id; check != ContainerID {
			t.Errorf("%s didn't store the expected (%v) config id, returned (%v)", action, ContainerID, check)
		}
		if check := parameters.fileSystemID; check != sys.ContainerFileSystemID {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, sys.ContainerFileSystemID, check)
		}
		if check := parameters.loaderID; check != ContainerLoaderID {
			t.Errorf("%s didn't store the expected (%v) loader id, returned (%v)", action, ContainerLoaderID, check)
		}
		if check := parameters.sourceFactoryID; check != ContainerSourceFactoryID {
			t.Errorf("%s didn't store the expected (%v) source factory id, returned (%v)", action, ContainerSourceFactoryID, check)
		}
		if check := parameters.decoderFactoryID; check != ContainerDecoderFactoryID {
			t.Errorf("%s didn't store the expected (%v) decoder factory id, returned (%v)", action, ContainerDecoderFactoryID, check)
		}
		if check := parameters.observeFrequency; check != ContainerObserveFrequency {
			t.Errorf("%s didn't store the expected (%v) observe frequency, returned (%v)", action, ContainerObserveFrequency, check)
		}
		if check := parameters.baseSourceID; check != ContainerBaseSourceID {
			t.Errorf("%s didn't store the expected (%v) base config id, returned (%v)", action, ContainerBaseSourceID, check)
		}
		if check := parameters.baseSourcePath; check != ContainerBaseSourcePath {
			t.Errorf("%s didn't store the expected (%v) base config path, returned (%v)", action, ContainerBaseSourcePath, check)
		}
		if check := parameters.baseSourceFormat; check != ContainerBaseSourceFormat {
			t.Errorf("%s didn't store the expected (%v) base config format, returned (%v)", action, ContainerBaseSourceFormat, check)
		}
	})

	t.Run("should creates a new config provider parameters with the base source path if defined in environment", func(t *testing.T) {
		action := "Creating the new config provider parameters without values and the base source path defined in environment"

		path := "__dummy_path__"
		os.Setenv(EnvironmentBaseSourcePath, path)
		defer os.Setenv(EnvironmentBaseSourcePath, "")

		parameters := NewDefaultProviderParameters().(*providerParameters)

		if check := parameters.baseSourcePath; check != path {
			t.Errorf("%s didn't store the expected (%v) base config path, returned (%v)", action, path, check)
		}
	})

	t.Run("should creates a new config provider parameters with the base source format if defined in environment", func(t *testing.T) {
		action := "Creating the new config provider parameters without values and the base source format defined in environment"

		format := "__dummy_format__"
		os.Setenv(EnvironmentBaseSourceFormat, format)
		defer os.Setenv(EnvironmentBaseSourceFormat, "")

		parameters := NewDefaultProviderParameters().(*providerParameters)

		if check := parameters.baseSourceFormat; check != format {
			t.Errorf("%s didn't store the expected (%v) base config format, returned (%v)", action, format, check)
		}
	})
}

func Test_ProviderParameters_GetID(t *testing.T) {
	t.Run("should correctly retrieve the stored config id", func(t *testing.T) {
		action := "Retrieving the config id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters(expected, "", "", "", "", 0, "", "", "")

		if check := parameters.GetID(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetID(t *testing.T) {
	t.Run("should assign a new config id", func(t *testing.T) {
		action := "Assigning the config id"

		id := "__dummy_id__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetID(); check != id {
			t.Errorf("%s didn't store the expected (%v) config id, returned (%v)", action, id, check)
		}
	})
}

func Test_ProviderParameters_GetFileSystemID(t *testing.T) {
	t.Run("should correctly retrieve the stored file system id", func(t *testing.T) {
		action := "Retrieving the file system id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", expected, "", "", "", 0, "", "", "")

		if check := parameters.GetFileSystemID(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetFileSystemID(t *testing.T) {
	t.Run("should assign a new file system id", func(t *testing.T) {
		action := "Assigning the file system id"

		id := "__dummy_id__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetFileSystemID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetFileSystemID(); check != id {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, id, check)
		}
	})
}

func Test_ProviderParameters_GetLoaderID(t *testing.T) {
	t.Run("should correctly retrieve the stored loader id", func(t *testing.T) {
		action := "Retrieving the loader id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", "", expected, "", "", 0, "", "", "")

		if check := parameters.GetLoaderID(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetLoaderID(t *testing.T) {
	t.Run("should assign a new loader id", func(t *testing.T) {
		action := "Assigning the loader id"

		id := "__dummy_id__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetLoaderID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetLoaderID(); check != id {
			t.Errorf("%s didn't store the expected (%v) loader id, returned (%v)", action, id, check)
		}
	})
}

func Test_ProviderParameters_GetSourceFactoryID(t *testing.T) {
	t.Run("should correctly retrieve the stored source factory id", func(t *testing.T) {
		action := "Retrieving the source factory id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", "", "", expected, "", 0, "", "", "")

		if check := parameters.GetSourceFactoryID(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetSourceFactoryID(t *testing.T) {
	t.Run("should assign a new source factory id", func(t *testing.T) {
		action := "Assigning the source factory id"

		id := "__dummy_id__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetSourceFactoryID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetSourceFactoryID(); check != id {
			t.Errorf("%s didn't store the expected (%v) source factory id, returned (%v)", action, id, check)
		}
	})
}

func Test_ProviderParameters_GetDecoderFactoryID(t *testing.T) {
	t.Run("should correctly retrieve the stored decoder factory id", func(t *testing.T) {
		action := "Retrieving the decoder factory id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", "", "", "", expected, 0, "", "", "")

		if check := parameters.GetDecoderFactoryID(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetDecoderFactoryID(t *testing.T) {
	t.Run("should assign a new decoder factory id", func(t *testing.T) {
		action := "Assigning the decoder factory id"

		id := "__dummy_id__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetDecoderFactoryID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetDecoderFactoryID(); check != id {
			t.Errorf("%s didn't store the expected (%v) decoder factory id, returned (%v)", action, id, check)
		}
	})
}

func Test_ProviderParameters_GetObserveFrequency(t *testing.T) {
	t.Run("should correctly retrieve the stored observe frequency value", func(t *testing.T) {
		action := "Retrieving the observe frequency"

		frequency := time.Second * 123

		parameters := NewProviderParameters("", "", "", "", "", frequency, "", "", "")

		if check := parameters.GetObserveFrequency(); check != frequency {
			t.Errorf("%s returned (%v), expected (%v)", action, check, frequency)
		}
	})
}

func Test_ProviderParameters_SetObserveFrequency(t *testing.T) {
	t.Run("should assign a new observe frequency", func(t *testing.T) {
		action := "Assigning the observe frequency"

		frequency := time.Second * 123

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetObserveFrequency(frequency); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetObserveFrequency(); check != frequency {
			t.Errorf("%s didn't store the expected (%v) observe frequency, returned (%v)", action, frequency, check)
		}
	})
}

func Test_ProviderParameters_GetBaseSourceID(t *testing.T) {
	t.Run("should correctly retrieve the stored base source id", func(t *testing.T) {
		action := "Retrieving the base source id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", "", "", "", "", 0, expected, "", "")

		if check := parameters.GetBaseSourceID(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetBaseSourceID(t *testing.T) {
	t.Run("should assign a new base source id", func(t *testing.T) {
		action := "Assigning the base source id"

		id := "__dummy_id__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetBaseSourceID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetBaseSourceID(); check != id {
			t.Errorf("%s didn't store the expected (%v) base source id, returned (%v)", action, id, check)
		}
	})
}

func Test_ProviderParameters_GetBaseSourcePath(t *testing.T) {
	t.Run("should correctly retrieve the stored base source path", func(t *testing.T) {
		action := "Retrieving the base source path"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", "", "", "", "", 0, "", expected, "")

		if check := parameters.GetBaseSourcePath(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetBaseSourcePath(t *testing.T) {
	t.Run("should assign a new base source path", func(t *testing.T) {
		action := "Assigning the base source path"

		path := "__dummy_path__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetBaseSourcePath(path); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetBaseSourcePath(); check != path {
			t.Errorf("%s didn't store the expected (%v) base source path, returned (%v)", action, path, check)
		}
	})
}

func Test_ProviderParameters_GetBaseSourceFormat(t *testing.T) {
	t.Run("should correctly retrieve the stored base source format", func(t *testing.T) {
		action := "Retrieving the base source format"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", "", "", "", "", 0, "", "", expected)

		if check := parameters.GetBaseSourceFormat(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetBaseSourceFormat(t *testing.T) {
	t.Run("should assign a new base source format", func(t *testing.T) {
		action := "Assigning the base source format"

		format := "__dummy_format__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetBaseSourceFormat(format); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetBaseSourceFormat(); check != format {
			t.Errorf("%s didn't store the expected (%v) base source format, returned (%v)", action, format, check)
		}
	})
}
