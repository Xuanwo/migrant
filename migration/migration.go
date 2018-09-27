package migration

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/Xuanwo/migrant/config"
	"github.com/Xuanwo/migrant/constants"
	"github.com/Xuanwo/migrant/model"
	"github.com/Xuanwo/migrant/service"
	"github.com/Xuanwo/migrant/service/mysql"
	"github.com/Xuanwo/migrant/service/redis"
)

var (
	src service.Source
	rec service.Recorder

	mysqlSrv service.Migrant
	redisSrv service.Migrant

	services map[string]service.Service

	sourceMigrations []model.Record
	existMigrations  []model.Record
)

// Setup will setup migration for migrant.
func Setup(c *config.Config) (err error) {
	services = make(map[string]service.Service)

	for k, v := range c.Services {
		content, err := yaml.Marshal(v.Options)
		if err != nil {
			return err
		}

		services[k], err = service.New(v.Type, content)
		if err != nil {
			return err
		}
	}

	var ok bool

	src, ok = services[c.Source].(service.Source)
	if !ok {
		return fmt.Errorf("service %s does not implement service.Source", c.Source)
	}

	rec, ok = services[c.Record].(service.Recorder)
	if !ok {
		return fmt.Errorf("service %s does not implement service.Recorder", c.Record)
	}

	if c.MySQLMigrant != "" {
		mysqlSrv, ok = services[c.MySQLMigrant].(*mysql.Client)
		if !ok {
			return fmt.Errorf("service %s is not a mysql service", c.MySQLMigrant)
		}
	}

	if c.RedisMigrant != "" {
		redisSrv, ok = services[c.RedisMigrant].(*redis.Client)
		if !ok {
			return fmt.Errorf("service %s is not a redis service", c.RedisMigrant)
		}
	}

	sourceMigrations, err = src.List()
	if err != nil {
		return
	}

	existMigrations, err = rec.List()
	if err != nil {
		return
	}

	if len(existMigrations) > len(sourceMigrations) {
		return constants.ErrMigrationMissing
	}

	for i := 0; i < len(existMigrations); i++ {
		sr := sourceMigrations[i]
		er := existMigrations[i]

		if er.ID != sr.ID {
			fmt.Printf("%s is not match with %s.", er.ID, sr.ID)
			return constants.ErrMigrationMismatch
		}
	}

	return nil
}

// Up will do the up migration.
func Up() (id string, err error) {
	if len(sourceMigrations) == len(existMigrations) {
		return "", nil
	}

	m := sourceMigrations[len(existMigrations)]

	content, err := src.Read(m.ID, m.Type)
	if err != nil {
		return
	}

	switch m.Type {
	case constants.RecordSQL:
		err = mysqlSrv.Up(content)
		if err != nil {
			return
		}
	case constants.RecordRedis:
		err = redisSrv.Up(content)
		if err != nil {
			return
		}
	default:
		return "", constants.ErrMigrationNotSupported
	}

	err = rec.Write(m.ID, m.Type)
	if err != nil {
		return
	}

	id = m.ID
	return
}

// Down will do the down migration.
func Down() (id string, err error) {
	l := len(existMigrations)

	if l == 0 {
		return "", nil
	}

	m := existMigrations[l-1]

	content, err := src.Read(m.ID, m.Type)
	if err != nil {
		return
	}

	switch m.Type {
	case constants.RecordSQL:
		err = mysqlSrv.Down(content)
		if err != nil {
			return
		}
	case constants.RecordRedis:
		err = redisSrv.Down(content)
		if err != nil {
			return
		}
	default:
		return "", constants.ErrMigrationNotSupported
	}

	err = rec.Delete(m.ID, m.Type)
	if err != nil {
		return
	}

	id = m.ID
	return
}

// Status will show current status.
func Status() (m []model.Record, err error) {
	if len(sourceMigrations) == len(existMigrations) {
		return existMigrations, nil
	}
	return append(existMigrations, sourceMigrations[len(existMigrations):]...), nil
}

// Sync will sync the migration.
func Sync() (m []model.Record, err error) {
	if len(sourceMigrations) == len(existMigrations) {
		return nil, nil
	}

	ms := sourceMigrations[len(existMigrations):]

	for _, v := range ms {
		content, err := src.Read(v.ID, v.Type)
		if err != nil {
			return nil, err
		}

		switch v.Type {
		case constants.RecordSQL:
			err = mysqlSrv.Up(content)
			if err != nil {
				return nil, err
			}
		case constants.RecordRedis:
			err = redisSrv.Up(content)
			if err != nil {
				return nil, err
			}
		default:
			return nil, constants.ErrMigrationNotSupported
		}

		err = rec.Write(v.ID, v.Type)
		if err != nil {
			return nil, err
		}
	}

	return
}
