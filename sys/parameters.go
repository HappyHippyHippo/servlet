package sys

import "os"

const (
	// ContainerFileSystemID defines the default id used to register the
	// application file system adapter instance in the application container.
	ContainerFileSystemID = "servlet.filesystem"

	// EnvContainerFileSystemID defines the environment variable used to
	// override the default value for the container file system adapter id
	EnvContainerFileSystemID = "SERVLET_CONTAINER_FILE_SYSTEM_ID"
)

// Parameters defines the system provider parameters storing structure
// that will be needed when instantiating a new provider
type Parameters struct {
	FileSystemID string
}

// NewParameters will instantiate a new file system provider
// parameters object with the requested values.
func NewParameters(fileSystemID string) Parameters {
	return Parameters{
		FileSystemID: fileSystemID,
	}
}

// NewDefaultParameters instantiate a new file system
// provider parameters object with the default values.
func NewDefaultParameters() Parameters {
	fileSystemID := ContainerFileSystemID
	if env := os.Getenv(EnvContainerFileSystemID); env != "" {
		fileSystemID = env
	}

	return Parameters{
		FileSystemID: fileSystemID,
	}
}
