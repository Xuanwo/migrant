package redis

import "fmt"

type redisConfig struct {
	Host    string `yaml:"host"`
	Port    int    `yaml:"port"`
	DB      int    `yaml:"db"`
	Timeout int    `yaml:"timeout"`
}

func (r *redisConfig) Check() error {
	if r.Host == "" {
		return fmt.Errorf("host not specified for redis")
	}
	if r.Port == 0 {
		return fmt.Errorf("post not specified for redis")
	}
	if r.Timeout == 0 {
		return fmt.Errorf("timeout not specified for redis")
	}

	return nil
}

type redisSentinelConfig struct {
	MasterName string   `yaml:"master_name"`
	Addresses  []string `yaml:"addresses"`
	DB         int      `yaml:"db"`
	Timeout    int      `yaml:"timeout"`
}

// Check checks the Redis sentinel configuration.
func (r *redisSentinelConfig) Check() error {
	if r.MasterName == "" {
		return fmt.Errorf("master name not specified for redis sentinel")
	}
	if len(r.Addresses) == 0 {
		return fmt.Errorf("addresses not specified for redis sentinel")
	}
	if r.Timeout == 0 {
		return fmt.Errorf("timeout not specified for redis")
	}

	return nil
}
