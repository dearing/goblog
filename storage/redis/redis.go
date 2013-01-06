package redis

import (
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"github.com/vmihailenco/redis"
	"html/template"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"
)

type Post struct {
	ID       string
	Title    string
	Content  template.HTML
	Created  time.Time
	Modified time.Time
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
	client.HSet(key, "created", fmt.Sprint(p.Created.Unix()))
	client.HSet(key, "modified", fmt.Sprint(p.Modified.Unix()))
	client.HIncrBy(key, "accessed", 1)

	z := redis.Z{Score: float64(p.Created.Unix()), Member: p.Title}
	client.ZAdd("posts", z)

	return e

}

func Get(id string, incr bool) (p Post, e error) {

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

	// Would think that there could be a mapping here in the github.com/vmihailenco/redis library?
	con := map[string]string{}
	for i := 0; i < len(v); i += 2 {
		con[v[i]] = v[i+1]
	}

	// Parse our UNIX timestamp (as string) into a int64 to be understood as a type of time
	created, e := strconv.ParseInt(con["created"], 10, 64)
	if e != nil {
		return p, e
	}
	mod, e := strconv.ParseInt(con["modified"], 10, 64)
	if e != nil {
		return p, e
	}

	p.ID = id
	p.Title = con["title"]
	p.Content = template.HTML(con["content"])
	p.Created = time.Unix(created, 0)
	p.Modified = time.Unix(mod, 0)
	p.Accessed = con["accessed"]

	if incr {
		client.HIncrBy(key, "accessed", 1)
	}

	return p, e
}

func Del(id string) (e error) {
	key := fmt.Sprintf("post:%s", id)
	client.Del(key)
	return e
}

// Return our post titles sorted in reverse of creation
func GetPosts() (keys *redis.StringSliceReq) {
	return client.ZRevRange("posts", "0", "-1")
}

// Return the latest post
func GetLatest() (p Post, e error) {

	a := client.ZRevRange("posts", "0", "0")
	if a.Err() != nil {
		log.Println(a.Err())
		return p, a.Err()
	}

	p, e = Get(a.Val()[0], false)

	return p, e
}

func getHTML(content string) template.HTML {
	return template.HTML(blackfriday.MarkdownCommon([]byte(content)))
}

func LoadDirectory(path string, suffix string) (e error) {
	x, e := ioutil.ReadDir(path)

	if e != nil {
		return e
	}

	for _, z := range x {
		if !z.IsDir() {

			b, e := ioutil.ReadFile(path + z.Name())
			if e != nil {
				log.Println(e)
				continue
			}

			id := strings.TrimRight(z.Name(), suffix)
			p := Post{
				ID:      id,
				Title:   id,
				Content: getHTML(string(b)),

				// TODO: figure out how the fuck to get the *created time* from FILE!
				// NOTE: seems to be an OS thang >> http://golang.org/pkg/os/#FileInfo 
				// -- Sys() interface{}   // underlying data source (can return nil) --
				Created:  time.Now(),
				Modified: z.ModTime(),
				Accessed: "0",
			}
			Set(p)
		}
	}

	return e
}
