package db

import (
	"time"
	"upper.io/db.v3"

	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
)

// Cond is a map that defines conditions for a query and satisfies the
// Constraints and Compound interfaces.
//
// Each entry of the map represents a condition (a column-value relation bound
// by a comparison operator). The comparison operator is optional and can be
// specified after the column name, if no comparison operator is provided the
// equality is used.
//
// Examples:
//
//  // Where age equals 18.
//  db.Cond{"age": 18}
//  //	// Where age is greater than or equal to 18.
//  db.Cond{"age >=": 18}
//
//  // Where id is in a list of ids.
//  db.Cond{"id IN": []{1, 2, 3}}
//
//  // Where age is lower than 18 (you could use this syntax when using
//  // mongodb).
//  db.Cond{"age $lt": 18}
//
//  // Where age > 32 and age < 35
//  db.Cond{"age >": 32, "age <": 35}
type Cond = db.Cond

// MySQL is a MySQL client.
type MySQL struct {
	sqlbuilder.Database
}

// MySQLOptions is options for MySQL client.
type MySQLOptions struct {
	Address  string
	Database string

	User     string
	Password string

	ConnectionTimeout  int
	MaxConnections     int
	MaxIdleConnections int
}

// NewMySQL creates new MySQL client.
func NewMySQL(opt *MySQLOptions) (m *MySQL, err error) {
	client, err := mysql.Open(mysql.ConnectionURL{
		Host:     opt.Address,
		Database: opt.Database,
		User:     opt.User,
		Password: opt.Password,
		Options:  map[string]string{"charset": "utf8mb4,utf8"},
	})
	if err != nil {
		return
	}

	client.SetConnMaxLifetime(time.Duration(opt.ConnectionTimeout) * time.Second)
	client.SetMaxOpenConns(opt.MaxConnections)
	client.SetMaxIdleConns(opt.MaxIdleConnections)

	m = &MySQL{Database: client}
	return
}
