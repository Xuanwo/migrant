package mysql

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/Xuanwo/migrant/common/db"
	"github.com/Xuanwo/migrant/model"
	"github.com/Xuanwo/migrant/service/mysql/sqlparse"
)

// Client is the struct for local file list endpoint.
type Client struct {
	client *db.MySQL

	migrationTable string
}

// Name implement service.Name
func (c *Client) Name() (name string) {
	return "mysql:" + c.client.Name()
}

// New create a new service client.
func New(opt []byte) (c *Client, err error) {
	c = &Client{}

	cfg := &mysqlConfig{}
	err = yaml.Unmarshal(opt, cfg)
	if err != nil {
		return
	}

	err = cfg.check()
	if err != nil {
		return
	}

	c.client, err = db.NewMySQL(&db.MySQLOptions{
		Address: fmt.Sprintf(
			"%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Database:           cfg.Database,
		User:               cfg.User,
		Password:           cfg.Password,
		ConnectionTimeout:  cfg.Timeout,
		MaxConnections:     cfg.MaxConnections,
		MaxIdleConnections: cfg.MaxIdleConnections,
	})
	if err != nil {
		return
	}

	c.migrationTable = cfg.MigrationTable

	if c.migrationTable != "" {
		_, err = c.client.Exec(schemaMigrationTable)
		if err != nil {
			return
		}
	}

	return
}

// Up implement service.Up
func (c *Client) Up(content []byte) (err error) {
	buf := bytes.NewReader(content)

	m, err := sqlparse.ParseMigration(buf)
	if err != nil {
		return
	}

	ctx := context.Background()

	tx, err := c.client.NewTx(ctx)
	if err != nil {
		return
	}
	defer tx.Close()

	for _, v := range m.UpStatements {
		_, err = tx.Exec(v)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}

	return
}

// Down implement service.Down
func (c *Client) Down(content []byte) (err error) {
	buf := bytes.NewReader(content)

	m, err := sqlparse.ParseMigration(buf)
	if err != nil {
		return
	}

	ctx := context.Background()

	tx, err := c.client.NewTx(ctx)
	if err != nil {
		return
	}
	defer tx.Close()

	for _, v := range m.DownStatements {
		_, err = tx.Exec(v)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}

	return
}

// List implement service.List
func (c *Client) List() (m []model.Record, err error) {
	if c.migrationTable == "" {
		// TODO: we should return error here.
		return nil, nil
	}

	m = []model.Record{}
	err = c.client.SelectFrom(c.migrationTable).All(&m)
	if err != nil {
		return
	}

	return
}

// Write implement service.Write
func (c *Client) Write(id, t string) (err error) {
	if c.migrationTable == "" {
		// TODO: we should return error here.
		return nil
	}

	m := &model.Record{
		ID:        id,
		Type:      t,
		AppliedAt: time.Now().Unix(),
	}
	_, err = c.client.InsertInto(c.migrationTable).Values(m).Exec()
	if err != nil {
		return
	}
	return
}

// Delete implement service.Delete
func (c *Client) Delete(id, t string) (err error) {
	if c.migrationTable == "" {
		// TODO: we should return error here.
		return nil
	}

	cond := db.Cond{
		"id":   id,
		"type": t,
	}
	_, err = c.client.DeleteFrom(c.migrationTable).Where(cond).Exec()
	if err != nil {
		return
	}
	return
}
