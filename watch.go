package main

import (
	"github.com/howeyc/fsnotify"
	"log"
)

var watcher fsnotify.Watcher

func watch(path string) {

	// Watch our content for changes
	// Straight from the author of github.com/howeyc/fsnotify:
	// Initialize our watcher and check for any errors...
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panicln(err)
	}

	// Spin up a goroutine that watches two channels:
	// watcher.Event for events of [Delete, Modify, Moved, New] and
	// watcher.Error for any errors behind the scenes.
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if config.Verbose {
					log.Printf("event:%v", ev)
				}
				if ev.IsModify() || ev.IsCreate() {
					err := push(ev.Name)
					if err != nil {
						log.Println(err)
					}
				}
				if ev.IsDelete() {
					err := drop(ev.Name)
					if err != nil {
						log.Println(err)
					}

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
				log.Panicln("error:", err)
			}
		}
	}()

	// Watch our articles for changes
	err = watcher.Watch(path)
	if err != nil {
		log.Fatal(err)
	}
}
