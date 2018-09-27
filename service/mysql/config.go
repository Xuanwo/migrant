package mysql

import "fmt"

type mysqlConfig struct {
	Host               string `yaml:"host"`
	Port               int    `yaml:"port"`
	Database           string `yaml:"database"`
	User               string `yaml:"user"`
	Password           string `yaml:"password"`
	Timeout            int    `yaml:"timeout"`
	MaxConnections     int    `yaml:"max_connections"`
	MaxIdleConnections int    `yaml:"max_idle_connections"`

	MigrationTable string `yaml:"migration_table"`
}

func (c *mysqlConfig) check() error {
	if c.Host == "" {
		return fmt.Errorf("host not specified for MySQL")
	}
	if c.Port == 0 {
		return fmt.Errorf("port not specified for MySQL")
	}
	if c.Database == "" {
		return fmt.Errorf("database not specified for MySQL")
	}
	if c.User == "" {
		return fmt.Errorf("user not specified for MySQL")
	}
	if c.Password == "" {
		return fmt.Errorf("password not specified for MySQL")
	}
	if c.Timeout == 0 {
		return fmt.Errorf("timeout not specified for MySQL")
	}
	if c.MaxConnections == 0 {
		return fmt.Errorf("max connections not specified for MySQL")
	}
	if c.MaxIdleConnections == 0 {
		return fmt.Errorf("max idle connections not specified for MySQL")
	}

	return nil
}
