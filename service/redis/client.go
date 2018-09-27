package redis

import (
	"bytes"
	"fmt"
	"log"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/Xuanwo/migrant/common/cache"
	"github.com/Xuanwo/migrant/constants"
)

// Client is the struct for local file list endpoint.
type Client struct {
	client *cache.Redis
}

// Name implement service.Name
func (c *Client) Name() (name string) {
	return "redis:" + c.client.Client.String()
}

// New create a new service client.
func New(t string, opt []byte) (c *Client, err error) {
	c = &Client{}

	switch t {
	case constants.ServiceRedis:
		cfg := &redisConfig{}
		err = yaml.Unmarshal(opt, cfg)
		if err != nil {
			return
		}

		c.client, err = cache.NewRedis(&cache.RedisOptions{
			Address: fmt.Sprintf(
				"%s:%d", cfg.Host, cfg.Port,
			),
			DB:      cfg.DB,
			Timeout: cfg.Timeout,
		})
		if err != nil {
			return
		}
	case constants.ServiceRedisSentinel:
		cfg := &redisSentinelConfig{}
		err = yaml.Unmarshal(opt, cfg)
		if err != nil {
			return
		}

		c.client, err = cache.NewRedisSentinel(&cache.RedisSentinelOptions{
			Addresses:  cfg.Addresses,
			MasterName: cfg.MasterName,
			DB:         cfg.DB,
			Timeout:    cfg.Timeout,
		})
		if err != nil {
			return
		}
	}

	return
}

// Up implement service.Up
func (c *Client) Up(content []byte) (err error) {
	m := &Migration{}

	err = yaml.Unmarshal(content, m)
	if err != nil {
		return
	}

	var arr [128]string
	keys := arr[0:0]
	buf := bytes.NewBuffer(make([]byte, 1024))

	for _, j := range m.Up {
		ch := make(chan string, 100)
		go c.scan(j.Pattern, ch)

		var tmpl []*template.Template
		for _, action := range j.Actions {
			nt := template.New("")
			nt, err = nt.Parse(action)
			if err != nil {
				return
			}

			tmpl = append(tmpl, nt)
		}

		for v := range ch {
			keys = append(keys, v)

			if len(keys) >= 100 {
				for _, t := range tmpl {
					buf.Reset()

					err = t.Execute(buf, keys)
					if err != nil {
						return
					}

					_, err = c.client.Do(convert(buf.String())...).Result()
					if err != nil {
						return
					}
				}

				keys = arr[0:0]
			}
		}

		if len(keys) > 0 {
			for _, t := range tmpl {
				buf.Reset()

				err = t.Execute(buf, keys)
				if err != nil {
					return
				}

				_, err = c.client.Do(convert(buf.String())...).Result()
				if err != nil {
					return
				}
			}

			keys = arr[0:0]
		}
	}

	return
}

// Down implement service.Down
func (c *Client) Down(content []byte) (err error) {
	m := &Migration{}

	err = yaml.Unmarshal(content, m)
	if err != nil {
		return
	}

	var arr [128]string
	keys := arr[0:0]
	buf := bytes.NewBuffer(make([]byte, 1024))

	for _, j := range m.Down {
		ch := make(chan string, 100)
		go c.scan(j.Pattern, ch)

		var tmpl []*template.Template
		for _, action := range j.Actions {
			nt := template.New("")
			nt, err = nt.Parse(action)
			if err != nil {
				return
			}

			tmpl = append(tmpl, nt)
		}

		for v := range ch {
			keys = append(keys, v)

			if len(keys) >= 100 {
				for _, t := range tmpl {
					buf.Reset()

					err = t.Execute(buf, keys)
					if err != nil {
						return
					}

					_, err = c.client.Do(convert(buf.String())...).Result()
					if err != nil {
						return
					}
				}

				keys = arr[0:0]
			}
		}

		if len(keys) > 0 {
			for _, t := range tmpl {
				buf.Reset()

				err = t.Execute(buf, keys)
				if err != nil {
					return
				}

				_, err = c.client.Do(convert(buf.String())...).Result()
				if err != nil {
					return
				}
			}

			keys = arr[0:0]
		}
	}

	return
}

func (c *Client) scan(pattern string, ch chan string) {
	defer close(ch)

	var keys []string
	var err error

	it := uint64(0)

	for {
		keys, it, err = c.client.Scan(it, pattern, 1000).Result()
		if err != nil {
			log.Fatalf("Redis failed err %v.", err)
		}

		for _, v := range keys {
			ch <- v
		}

		if it == 0 {
			break
		}
	}
}
