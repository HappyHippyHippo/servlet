package servlet

const (
	// ContainerFileSystemID defines the default id used to register the
	// application file system adapter instance in the application container.
	ContainerFileSystemID = "servlet.filesystem"

	// EnvContainerFileSystemID defines the environment variable used to
	// override the default value for the container file system adapter id
	EnvContainerFileSystemID = "SERVLET_CONTAINER_FILE_SYSTEM_ID"
)
