/*

	Jacob Dearing
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/russross/blackfriday"
	"github.com/vmihailenco/redis"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Represents and blog post.
// In the future timestamps, author information etc will be implemented.
type Article struct {
	Title string        // just the title
	Body  template.HTML // we consider the storage to be safe enough to generate HTML from (after markdown processing)
}

// Push a files contents up to the redis db after processing as markdown
func push(filename string) error {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	test := client.Set(filename, string(blackfriday.MarkdownCommon(body)))

	if test.Err() != nil {
		return test.Err()
	}

	if *verbose {
		log.Printf("pushed %s", filename)
	}

	return nil
}

// Reads a folder for files that are not folders itself (one level only) and pushes to the redis server as the folder + filename.
// TODO: needs a better naming scheme that will unfold when I get around to organizing data on the db
func pushall(folder string) error {

	files, _ := ioutil.ReadDir("articles")

	// for each file in the folder that, isn't a folder itself, push the parsed contents up
	for _, file := range files {
		if !file.IsDir() {
			err := push("articles/" + file.Name())
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

// Load an article to the db by reading it then parsing it as markdown and then pushing it up to the redis db by name
func pull(title string) (*Article, error) {

	key := fmt.Sprintf(title)

	if !client.Exists(key).Val() {
		if *verbose {
			log.Printf("not found: %v", key)
			return nil, errors.New(fmt.Sprintf("db does not contain key: %v", key))
		}
	}

	output := client.Get(key).Val()
	return &Article{Title: title, Body: template.HTML(output)}, nil
}

// Generate a simple list of article titles and links from redis.
// TODO: better naming scheme as at PUSHALL
func tocHandler(w http.ResponseWriter, r *http.Request) {
	title := "table of contents"

	keys := client.Keys("articles/*")

	templates.ExecuteTemplate(w, "head", title)
	templates.ExecuteTemplate(w, "bar", nil)
	templates.ExecuteTemplate(w, "toc-head", nil)

	// for each key we add a list element
	for _, element := range keys.Val() {
		url := strings.Replace(element, ".text", "", 1)
		url = strings.Replace(url, "articles/", "", 1)
		templates.ExecuteTemplate(w, "toc-item", url)
	}

	templates.ExecuteTemplate(w, "toc-foot", nil)
}

// Load and display an article from our redis db.
func articleHandler(w http.ResponseWriter, r *http.Request) {
	// Extract a meaningful title from the path.
	title := r.URL.Path[len("/article/"):]

	p, err := pull("articles/"+title+".text")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		log.Printf("error : %v\n", err)
		return
	}

	//t, err := template.ParseFiles("templates/common.html", "templates/article.html")
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		log.Printf("error : %v\n", err)
		return
	}

	templates.ExecuteTemplate(w, "head", p)
	templates.ExecuteTemplate(w, "bar", p)
	templates.ExecuteTemplate(w, "article", p)
	templates.ExecuteTemplate(w, "foot", p)
}

// ARGS
var host = flag.String("host", ":8080", "host to bind to")
var root = flag.String("root", "wwwroot", "webserver document root folder")
var redis_host = flag.String("rh", "192.168.1.150:6379", "redis host")
var redis_pass = flag.String("rp", "", "redis password")
var redis_db = flag.Int64("rdb", -1, "redis db index")
var verbose = flag.Bool("verbose", false, "log common operations and not just errors")

// MISC
var templates = template.Must(template.ParseGlob("templates/*.html"))
var client = redis.NewTCPClient(*redis_host, *redis_pass, *redis_db)

//  MAIN
func main() {

	// First we parse our env args for use down road.
	flag.Parse()

	// Initialize contact with the server using our arguments or defaults.
	//client = redis.NewTCPClient(*redis_host, *redis_pass, *redis_db)

	// If we can ping wihtout an error then we can move on.
	if ping := client.Ping(); ping.Err() != nil {
		log.Panicf("%v", ping.Err())
	}

	defer client.Close()

	// Straight from the author of github.com/howeyc/fsnotify:
	// Initialize our watcher and check for any errors...
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	// Spin up a goroutine that watches two channels:
	// watcher.Event for events of [Delete, Modify, Moved, New] and
	// watcher.Error for any errors behind the scenes.
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if *verbose {
					log.Println("event:", ev)
				}
				pull(ev.Name)
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	// Watch our articles for changes
	err = watcher.Watch("articles")
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Push our working articles to our redis db
	pushall("articles")

	//	Setup our handlers and get cracking...
	http.Handle("/", http.FileServer(http.Dir(*root)))
	http.HandleFunc("/article/", articleHandler)
	http.HandleFunc("/toc/", tocHandler)

	fmt.Printf("listening on %s // root=%s\n", *host, *root)

	if err = http.ListenAndServe(*host, nil); err != nil {
		log.Fatalf("%v", err)
	}
}
