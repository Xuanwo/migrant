package fs

import (
	"github.com/Xuanwo/migrant/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

// Client is the struct for local file list endpoint.
type Client struct {
	Path string `yaml:"path"`
}

// Name implement service.Name
func (c *Client) Name() (name string) {
	return "fs:" + c.Path
}

// New create a new service client.
func New(opt []byte) (c *Client, err error) {
	c = &Client{}

	err = yaml.Unmarshal(opt, c)
	if err != nil {
		return
	}

	return
}

// List implement service.List
func (c *Client) List() (r []model.Record, err error) {
	files, err := ioutil.ReadDir(c.Path)
	if err != nil {
		return
	}

	r = make([]model.Record, len(files))
	for k, v := range files {
		name := strings.Split(v.Name(), ".")

		if len(name) != 2 {
			// TODO: we should return an error here.
			return nil, nil
		}

		r[k] = model.Record{
			ID:   name[0],
			Type: name[1],
		}
	}

	return
}

// Read implement service.Read
func (c *Client) Read(id, t string) (content []byte, err error) {
	name := strings.Join([]string{id, t}, ".")

	path := filepath.Join(c.Path, name)

	content, err = ioutil.ReadFile(path)
	if err != nil {
		return
	}
	return
}
