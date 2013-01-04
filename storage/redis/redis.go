package redis

import (
	"errors"
	"fmt"
	"github.com/vmihailenco/redis"
	"io/ioutil"
)

type Post struct {
	ID       string
	Title    string // just the title
	Content  string // we consider the storage to be safe enough to generate HTML from (after markdown processing)
	Author   string
	Created  string
	Modified string
	Accessed string
}

var client *redis.Client

// TODO: error handling on failed connection
func Connect(host string, pass string, db int64) (e error) {
	client = redis.NewTCPClient(host, pass, db)
	return e
}

func Close() (e error) {
	e = client.Close()
	return e
}

// TODO: error handling
func Set(p Post) (e error) {

	//log.Print(p)
	key := fmt.Sprintf("post:" + p.ID)

	client.HSet(key, "title", p.Title)
	client.HSet(key, "content", string(p.Content))
	client.HSet(key, "author", p.Author)
	client.HSet(key, "created", string(p.Created))
	client.HSet(key, "modified", string(p.Modified))
	client.HIncrBy(key, "accessed", 1)

	return e

}

func Get(id string) (p Post, e error) {

	key := fmt.Sprintf("post:%s", id)

	if !client.Exists(key).Val() {
		return p, errors.New("key doesn't exist : " + key)
	}

	get := client.HGetAll(key)
	e = get.Err()
	if e != nil {
		return p, e
	}

	v := get.Val()

	// Build our post now
	// Would think that there could be a mapping here in the github.com/vmihailenco/redis library?
	// BUG(dearing): HASHES are unsorted so this should fail at some point down the road.
	p.ID = id
	p.Title = v[1]
	p.Content = v[3]
	p.Author = v[5]
	p.Created = v[7]
	p.Modified = v[9]
	p.Accessed = v[11]

	return p, e
}

func Del(id string) (e error) {
	key := fmt.Sprintf("post:%s", id)
	client.Del(key)
	return e
}

func Keys(pattern string) (keys *redis.StringSliceReq) {
	return client.Keys(pattern)
}

func LoadDirectory(path string) (e error) {

	x, e := ioutil.ReadDir(path)

	if e != nil {
		return e
	}

	for _, z := range x {
		if !z.IsDir() {

			b, _ := ioutil.ReadFile(z.Name())

			id := client.Incr("global:nextPostID")
			p := Post{
				ID:       fmt.Sprintf("%v", id.Val()),
				Title:    z.Name(),
				Content:  string(b),
				Created:  fmt.Sprintf("%v", z.ModTime().Unix()),
				Modified: fmt.Sprintf("%v", z.ModTime().Unix()),
				Accessed: "0",
			}
			Set(p)
		}
	}

	return e
}
