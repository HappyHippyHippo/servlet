package servlet

import (
	"time"
)

const (
	// ConfigDecoderFormatYAML defines the value to be used to declare a YAML
	// config source format.
	ConfigDecoderFormatYAML = "yaml"

	// ConfigSourceTypeFile defines the value to be used to declare a
	// simple file config source type.
	ConfigSourceTypeFile = "file"

	// ConfigSourceTypeObservableFile defines the value to be used to declare a
	// observable file config source type.
	ConfigSourceTypeObservableFile = "observable_file"

	// ConfigSourceTypeEnv defines the value to be used to declare a
	// environment config source type.
	ConfigSourceTypeEnv = "env"

	// ContainerConfigID defines the id to be used as the default of a
	// config instance in the application container.
	ContainerConfigID = "servlet.config"

	// EnvContainerConfigID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// config id.
	EnvContainerConfigID = "SERVLET_CONTAINER_CONFIG_ID"

	// ContainerConfigDecoderFactoryStrategyYamlID defines the id to be used
	// as the default of a yaml config decoder factory strategy instance in
	// the application container.
	ContainerConfigDecoderFactoryStrategyYamlID = "servlet.config.factory.decoder.yaml"

	// EnvContainerConfigDecoderFactoryStrategyYamlID defines the name of
	// the environment variable to be checked for a overriding value for
	// the application container yaml config decoder factory strategy id.
	EnvContainerConfigDecoderFactoryStrategyYamlID = "SERVLET_CONTAINER_CONFIG_DECODER_FACTORY_STRATEGY_YAML_ID"

	// ContainerConfigDecoderFactoryID defines the id to be used as the
	// default of a config decoder factory instance in the application
	// container.
	ContainerConfigDecoderFactoryID = "servlet.config.factory.decoder"

	// EnvContainerConfigDecoderFactoryID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container config decoder factory id.
	EnvContainerConfigDecoderFactoryID = "SERVLET_CONTAINER_CONFIG_DECODER_FACTORY_ID"

	// ContainerConfigSourceFactoryStrategyFileID defines the id to be used as
	// the default of a config file source factory strategy instance in the
	// application container.
	ContainerConfigSourceFactoryStrategyFileID = "servlet.config.factory.source.file"

	// EnvContainerConfigSourceFactoryStrategyFileID defines the name of the
	// environment variable to be checked for a overriding value for the
	// application container config file source factory strategy id.
	EnvContainerConfigSourceFactoryStrategyFileID = "SERVLET_CONTAINER_CONFIG_SOURCE_FACTORY_STRATEGY_FILE_ID"

	// ContainerConfigSourceFactoryStrategyObservableFileID defines the id to
	// the default of a config observable file source factory strategy instance
	// in the application container.
	ContainerConfigSourceFactoryStrategyObservableFileID = "servlet.config.factory.source.observable_file"

	// EnvContainerConfigSourceFactoryStrategyObservableFileID defines the name
	// of the environment variable to be checked for a overriding value for the
	// application container config observable file source factory strategy id.
	EnvContainerConfigSourceFactoryStrategyObservableFileID = "SERVLET_CONTAINER_CONFIG_SOURCE_FACTORY_STRATEGY_OBSERVABLE_FILE_ID"

	// ContainerConfigSourceFactoryStrategyEnvironmentID defines the id to the default
	// of a config environment source factory strategy instance in the
	// application container.
	ContainerConfigSourceFactoryStrategyEnvironmentID = "servlet.config.factory.source.environment"

	// EnvContainerConfigSourceFactoryStrategyEnvironmentID defines the name of the
	// environment variable to be checked for a overriding value for the
	// application container config environment source factory strategy id.
	EnvContainerConfigSourceFactoryStrategyEnvironmentID = "SERVLET_CONTAINER_CONFIG_SOURCE_FACTORY_STRATEGY_ENVIRONMENT_ID"

	// ContainerConfigSourceFactoryID defines the id to be used as the default
	// of a config source factory instance in the application container.
	ContainerConfigSourceFactoryID = "servlet.config.factory.source"

	// EnvContainerConfigSourceFactoryID defines the name of the environment
	// variable to be checked for a overriding value for the application
	// container config source factory id.
	EnvContainerConfigSourceFactoryID = "SERVLET_CONTAINER_CONFIG_SOURCE_FACTORY_ID"

	// ContainerConfigLoaderID defines the id to be used as the default of a
	// config loader instance in the application container.
	ContainerConfigLoaderID = "servlet.config.loader"

	// EnvContainerConfigLoaderID defines the name of the environment variable
	// to be checked for a overriding value for the application container
	// config loaded id.
	EnvContainerConfigLoaderID = "SERVLET_CONTAINER_CONFIG_LOADER_ID"

	// ConfigObserveFrequency defines the id to be used as the default of a
	// config observable source frequency time.
	ConfigObserveFrequency = time.Second * 0

	// EnvConfigObserveFrequency defines the name of the environment variable
	// to be checked for a overriding value for the config observe frequency.
	EnvConfigObserveFrequency = "SERVLET_CONFIG_OBSERVE_FREQUENCY"

	// ConfigEntrySourceActive defines the entry config source active flag
	// used to signal the config loader to load the entry source or not
	ConfigEntrySourceActive = true

	// EnvConfigEntrySourceActive defines the name of the environment variable
	// to be checked for a overriding value for the config entry source active.
	EnvConfigEntrySourceActive = "SERVLET_CONFIG_ENTRY_SOURCE_ACTIVE"

	// ConfigEntrySourceID defines the id to be used as the default of the
	// entry config source id to be used as the loader entry.
	ConfigEntrySourceID = "entry"

	// EnvConfigEntrySourceID defines the name of the environment variable
	// to be checked for a overriding value for the config entry source id.
	EnvConfigEntrySourceID = "SERVLET_CONFIG_ENTRY_SOURCE_ID"

	// ConfigEntrySourcePath defines the entry config source path
	// to be used as the loader entry.
	ConfigEntrySourcePath = "config/config.yaml"

	// EnvConfigEntrySourcePath defines the name of the environment variable
	// to be checked for a overriding value for the config entry source path.
	EnvConfigEntrySourcePath = "SERVLET_CONFIG_ENTRY_SOURCE_PATH"

	// ConfigEntrySourceFormat defines the entry config source format
	// to be used as the loader entry.
	ConfigEntrySourceFormat = ConfigDecoderFormatYAML

	// EnvConfigEntrySourceFormat defines the name of the environment variable
	// to be checked for a overriding value for the config entry source format.
	EnvConfigEntrySourceFormat = "SERVLET_CONFIG_ENTRY_SOURCE_FORMAT"
)
