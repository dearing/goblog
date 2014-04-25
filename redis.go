package main

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

var pool *redis.Pool

func newPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

type Page struct {
	UUID     string
	Title    string
	Content  string
	Author   string
	Created  uint64
	Modified uint64
	Views    int64
}

var Default = &Page{
	UUID:     "0",
	Title:    "not found",
	Content:  "no results",
	Author:   "server",
	Created:  0,
	Modified: 0,
	Views:    0,
}

func create() (p *Page) {
	p = &Page{
		UUID: uuid.New(),
	}

	log.Printf("%s create\n", p.UUID)
	return
}

func (p *Page) save() (err error) {
	c := pool.Get()
	defer c.Close()

	c.Send("hset", p.UUID, "Author", p.Author)
	c.Send("hset", p.UUID, "Content", p.Content)
	c.Send("hset", p.UUID, "Created", p.Created)
	c.Send("hset", p.UUID, "Title", p.Title)
	c.Send("hset", p.UUID, "Modified", p.Modified)
	c.Send("hset", p.UUID, "Views", p.Views)

	c.Flush()

	log.Printf("%s save\n", p.UUID)
	return nil
}

func (p *Page) load() (err error) {

	c := pool.Get()
	defer c.Close()

	reply, err := redis.Bool(c.Do("exists", p.UUID))
	if err != nil {
		return err
	}

	if reply {
		log.Printf("%s load\n", p.UUID)
	}

	return nil
}

func (p *Page) delete() (err error) {
	c := pool.Get()
	defer c.Close()

	reply, err := redis.Bool(c.Do("exists", p.UUID))
	if err != nil {
		return err
	}

	if reply {
		defer log.Printf("%s delete\n", p.UUID)
		c.Do("del", p.UUID)
	}

	return nil
}

func exists(uuid string) bool {
	c := pool.Get()
	defer c.Close()
	reply, _ := redis.Bool(c.Do("exists", uuid))
	return reply
}
