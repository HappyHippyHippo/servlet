package servlet

// ConfigSourceFactoryStrategyEnvironment defines a environment config source
// instantiation strategy to be used by the config sources factory
// instance.
type ConfigSourceFactoryStrategyEnvironment struct{}

// NewConfigSourceFactoryStrategyEnvironment instantiate a new environment
// source factory strategy that will enable the source factory to instantiate
// a new observable file configuration source.
func NewConfigSourceFactoryStrategyEnvironment() (*ConfigSourceFactoryStrategyEnvironment, error) {
	return &ConfigSourceFactoryStrategyEnvironment{}, nil
}

// Accept will check if the source factory strategy can instantiate a
// new source of the requested type. Also, validates that there is the path
// and content format extra parameters, and thar this parameters are strings.
func (ConfigSourceFactoryStrategyEnvironment) Accept(sourceType string, args ...interface{}) bool {
	if sourceType != ConfigSourceTypeEnv || len(args) < 1 {
		return false
	}

	switch args[0].(type) {
	case map[string]string:
	default:
		return false
	}

	return true
}

// AcceptConfig will check if the source factory strategy can instantiate a
// source where the data to check comes from a configuration partial instance.
func (s ConfigSourceFactoryStrategyEnvironment) AcceptConfig(conf ConfigPartial) (check bool) {
	defer func() {
		if r := recover(); r != nil {
			check = false
		}
	}()

	sourceType := conf.String("type")
	mapping := conf.Get("mapping")

	return s.Accept(sourceType, mapping)
}

// Create will instantiate the desired environment source instance.
func (s ConfigSourceFactoryStrategyEnvironment) Create(args ...interface{}) (source ConfigSource, err error) {
	defer func() {
		if r := recover(); r != nil {
			source = nil
			err = r.(error)
		}
	}()

	mappings := args[0].(map[string]string)

	return NewConfigSourceEnvironment(mappings)
}

// CreateConfig will instantiate the desired environment source instance
// where the initialization data comes from a configuration partial instance.
func (s ConfigSourceFactoryStrategyEnvironment) CreateConfig(conf ConfigPartial) (source ConfigSource, err error) {
	defer func() {
		if r := recover(); r != nil {
			source = nil
			err = r.(error)
		}
	}()

	mapping := map[string]string{}
	for k, v := range conf.Get("mapping").(ConfigPartial) {
		mapping[k.(string)] = v.(string)
	}

	return s.Create(mapping)
}
