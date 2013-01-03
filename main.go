/*
	Copyright (c) 2012 Jacob Dearing
*/
package main

import (
	"flag"
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/vmihailenco/redis"
	"html/template"
	"log"
	"net/http"
)

// Represents and blog post.
// In the future timestamps, author information etc will be implemented.
type Article struct {
	Title string        // just the title
	Body  template.HTML // we consider the storage to be safe enough to generate HTML from (after markdown processing)
}

// ARGS
var articles = flag.String("articles", "articles", "markdown posts")
var templates = flag.String("templates", "templates", "templates posts")
var suffix = flag.String("suffix", ".md", "filtered extension")
var host = flag.String("wwwhost", ":8080", "host to bind to")
var root = flag.String("wwwroot", "wwwroot", "webserver document root folder")
var redis_host = flag.String("redis-host", "localhost:6379", "redis host")
var redis_pass = flag.String("redis-pass", "", "redis password")
var redis_db = flag.Int64("redis-db", -1, "redis db index")
var verbose = flag.Bool("verbose", false, "log common operations and not just errors")

var client *redis.Client

//  MAIN
func main() {

	// First we parse our env args for use down road.
	flag.Parse()

	// Initialize contact with the server using our arguments or defaults.
	client = redis.NewTCPClient(*redis_host, *redis_pass, *redis_db)

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
					log.Printf("event:%v", ev)
				}
				if ev.IsModify() || ev.IsCreate() {
					push(ev.Name)
				}
				if ev.IsDelete() {
					drop(ev.Name)
				}

				// TODO: Need to work out how to know the old name and the new.
				// If it isn't supported in fsnotify then I'll need to make
				// a routine scan the present state and compare that to the
				// database.
				/*
					if ev.IsRename() {
						push(ev.Name)
					}
				*/

			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	// Watch our articles for changes
	err = watcher.Watch(*articles)
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Push our working articles to our redis db
	pushall(*articles)

	//	Setup our handlers and get cracking...
	http.Handle("/static/", http.FileServer(http.Dir(*root)))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/p/", articleHandler)
	http.HandleFunc("/toc/", tocHandler)

	fmt.Printf("listening on %s // root=%s\n", *host, *root)

	if err = http.ListenAndServe(*host, nil); err != nil {
		log.Fatalf("%v", err)
	}
}
