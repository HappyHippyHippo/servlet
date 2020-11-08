package servlet

const (
	// LogFormatterFormatJSON defines the value to be used to declare a JSON
	// log formatter format.
	LogFormatterFormatJSON = "json"

	// LogStreamTypeFile defines the value to be used to declare a file
	// log stream type.
	LogStreamTypeFile = "file"

	// ContainerLoggerID defines the id to be used as the default of a
	// logger instance in the application container.
	ContainerLoggerID = "servlet.log"

	// EnvContainerLoggerID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// logger id.
	EnvContainerLoggerID = "SERVLET_CONTAINER_LOGGER_ID"

	// ContainerLogFormatterFactoryStrategyJSONID defines the id to be used as
	// the default of a logger json formatter factory strategy instance in the
	// application container.
	ContainerLogFormatterFactoryStrategyJSONID = "servlet.log.factory.formatter.json"

	// EnvContainerLogFormatterFactoryStrategyJSONID defines the name of the
	// environment variable to be checked for a overriding value for the
	// application container logger json formatter factory strategy id.
	EnvContainerLogFormatterFactoryStrategyJSONID = "SERVLET_CONTAINER_LOGGER_FORMATTER_FACTORY_STRATEGY_JSON_ID"

	// ContainerLogFormatterFactoryID defines the id to be used as the
	// default of a logger formatter factory instance in the application
	// container.
	ContainerLogFormatterFactoryID = "servlet.log.factory.formatter"

	// EnvContainerLogFormatterFactoryID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container logger formatter factory id.
	EnvContainerLogFormatterFactoryID = "SERVLET_CONTAINER_LOGGER_FORMATTER_FACTORY_ID"

	// ContainerLogStreamFactoryStrategyFileID defines the id to be used as the
	// default of a logger file stream factory strategy instance in the
	// application container.
	ContainerLogStreamFactoryStrategyFileID = "servlet.log.factory.stream.file"

	// EnvContainerLogStreamFactoryStrategyFileID defines the name of the
	// environment variable to be checked for a overriding value for the
	// application container logger file stream factory strategy id.
	EnvContainerLogStreamFactoryStrategyFileID = "SERVLET_CONTAINER_LOGGER_STREAM_FACTORY_STRATEGY_FILE_ID"

	// ContainerLogStreamFactoryID defines the id to be used as the default
	// of a logger stream factory instance in the application container.
	ContainerLogStreamFactoryID = "servlet.log.factory.stream"

	// EnvContainerLogStreamFactoryID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container logger stream factory id.
	EnvContainerLogStreamFactoryID = "SERVLET_CONTAINER_LOGGER_STREAM_FACTORY_ID"

	// ContainerLogLoaderID defines the id to be used as the default of a
	// logger loader instance in the application container.
	ContainerLogLoaderID = "servlet.log.loader"

	// EnvContainerLogLoaderID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container logger loader id.
	EnvContainerLogLoaderID = "SERVLET_CONTAINER_LOGGER_LOADER_ID"
)
