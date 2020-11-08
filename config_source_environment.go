package servlet

import (
	"os"
	"strings"
	"sync"
)

// ConfigSourceEnvironment defines an instance of a environment variables
// stream configuration source.
type ConfigSourceEnvironment struct {
	ConfigSourceBase
	mapper map[string]string
}

// NewConfigSourceEnvironment instantiate a new source that read a list of
// environment variables into mapped config paths.
func NewConfigSourceEnvironment(mapper map[string]string) (*ConfigSourceEnvironment, error) {
	s := &ConfigSourceEnvironment{
		ConfigSourceBase: ConfigSourceBase{
			mutex:   &sync.Mutex{},
			partial: ConfigPartial{},
		},
		mapper: mapper,
	}

	_ = s.load()

	return s, nil
}

func (s *ConfigSourceEnvironment) load() error {
	for v, p := range s.mapper {
		if env := os.Getenv(v); env != "" {
			step := s.partial
			sections := strings.Split(p, ".")
			for i, section := range sections {
				if i != len(sections)-1 {
					if _, ok := step[section]; ok == false {
						step[section] = ConfigPartial{}
					}

					switch step[section].(type) {
					case ConfigPartial:
					default:
						step[section] = ConfigPartial{}
					}

					step = step[section].(ConfigPartial)
				} else {
					step[section] = env
				}
			}
		}
	}

	return nil
}
