package sys

import (
	"testing"
)

func Test_NewFileSystemProviderParameters(t *testing.T) {
	t.Run("should creates a new file system provider parameters", func(t *testing.T) {
		action := "Creating the new file system provider parameters"

		id := "__dummy_id__"

		parameters := NewFileSystemProviderParameters(id).(*fileSystemProviderParameters)

		if check := parameters.id; check != id {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, id, check)
		}
	})
}

func Test_NewDefaultFileSystemProviderParameters(t *testing.T) {
	t.Run("should creates a new file system provider parameters with the default values", func(t *testing.T) {
		action := "Creating the new file system provider parameters without values"

		parameters := NewDefaultFileSystemProviderParameters().(*fileSystemProviderParameters)

		if check := parameters.id; check != ContainerFileSystemID {
			t.Errorf("%s didn't store the expected (%v) file system id, returned (%v)", action, ContainerFileSystemID, check)
		}
	})
}

func Test_FileSystemProviderParameters_GetID(t *testing.T) {
	t.Run("should retrieve the assigned file system id given at creation", func(t *testing.T) {
		action := "Retrieving the file system id"

		id := "__dummy_id__"

		parameters := NewFileSystemProviderParameters(id)

		if check := parameters.GetID(); check != id {
			t.Errorf("%s didn't store the expected (%v) logger id, returned (%v)", action, id, check)
		}
	})
}

func Test_FileSystemProviderParameters_SetID(t *testing.T) {
	t.Run("should assign a new file system id", func(t *testing.T) {
		action := "Assigning the file system id"

		id := "__dummy_id__"

		parameters := NewDefaultFileSystemProviderParameters()
		parameters.SetID(id)

		if check := parameters.GetID(); check != id {
			t.Errorf("%s didn't store the expected (%v) logger id, returned (%v)", action, id, check)
		}
	})
}
