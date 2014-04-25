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
	Created  int64
	Modified int64
	Views    int64
}

func create() (p *Page) {
	p = &Page{
		UUID: uuid.New(),
	}

	p.Created = time.Now().Unix()
	p.Modified = time.Now().Unix()
	p.Views = 1

	log.Printf("%s create\n", p.UUID)
	return
}

func (p *Page) save() (err error) {
	c := pool.Get()
	defer c.Close()

	_, err = c.Do("HMSET", redis.Args{}.Add(p.UUID).AddFlat(p)...)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("%s save\n", p.UUID)
	return nil
}

func (p *Page) load() (err error) {

	c := pool.Get()
	defer c.Close()

	if exists(p.UUID) {
		c.Do("HINCRBY", p.UUID, "Views", 1)
		reply, _ := redis.Values(c.Do("HGETALL", p.UUID))
		redis.ScanStruct(reply, p)
		log.Printf("%s load\n", p.UUID)
	}

	return nil
}

func (p *Page) delete() (err error) {
	c := pool.Get()
	defer c.Close()

	reply, err := redis.Bool(c.Do("EXISTS", p.UUID))
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
