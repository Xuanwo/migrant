package config

import "github.com/Xuanwo/migrant/constants"

// Service stores a configuration for a service.
type Service struct {
	Type    string                 `yaml:"type"`
	Options map[string]interface{} `yaml:"options"`
}

// Check checks the configuration.
func (s *Service) Check() error {
	// TODO: we should check service config here.
	switch s.Type {
	case constants.ServiceFs:
	case constants.ServiceMySQL:
	case constants.ServiceRedis:
	case constants.ServiceRedisSentinel:
	default:
		return constants.ErrServiceInvalid
	}
	return nil
}
