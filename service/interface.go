package service

import (
	"github.com/Xuanwo/migrant/constants"
	"github.com/Xuanwo/migrant/model"
	"github.com/Xuanwo/migrant/service/fs"
	"github.com/Xuanwo/migrant/service/mysql"
	"github.com/Xuanwo/migrant/service/redis"
)

// Service is represent a service.
type Service interface {
	Name() string
}

// Source is a service that provides migration sources.
type Source interface {
	Service

	List() ([]model.Record, error)
	Read(id, t string) (content []byte, err error)
}

// Recorder is a service that provides migration records.
type Recorder interface {
	Service

	List() ([]model.Record, error)
	Write(id, t string) error
	Delete(id, t string) error
}

// Migrant is a service that operate migration.
type Migrant interface {
	Service

	Up(content []byte) (err error)
	Down(content []byte) (err error)
}

// New will create a new service.
func New(t string, opt []byte) (Service, error) {
	switch t {
	case constants.ServiceFs:
		return fs.New(opt)
	case constants.ServiceMySQL:
		return mysql.New(opt)
	case constants.ServiceRedis:
		return redis.New(constants.ServiceRedis, opt)
	case constants.ServiceRedisSentinel:
		return redis.New(constants.ServiceRedisSentinel, opt)
	default:
		return nil, constants.ErrServiceInvalid
	}
}
