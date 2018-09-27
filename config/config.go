package config

import (
	"fmt"
	"io/ioutil"
	"os/user"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/Xuanwo/migrant/constants"
)

// Config stores a configuration.
type Config struct {
	Source string `yaml:"source"`
	Record string `yaml:"record"`

	MySQLMigrant string `yaml:"mysql_migrant"`
	RedisMigrant string `yaml:"redis_migrant"`

	Services map[string]Service `yaml:"services"`
}

// New create a global config.
func New() (*Config, error) {
	return &Config{}, nil
}

// LoadFromFilePath loads configuration from a specified local path.
// It returns error if file not found or yaml decode failed.
func (s *Config) LoadFromFilePath(filePath string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	if strings.Index(filePath, "~/") == 0 {
		filePath = strings.Replace(filePath, "~/", usr.HomeDir+"/", 1)
	}

	c, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	return s.LoadFromContent(c)
}

// LoadFromContent loads configuration from a given bytes.
// It returns error if yaml decode failed.
func (s *Config) LoadFromContent(content []byte) error {
	return yaml.Unmarshal(content, s)
}

// Check checks the configuration.
func (s *Config) Check() error {
	if len(s.Services) == 0 {
		return fmt.Errorf("services are not set")
	}

	for _, v := range s.Services {
		if err := v.Check(); err != nil {
			return constants.ErrServiceInvalid
		}
	}

	// Check source.
	if s.Source == "" {
		return fmt.Errorf("source is not set")
	}
	_, ok := s.Services[s.Source]
	if !ok {
		return fmt.Errorf("source is not found")
	}

	// Check record.
	if s.Record == "" {
		return fmt.Errorf("record is not set")
	}
	_, ok = s.Services[s.Record]
	if !ok {
		return fmt.Errorf("record is not found")
	}

	// Check mysql migrant.
	if s.MySQLMigrant != "" {
		_, ok = s.Services[s.MySQLMigrant]
		if !ok {
			return fmt.Errorf("mysql migrant is not found")
		}
	}

	// Check redis migrant.
	if s.RedisMigrant != "" {
		_, ok = s.Services[s.RedisMigrant]
		if !ok {
			return fmt.Errorf("redis migrant is not found")
		}
	}

	return nil
}
