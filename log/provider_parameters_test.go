package log

import (
	"reflect"
	"testing"

	"github.com/happyhippyhippo/servlet/config"
	"github.com/happyhippyhippo/servlet/sys"
)

func Test_NewProviderParameters(t *testing.T) {
	t.Run("should creates a new logger provider parameters", func(t *testing.T) {
		action := "Creating the new logger provider parameters"

		id := "__dummy_id__"
		fileSystemID := "__dummy_file_system_id__"
		configID := "__dummy_config_id__"
		formatterFactoryID := "__dummy_formatter_factory_id__"
		streamFactoryID := "__dummy_stream_factory_id__"
		loaderID := "__dummy_loader_id__"

		parameters := NewProviderParameters(
			id,
			fileSystemID,
			configID,
			formatterFactoryID,
			streamFactoryID,
			loaderID,
		).(*providerParameters)

		if check := parameters.id; check != id {
			t.Errorf("%s didn't store the expected (%v) logger id, returned (%v)", action, id, check)
		}
		if check := parameters.fileSystemID; check != fileSystemID {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, fileSystemID, check)
		}
		if check := parameters.configID; check != configID {
			t.Errorf("%s didn't store the expected (%v) config id, returned (%v)", action, configID, check)
		}
		if check := parameters.formatterFactoryID; check != formatterFactoryID {
			t.Errorf("%s didn't store the expected (%v) formatter factory id, returned (%v)", action, formatterFactoryID, check)
		}
		if check := parameters.streamFactoryID; check != streamFactoryID {
			t.Errorf("%s didn't store the expected (%v) stream factory id, returned (%v)", action, streamFactoryID, check)
		}
		if check := parameters.loaderID; check != loaderID {
			t.Errorf("%s didn't store the expected (%v) loader id, returned (%v)", action, loaderID, check)
		}
	})
}

func Test_NewDefaultProviderParameters(t *testing.T) {
	t.Run("should creates a new logger provider parameters with the default values", func(t *testing.T) {
		action := "Creating the new logger provider parameters without values"

		parameters := NewDefaultProviderParameters().(*providerParameters)

		if check := parameters.id; check != ContainerID {
			t.Errorf("%s didn't store the expected (%v) logger id, returned (%v)", action, ContainerID, check)
		}
		if check := parameters.fileSystemID; check != sys.ContainerFileSystemID {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, sys.ContainerFileSystemID, check)
		}
		if check := parameters.configID; check != config.ContainerID {
			t.Errorf("%s didn't store the expected (%v) config id, returned (%v)", action, config.ContainerID, check)
		}
		if check := parameters.formatterFactoryID; check != ContainerFormatterFactoryID {
			t.Errorf("%s didn't store the expected (%v) formatter factory id, returned (%v)", action, ContainerFormatterFactoryID, check)
		}
		if check := parameters.streamFactoryID; check != ContainerStreamFactoryID {
			t.Errorf("%s didn't store the expected (%v) stream factory id, returned (%v)", action, ContainerStreamFactoryID, check)
		}
		if check := parameters.loaderID; check != ContainerLoaderID {
			t.Errorf("%s didn't store the expected (%v) loader id, returned (%v)", action, ContainerLoaderID, check)
		}
	})
}

func Test_ProviderParameters_GetID(t *testing.T) {
	t.Run("should correctly retrieve the stored logger id", func(t *testing.T) {
		action := "Retrieving the logger id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters(expected, "", "", "", "", "")

		if check := parameters.GetID(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetID(t *testing.T) {
	t.Run("should assign a new logger id", func(t *testing.T) {
		action := "Assigning the logger id"

		id := "__dummy_id__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetID(); check != id {
			t.Errorf("%s didn't store the expected (%v) logger id, returned (%v)", action, id, check)
		}
	})
}

func Test_ProviderParameters_GetFileSystemID(t *testing.T) {
	t.Run("should correctly retrieve the stored file system id", func(t *testing.T) {
		action := "Retrieving the file system id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", expected, "", "", "", "")

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

func Test_ProviderParameters_GetConfigID(t *testing.T) {
	t.Run("should correctly retrieve the stored config id", func(t *testing.T) {
		action := "Retrieving the config id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", "", expected, "", "", "")

		if check := parameters.GetConfigID(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetConfigID(t *testing.T) {
	t.Run("should assign a new config id", func(t *testing.T) {
		action := "Assigning the config id"

		id := "__dummy_id__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetConfigID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetConfigID(); check != id {
			t.Errorf("%s didn't store the expected (%v) config id, returned (%v)", action, id, check)
		}
	})
}

func Test_ProviderParameters_GetFormatterFactoryID(t *testing.T) {
	t.Run("should correctly retrieve the stored formatter factory id", func(t *testing.T) {
		action := "Retrieving the formatter factory id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", "", "", expected, "", "")

		if check := parameters.GetFormatterFactoryID(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetFormatterFactoryID(t *testing.T) {
	t.Run("should assign a new formatter factory id", func(t *testing.T) {
		action := "Assigning the formatter factory id"

		id := "__dummy_id__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetFormatterFactoryID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetFormatterFactoryID(); check != id {
			t.Errorf("%s didn't store the expected (%v) formatter factory id, returned (%v)", action, id, check)
		}
	})
}

func Test_ProviderParameters_GetStreamFactoryID(t *testing.T) {
	t.Run("should correctly retrieve the stored stream factory id", func(t *testing.T) {
		action := "Retrieving the stream factory id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", "", "", "", expected, "")

		if check := parameters.GetStreamFactoryID(); check != expected {
			t.Errorf("%s returned (%v), expected (%v)", action, check, expected)
		}
	})
}

func Test_ProviderParameters_SetStreamFactoryID(t *testing.T) {
	t.Run("should assign a new stream factory id", func(t *testing.T) {
		action := "Assigning the stream factory id"

		id := "__dummy_id__"

		parameters := NewDefaultProviderParameters()
		if check := parameters.SetStreamFactoryID(id); !reflect.DeepEqual(check, parameters) {
			t.Errorf("%s didn't returned the parameters instance", action)
		}

		if check := parameters.GetStreamFactoryID(); check != id {
			t.Errorf("%s didn't store the expected (%v) stream factory id, returned (%v)", action, id, check)
		}
	})
}

func Test_ProviderParameters_GetLoaderID(t *testing.T) {
	t.Run("should correctly retrieve the stored loader id", func(t *testing.T) {
		action := "Retrieving the loader id"

		expected := "__dummy_value__"

		parameters := NewProviderParameters("", "", "", "", "", expected)

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
